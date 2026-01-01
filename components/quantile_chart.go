package components

import (
	"fmt"
	"math"
	"strings"

	"github.com/HatiCode/kedastral-tui/client"
	"github.com/charmbracelet/lipgloss"
)

// QuantileChart renders a forecast chart with P10/P50/P90 quantiles.
type QuantileChart struct {
	width, height int
}

// NewQuantileChart creates a new quantile chart.
func NewQuantileChart(width, height int) *QuantileChart {
	return &QuantileChart{width: width, height: height}
}

// Render renders the quantile chart.
func (c *QuantileChart) Render(snapshot *client.QuantileSnapshotData) string {
	if snapshot == nil || len(snapshot.Snapshot.DesiredReplicas) == 0 {
		return c.renderEmpty()
	}

	// Check if quantiles are available (API v2)
	if len(snapshot.Snapshot.Quantiles) > 0 && snapshot.APIVersion >= 2 {
		return c.renderQuantiles(snapshot)
	}

	// Fallback to single line (API v1)
	return c.renderSingleLine(snapshot)
}

// renderQuantiles renders P10/P50/P90 quantile lines.
func (c *QuantileChart) renderQuantiles(snapshot *client.QuantileSnapshotData) string {
	p10, hasP10 := snapshot.Snapshot.Quantiles["p10"]
	p50, hasP50 := snapshot.Snapshot.Quantiles["p50"]
	p90, hasP90 := snapshot.Snapshot.Quantiles["p90"]

	if !hasP50 {
		return c.renderEmpty()
	}

	// Find global min/max across all quantiles
	minVal, maxVal := findMinMaxAcross(p10, p50, p90)

	if minVal == maxVal {
		maxVal = minVal + 1
	}

	// Create chart grid
	chartHeight := c.height - 4 // Reserve space for title and legend
	chartWidth := c.width - 10  // Reserve space for Y-axis labels

	var lines []string
	titleStyle := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("39"))
	lines = append(lines, titleStyle.Render("Forecast Timeline (P10/P50/P90)"))
	lines = append(lines, "")

	// Render Y-axis and plot
	for row := chartHeight; row >= 0; row-- {
		yVal := minVal + float64(row)/float64(chartHeight)*(maxVal-minVal)

		// Y-axis label (every few rows)
		yLabel := ""
		if row == chartHeight || row == chartHeight/2 || row == 0 {
			yLabel = fmt.Sprintf("%6.1f", yVal)
		} else {
			yLabel = "      "
		}

		// Plot line
		var plotLine strings.Builder
		plotLine.WriteString(" ")

		for col := 0; col < chartWidth; col++ {
			idx := int(float64(col) / float64(chartWidth) * float64(len(p50)-1))
			if idx >= len(p50) {
				idx = len(p50) - 1
			}

			// Check which quantile to plot at this position
			char := " "
			style := lipgloss.NewStyle()

			// Determine Y position for each quantile
			p10Y := -1.0
			if hasP10 && idx < len(p10) {
				p10Y = (p10[idx] - minVal) / (maxVal - minVal) * float64(chartHeight)
			}
			p50Y := (p50[idx] - minVal) / (maxVal - minVal) * float64(chartHeight)
			p90Y := -1.0
			if hasP90 && idx < len(p90) {
				p90Y = (p90[idx] - minVal) / (maxVal - minVal) * float64(chartHeight)
			}

			// Plot character based on which line is closest to current row
			currentY := float64(row)
			minDist := 999999.0
			plotType := ""

			if hasP10 && p10Y >= 0 && math.Abs(currentY-p10Y) < minDist {
				minDist = math.Abs(currentY - p10Y)
				plotType = "p10"
			}
			if math.Abs(currentY-p50Y) < minDist {
				minDist = math.Abs(currentY - p50Y)
				plotType = "p50"
			}
			if hasP90 && p90Y >= 0 && math.Abs(currentY-p90Y) < minDist {
				minDist = math.Abs(currentY - p90Y)
				plotType = "p90"
			}

			if minDist < 0.5 {
				switch plotType {
				case "p10":
					char = "·"
					style = style.Foreground(lipgloss.Color("39")) // Cyan
				case "p50":
					char = "●"
					style = style.Foreground(lipgloss.Color("42")) // Green
				case "p90":
					char = "■"
					style = style.Foreground(lipgloss.Color("196")) // Red
				}
			}

			plotLine.WriteString(style.Render(char))
		}

		lines = append(lines, yLabel+"┤"+plotLine.String())
	}

	// X-axis
	xAxis := "       └" + strings.Repeat("─", chartWidth)
	lines = append(lines, xAxis)

	// X-axis labels
	xLabels := "        Now" + strings.Repeat(" ", chartWidth-20) + fmt.Sprintf("+%dm", snapshot.Snapshot.HorizonSeconds/60)
	lines = append(lines, xLabels)

	// Legend
	p10Style := lipgloss.NewStyle().Foreground(lipgloss.Color("39"))
	p50Style := lipgloss.NewStyle().Foreground(lipgloss.Color("42"))
	p90Style := lipgloss.NewStyle().Foreground(lipgloss.Color("196"))

	legend := fmt.Sprintf("Legend: %s P10  %s P50  %s P90",
		p10Style.Render("···"),
		p50Style.Render("●●●"),
		p90Style.Render("■■■"),
	)
	lines = append(lines, "")
	lines = append(lines, legend)

	return strings.Join(lines, "\n")
}

