// Package components provides UI components for the TUI.
package components

import (
	"fmt"
	"math"
	"strings"

	"github.com/HatiCode/kedastral-tui/client"
	"github.com/charmbracelet/lipgloss"
)

// ForecastChart renders an ASCII line chart of forecast values.
type ForecastChart struct {
	width  int
	height int
}

// NewForecastChart creates a new forecast chart component.
func NewForecastChart(width, height int) *ForecastChart {
	return &ForecastChart{
		width:  width,
		height: height,
	}
}

// Render renders the forecast chart.
func (c *ForecastChart) Render(snapshot *client.SnapshotData) string {
	if snapshot == nil || len(snapshot.Snapshot.Values) == 0 {
		return "No forecast data available"
	}

	var s strings.Builder

	headerStyle := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("39"))
	s.WriteString(headerStyle.Render("FORECAST TIMELINE"))
	s.WriteString(fmt.Sprintf(" (next %dm)\n\n", snapshot.Snapshot.HorizonSeconds/60))

	values := snapshot.Snapshot.Values
	metric := snapshot.Snapshot.Metric

	chartWidth := c.width - 15
	chartHeight := max(c.height, 5)
	if chartWidth < 20 {
		chartWidth = 20
	}

	minVal, maxVal := findMinMax(values)

	if maxVal == minVal {
		maxVal = minVal + 1
	}

	valueRange := maxVal - minVal
	s.WriteString(fmt.Sprintf("%-8s\n", metric))

	for row := chartHeight - 1; row >= 0; row-- {
		y := minVal + (float64(row)/float64(chartHeight-1))*valueRange

		s.WriteString(fmt.Sprintf("%6.0f ┤", y))

		for col := 0; col < chartWidth; col++ {
			dataIndex := int(float64(col) / float64(chartWidth-1) * float64(len(values)-1))
			if dataIndex >= len(values) {
				dataIndex = len(values) - 1
			}

			value := values[dataIndex]

			normalizedValue := (value - minVal) / valueRange
			normalizedY := float64(row) / float64(chartHeight-1)

			diff := math.Abs(normalizedValue - normalizedY)

			if diff < 0.05 {
				s.WriteString("●")
			} else if diff < 0.1 {
				s.WriteString("·")
			} else {
				s.WriteString(" ")
			}
		}
		s.WriteString("\n")
	}

	s.WriteString("       └")
	s.WriteString(strings.Repeat("─", chartWidth))
	s.WriteString("\n")

	horizonMin := snapshot.Snapshot.HorizonSeconds / 60
	s.WriteString(fmt.Sprintf("        Now%s+%dm\n",
		strings.Repeat(" ", chartWidth-10),
		horizonMin))

	return s.String()
}

func findMinMax(values []float64) (float64, float64) {
	if len(values) == 0 {
		return 0, 0
	}

	min := values[0]
	max := values[0]

	for _, v := range values {
		if v < min {
			min = v
		}
		if v > max {
			max = v
		}
	}

	margin := (max - min) * 0.1
	if margin == 0 {
		margin = max * 0.1
	}
	if margin == 0 {
		margin = 1
	}

	return min - margin, max + margin
}
