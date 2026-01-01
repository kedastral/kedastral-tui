package ui

import (
	"fmt"
	"strings"

	"github.com/HatiCode/kedastral-tui/client"
	"github.com/HatiCode/kedastral-tui/components"
	"github.com/charmbracelet/lipgloss"
)

var (
	mutedStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("241"))
)

func (m Model) View() string {
	if !m.ready {
		return "Initializing..."
	}

	// Check minimum terminal size
	if m.width < 80 || m.height < 24 {
		errorStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("196")).Bold(true)
		return errorStyle.Render(
			"Terminal too small.\n" +
				fmt.Sprintf("Current: %dx%d\n", m.width, m.height) +
				"Minimum required: 80x24\n\n" +
				"Please resize your terminal.",
		)
	}

	if m.showHelp {
		help := components.NewHelp(m.width)
		return help.Render()
	}

	// Compute layout dimensions
	layout := m.layoutMgr.Compute()

	// Render sidebar panel
	sidebarContent := m.renderSidebar(layout.Sidebar.W, layout.Sidebar.H)

	// Render main panel
	mainContent := m.renderMain(layout.Main.W, layout.Main.H)

	// Render bottom panel
	bottomContent := m.renderBottom(layout.Bottom.W, layout.Bottom.H)

	// Compose the layout
	var result string

	if layout.Sidebar.W > 0 {
		// Sidebar + Main vertically stacked, then Bottom below Main
		if layout.Bottom.H > 0 {
			// Three panel layout: sidebar | (main over bottom)
			rightColumn := lipgloss.JoinVertical(lipgloss.Left, mainContent, bottomContent)
			result = lipgloss.JoinHorizontal(lipgloss.Top, sidebarContent, rightColumn)
		} else {
			// Two panel layout: sidebar | main
			result = lipgloss.JoinHorizontal(lipgloss.Top, sidebarContent, mainContent)
		}
	} else {
		// No sidebar
		if layout.Bottom.H > 0 {
			// Two panel layout: main over bottom
			result = lipgloss.JoinVertical(lipgloss.Left, mainContent, bottomContent)
		} else {
			// Single panel layout: just main
			result = mainContent
		}
	}

	if m.toastManager != nil && m.toastManager.HasToasts() {
		toastContent := m.toastManager.Render()
		toastOverlay := lipgloss.Place(
			m.width, 3,
			lipgloss.Right, lipgloss.Top,
			toastContent,
		)
		lines := strings.Split(result, "\n")
		if len(lines) > 0 {
			lines[0] = toastOverlay
		}
		result = strings.Join(lines, "\n")
	}

	return result
}

func (m Model) renderSidebar(width, height int) string {
	if width == 0 {
		return ""
	}

	focused := m.focusedPanel == PanelSidebar

	borderStyle := lipgloss.NewStyle().
		Width(width).
		Height(height).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(borderColor(focused))

	var content string
	if m.sidebar != nil {
		// Use actual sidebar
		content = m.sidebar.View()
	} else {
		// Fallback: show current workload
		content = lipgloss.NewStyle().
			Width(width - 4).
			Render(
				lipgloss.NewStyle().Bold(true).Render("Workloads") + "\n\n" +
					"> " + m.currentWorkload + "\n" +
					"\n" +
					"Loading workloads...",
			)
	}

	return borderStyle.Render(content)
}