// renderSingleLine renders a single forecast line (API v1 fallback).
func (c *QuantileChart) renderSingleLine(snapshot *client.QuantileSnapshotData) string {
	values := snapshot.Snapshot.Values
	if len(values) == 0 {
		// Try to get P50 from quantiles
		if p50, ok := snapshot.Snapshot.Quantiles["p50"]; ok {
			values = p50
		} else {
			return c.renderEmpty()
		}
	}

	minVal, maxVal := findMinMaxQuantile(values)
	if minVal == maxVal {
		maxVal = minVal + 1
	}

	chartHeight := c.height - 4
	chartWidth := c.width - 10

	var lines []string
	titleStyle := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("39"))
	warningStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("220"))

	lines = append(lines, titleStyle.Render("Forecast Timeline"))
	lines = append(lines, warningStyle.Render("⚠ Quantiles unavailable. Showing single-point forecast."))
	lines = append(lines, "")

	// Render Y-axis and plot
	for row := chartHeight; row >= 0; row-- {
		yVal := minVal + float64(row)/float64(chartHeight)*(maxVal-minVal)

		yLabel := ""
		if row == chartHeight || row == chartHeight/2 || row == 0 {
			yLabel = fmt.Sprintf("%6.1f", yVal)
		} else {
			yLabel = "      "
		}

		var plotLine strings.Builder
		plotLine.WriteString(" ")

		for col := 0; col < chartWidth; col++ {
			idx := int(float64(col) / float64(chartWidth) * float64(len(values)-1))
			if idx >= len(values) {
				idx = len(values) - 1
			}

			valueY := (values[idx] - minVal) / (maxVal - minVal) * float64(chartHeight)
			currentY := float64(row)

			if math.Abs(currentY-valueY) < 0.5 {
				plotLine.WriteString("●")
			} else {
				plotLine.WriteString(" ")
			}
		}

		lines = append(lines, yLabel+"┤"+plotLine.String())
	}

	// X-axis
	xAxis := "       └" + strings.Repeat("─", chartWidth)
	lines = append(lines, xAxis)

	xLabels := "        Now" + strings.Repeat(" ", chartWidth-20) + fmt.Sprintf("+%dm", snapshot.Snapshot.HorizonSeconds/60)
	lines = append(lines, xLabels)

	return strings.Join(lines, "\n")
}

// renderEmpty renders an empty chart placeholder.
func (c *QuantileChart) renderEmpty() string {
	return lipgloss.NewStyle().
		Foreground(lipgloss.Color("241")).
		Render("No forecast data available")
}

// findMinMaxQuantile finds min and max values in a slice (quantile version).
func findMinMaxQuantile(values []float64) (float64, float64) {
	if len(values) == 0 {
		return 0, 0
	}

	minVal := values[0]
	maxVal := values[0]

	for _, v := range values {
		if v < minVal {
			minVal = v
		}
		if v > maxVal {
			maxVal = v
		}
	}

	return minVal, maxVal
}

// findMinMaxAcross finds min and max across multiple quantile slices.
func findMinMaxAcross(slices ...[]float64) (float64, float64) {
	minVal := math.MaxFloat64
	maxVal := -math.MaxFloat64

	for _, slice := range slices {
		if len(slice) == 0 {
			continue
		}
		for _, v := range slice {
			if v < minVal {
				minVal = v
			}
			if v > maxVal {
				maxVal = v
			}
		}
	}

	if minVal == math.MaxFloat64 {
		return 0, 0
	}

	return minVal, maxVal
}
