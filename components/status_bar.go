// Package components provides UI components for the TUI.
package components

import (
	"fmt"
	"strings"
	"time"

	"github.com/charmbracelet/lipgloss"
	"github.com/HatiCode/kedastral-tui/client"
)

// StatusBar renders the header status bar.
type StatusBar struct {
	width int
}

// NewStatusBar creates a new status bar component.
func NewStatusBar(width int) *StatusBar {
	return &StatusBar{width: width}
}

// Render renders the status bar.
func (s *StatusBar) Render(
	workload string,
	mode string,
	lastUpdate time.Time,
	snapshot *client.SnapshotData,
	scalerMetrics *client.ScalerMetrics,
	forecasterHealthy, scalerHealthy bool,
	err error,
) string {
	var b strings.Builder

	titleStyle := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("39"))
	modeStyle := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("220"))
	successStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("42"))
	errorStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("196"))
	mutedStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("241"))

	b.WriteString(titleStyle.Render("Kedastral Monitor"))
	b.WriteString(" - ")
	b.WriteString(fmt.Sprintf("workload: %s", workload))

	b.WriteString("  ")
	b.WriteString(modeStyle.Render(fmt.Sprintf("[%s]", mode)))

	if !lastUpdate.IsZero() {
		b.WriteString("  ")
		b.WriteString(mutedStyle.Render(fmt.Sprintf("Last: %s ago", time.Since(lastUpdate).Round(time.Second))))
	}

	b.WriteString("\n")

	statusLine := "Status: "
	if forecasterHealthy {
		statusLine += successStyle.Render("Forecaster ✓")
	} else {
		statusLine += errorStyle.Render("Forecaster ✗")
	}

	statusLine += "  "
	if scalerHealthy {
		statusLine += successStyle.Render("Scaler ✓")
	} else {
		statusLine += errorStyle.Render("Scaler ✗")
	}

	if snapshot != nil {
		statusLine += fmt.Sprintf("  Forecast age: %s", snapshot.ForecastAge.Round(time.Second))
		if snapshot.Stale {
			statusLine += "  " + errorStyle.Render("[STALE]")
		}
	}

	b.WriteString(statusLine)
	b.WriteString("\n")

	if err != nil {
		b.WriteString(errorStyle.Render(fmt.Sprintf("Error: %v", err)))
		b.WriteString("\n")
	}

	b.WriteString(strings.Repeat("─", s.width))
	b.WriteString("\n")

	return b.String()
}