func (m Model) renderMain(width, height int) string {
	if width == 0 || height == 0 {
		return ""
	}

	focused := m.focusedPanel == PanelMain

	borderStyle := lipgloss.NewStyle().
		Width(width).
		Height(height).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(borderColor(focused))

	// Status bar
	modeStr := "LIVE"
	if m.mode == ModePaused {
		modeStr = "PAUSED"
	}

	statusBar := components.NewStatusBar(width - 4)
	statusBarContent := statusBar.Render(
		m.currentWorkload,
		modeStr,
		m.lastUpdate,
		m.snapshot,
		m.scalerMetrics,
		m.forecasterHealthy,
		m.scalerHealthy,
		m.loading,
		m.spinner.View(),
		m.err,
	)

	// Tab bar
	var tabBar string
	if m.mainTabs != nil {
		tabBar = m.mainTabs.View()
	} else {
		tabBar = lipgloss.NewStyle().Bold(true).Render(
			"[■ Charts] | Tables | Config | Logs",
		)
	}

	vp := m.tabViewports[m.activeTab]
	var tabContent string
	switch m.activeTab {
	case TabCharts:
		chartHeight := 10
		if height > 30 {
			chartHeight = 12
		}
		if m.quantileSnapshot != nil {
			quantileChart := components.NewQuantileChart(width-4, chartHeight)
			tabContent = quantileChart.Render(m.quantileSnapshot)
		} else {
			forecastChart := components.NewForecastChart(width-4, chartHeight)
			tabContent = forecastChart.Render(m.snapshot)
		}

	case TabTables:
		replicaTable := components.NewReplicaTable(width - 4)
		tabContent = replicaTable.Render(m.snapshot)

	case TabConfig:
		tabContent = m.renderConfigView(width - 4)

	case TabLogs:
		tabContent = "Logs view\n(Will show forecast/scaler logs)"

	default:
		tabContent = "Unknown tab"
	}

	vp.SetContent(tabContent)
	mainContent := vp.View()

	footer := mutedStyle.Render("[Tab] focus  [1-4] tabs  [h/l] navigate  [SPACE] pause  [H] help  [Q] quit")

	content := lipgloss.JoinVertical(lipgloss.Left,
		statusBarContent,
		tabBar,
		"",
		mainContent,
		"",
		footer,
	)

	// Ensure content fits in the bordered area
	content = lipgloss.NewStyle().
		Width(width - 4).
		Height(height - 4).
		Render(content)

	return borderStyle.Render(content)
}

func (m Model) renderBottom(width, height int) string {
	if width == 0 || height == 0 {
		return ""
	}

	focused := m.focusedPanel == PanelBottom

	borderStyle := lipgloss.NewStyle().
		Width(width).
		Height(height).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(borderColor(focused))

	var content string
	if m.bottomPanel != nil {
		// Use actual bottom panel
		content = m.bottomPanel.View()
	} else {
		// Fallback
		content = "Loading..."
	}

	// Ensure content fits in bordered area
	content = lipgloss.NewStyle().
		Width(width - 4).
		Height(height - 4).
		Render(content)

	return borderStyle.Render(content)
}

func borderColor(focused bool) lipgloss.Color {
	if focused {
		return lipgloss.Color("39") // Bright blue
	}
	return lipgloss.Color("241") // Muted gray
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

func (m Model) renderConfigView(width int) string {
	var s strings.Builder

	titleStyle := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("39"))

	s.WriteString(titleStyle.Render("Workload Configuration"))
	s.WriteString("\n\n")

	if m.snapshot != nil {
		s.WriteString(fmt.Sprintf("  Workload:        %s\n", m.cfg.Workload))
		s.WriteString(fmt.Sprintf("  Metric:          %s\n", m.snapshot.Snapshot.Metric))
		s.WriteString(fmt.Sprintf("  Step Duration:   %ds\n", m.snapshot.Snapshot.StepSeconds))
		s.WriteString(fmt.Sprintf("  Horizon:         %ds\n", m.snapshot.Snapshot.HorizonSeconds))
		s.WriteString(fmt.Sprintf("  Generated At:    %s\n", m.snapshot.Snapshot.GeneratedAt.Format("15:04:05")))
	}

	s.WriteString("\n")
	s.WriteString(titleStyle.Render("Scaler Configuration"))
	s.WriteString("\n\n")

	if m.scalerMetrics != nil {
		s.WriteString(renderScalerInfo(m.scalerMetrics))
	}

	s.WriteString("\n\n")
	s.WriteString(titleStyle.Render("TUI Configuration"))
	s.WriteString("\n\n")
	s.WriteString(fmt.Sprintf("  Forecaster URL:  %s\n", m.cfg.ForecasterURL))
	s.WriteString(fmt.Sprintf("  Scaler URL:      %s\n", m.cfg.ScalerURL))
	s.WriteString(fmt.Sprintf("  Refresh Interval: %s\n", m.cfg.RefreshInterval))
	s.WriteString(fmt.Sprintf("  Lead Time:       %s\n", m.cfg.LeadTime))

	return s.String()
}
