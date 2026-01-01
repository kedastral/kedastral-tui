package panels

import (
	"fmt"
	"strings"
	"time"

	"github.com/HatiCode/kedastral-tui/client"
	"github.com/HatiCode/kedastral-tui/config"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type BottomPanelMode int

const (
	BottomLogs BottomPanelMode = iota
	BottomMetrics
	BottomEvents
	BottomInfo
)

type NewLogMsg struct {
	Log string
}

type BottomPanelModel struct {
	mode       BottomPanelMode
	viewport   viewport.Model
	width      int
	height     int
	logs       []string
	metrics    *client.ScalerMetrics
	cfg        *config.Config
	apiVersion int
}

func NewBottomPanel(width, height int, cfg *config.Config) BottomPanelModel {
	vp := viewport.New(width-4, height-4)
	vp.SetContent("No logs yet")

	return BottomPanelModel{
		mode:     BottomLogs,
		viewport: vp,
		width:    width,
		height:   height,
		logs:     make([]string, 0, 1000),
		cfg:      cfg,
	}
}

func (b BottomPanelModel) Update(msg tea.Msg) (BottomPanelModel, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "b":
			// Cycle through modes
			b.mode = (b.mode + 1) % 4
			b.updateViewportContent()
			return b, nil
		case "j", "down":
			b.viewport, cmd = b.viewport.Update(msg)
			return b, cmd
		case "k", "up":
			b.viewport, cmd = b.viewport.Update(msg)
			return b, cmd
		}

	case NewLogMsg:
		b.addLog(msg.Log)
		if b.mode == BottomLogs {
			b.updateViewportContent()
		}
	}

	return b, nil
}

func (b BottomPanelModel) View() string {
	title := fmt.Sprintf("─ %s ", b.modeTitle())
	titleStyle := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("39"))

	return lipgloss.JoinVertical(lipgloss.Left,
		titleStyle.Render(title),
		b.viewport.View(),
	)
}

func (b *BottomPanelModel) SetSize(width, height int) {
	b.width = width
	b.height = height
	b.viewport.Width = width - 4
	b.viewport.Height = height - 4
}

func (b *BottomPanelModel) UpdateMetrics(metrics *client.ScalerMetrics) {
	b.metrics = metrics
	if b.mode == BottomMetrics {
		b.updateViewportContent()
	}
}

func (b *BottomPanelModel) UpdateAPIVersion(version int) {
	b.apiVersion = version
	if b.mode == BottomInfo {
		b.updateViewportContent()
	}
}

func (b *BottomPanelModel) addLog(log string) {
	timestamp := time.Now().Format("15:04:05")
	entry := fmt.Sprintf("[%s] %s", timestamp, log)

	b.logs = append(b.logs, entry)

	if len(b.logs) > 1000 {
		b.logs = b.logs[len(b.logs)-1000:]
	}
}

func (b *BottomPanelModel) updateViewportContent() {
	var content string

	switch b.mode {
	case BottomLogs:
		content = b.renderLogs()
	case BottomMetrics:
		content = b.renderMetrics()
	case BottomEvents:
		content = b.renderEvents()
	case BottomInfo:
		content = b.renderInfo()
	}

	b.viewport.SetContent(content)
}

func (b BottomPanelModel) modeTitle() string {
	switch b.mode {
	case BottomLogs:
		return "Logs"
	case BottomMetrics:
		return "Metrics"
	case BottomEvents:
		return "Events"
	case BottomInfo:
		return "Info"
	default:
		return "Unknown"
	}
}

func (b *BottomPanelModel) renderLogs() string {
	if len(b.logs) == 0 {
		return lipgloss.NewStyle().
			Foreground(lipgloss.Color("241")).
			Render("No logs yet. Logs will appear here as events occur.")
	}

	start := 0
	if len(b.logs) > 50 {
		start = len(b.logs) - 50
	}

	return strings.Join(b.logs[start:], "\n")
}

func (b *BottomPanelModel) renderMetrics() string {
	if b.metrics == nil {
		return "No metrics available"
	}

	var s strings.Builder
	titleStyle := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("39"))
	successStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("42"))
	errorStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("196"))

	s.WriteString(titleStyle.Render("Scaler Metrics"))
	s.WriteString("\n\n")

	status := "Active: "
	if b.metrics.Active {
		status += successStyle.Render("✓ Yes")
	} else {
		status += errorStyle.Render("✗ No")
	}
	s.WriteString(status + "\n")
	s.WriteString(fmt.Sprintf("Desired replicas: %d\n", b.metrics.DesiredReplicas))
	s.WriteString(fmt.Sprintf("Forecast age seen: %.1fs\n", b.metrics.ForecastAgeSeen))

	connection := "Connection: "
	if b.metrics.ConnectionHealthy {
		connection += successStyle.Render("✓ Healthy")
	} else {
		connection += errorStyle.Render("✗ Unhealthy")
	}
	s.WriteString(connection + "\n")

	return s.String()
}

func (b *BottomPanelModel) renderEvents() string {
	// Placeholder for events
	return `Recent Events:

This view will show important events such as:
- Forecast updates
- Scaling decisions
- API errors
- Configuration changes

Events will be implemented in future iterations.`
}

func (b *BottomPanelModel) renderInfo() string {
	var s strings.Builder
	titleStyle := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("39"))
	checkmark := lipgloss.NewStyle().Foreground(lipgloss.Color("42")).Render("✓")
	xmark := lipgloss.NewStyle().Foreground(lipgloss.Color("196")).Render("✗")

	s.WriteString(titleStyle.Render("System Information"))
	s.WriteString("\n\n")

	if b.cfg != nil {
		s.WriteString(fmt.Sprintf("Forecaster URL:  %s\n", b.cfg.ForecasterURL))
		s.WriteString(fmt.Sprintf("Scaler URL:      %s\n", b.cfg.ScalerURL))
		s.WriteString(fmt.Sprintf("API Version:     v%d\n", b.apiVersion))
		s.WriteString(fmt.Sprintf("Refresh Interval: %s\n", b.cfg.RefreshInterval))
		s.WriteString(fmt.Sprintf("Lead Time:       %s\n", b.cfg.LeadTime))
	}

	s.WriteString("\n")
	s.WriteString(titleStyle.Render("Features"))
	s.WriteString("\n\n")

	quantilesIcon := xmark
	if b.apiVersion >= 2 {
		quantilesIcon = checkmark
	}

	s.WriteString(fmt.Sprintf("%s Quantile forecasts (P10/P50/P90)\n", quantilesIcon))
	s.WriteString(fmt.Sprintf("%s Multi-workload switching\n", checkmark))
	s.WriteString(fmt.Sprintf("%s Tab-based views\n", checkmark))
	s.WriteString(fmt.Sprintf("%s Bottom panel (logs/metrics/events/info)\n", checkmark))
	s.WriteString(fmt.Sprintf("%s Interactive sidebar\n", checkmark))

	return s.String()
}
