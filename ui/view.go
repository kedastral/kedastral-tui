package ui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/HatiCode/kedastral-tui/client"
	"github.com/HatiCode/kedastral-tui/components"
)

var (
	mutedStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("241"))
)

// View renders the UI.
func (m Model) View() string {
	if !m.ready {
		return "Initializing..."
	}

	if m.showHelp {
		help := components.NewHelp(m.width)
		return help.Render()
	}

	var s strings.Builder

	modeStr := "LIVE"
	if m.mode == ModePaused {
		modeStr = "PAUSED"
	}

	statusBar := components.NewStatusBar(m.width)
	s.WriteString(statusBar.Render(
		m.cfg.Workload,
		modeStr,
		m.lastUpdate,
		m.snapshot,
		m.scalerMetrics,
		m.forecasterHealthy,
		m.scalerHealthy,
		m.err,
	))
	s.WriteString("\n")

	chartHeight := 10
	if m.height > 30 {
		chartHeight = 12
	}

	forecastChart := components.NewForecastChart(m.width, chartHeight)
	s.WriteString(forecastChart.Render(m.snapshot))
	s.WriteString("\n")

	replicaTable := components.NewReplicaTable(m.width)
	s.WriteString(replicaTable.Render(m.snapshot))
	s.WriteString("\n\n")

	if m.scalerMetrics != nil {
		s.WriteString(renderScalerInfo(m.scalerMetrics))
		s.WriteString("\n\n")
	}

	s.WriteString(mutedStyle.Render("[SPACE] pause  [R] refresh  [H] help  [Q] quit"))

	return s.String()
}

func renderScalerInfo(metrics *client.ScalerMetrics) string {
	var s strings.Builder

	titleStyle := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("39"))
	successStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("42"))
	errorStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("196"))

	s.WriteString(titleStyle.Render("SCALER STATUS"))
	s.WriteString("\n")

	status := "Active: "
	if metrics.Active {
		status += successStyle.Render("✓")
	} else {
		status += errorStyle.Render("✗")
	}
	s.WriteString(status)
	s.WriteString(fmt.Sprintf("  Desired replicas: %d", metrics.DesiredReplicas))
	s.WriteString(fmt.Sprintf("  Forecast age seen: %.1fs", metrics.ForecastAgeSeen))

	return s.String()
}
