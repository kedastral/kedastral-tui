// Package components provides UI components for the TUI.
package components

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// Help renders the help screen.
type Help struct {
	width int
}

// NewHelp creates a new help component.
func NewHelp(width int) *Help {
	return &Help{width: width}
}

// Render renders the help screen.
func (h *Help) Render() string {
	var s strings.Builder

	titleStyle := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("39"))
	keyStyle := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("42"))
	descStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("252"))

	s.WriteString(titleStyle.Render("KEDASTRAL TUI - HELP"))
	s.WriteString("\n\n")

	s.WriteString(titleStyle.Render("KEYBOARD SHORTCUTS"))
	s.WriteString("\n\n")

	shortcuts := []struct {
		key  string
		desc string
	}{
		{"Tab", "Switch panel focus (Sidebar → Main → Bottom)"},
		{"Shift+Tab", "Switch panel focus (reverse)"},
		{"W", "Jump to sidebar (workload list)"},
		{"M", "Jump to main panel"},
		{"", ""},
		{"1-4", "Jump to tab (Charts/Tables/Config/Logs)"},
		{"H, L or ←, →", "Navigate tabs left/right"},
		{"J, K or ↑, ↓", "Scroll content up/down"},
		{"G", "Jump to top of scrollable content"},
		{"Shift+G", "Jump to bottom of scrollable content"},
		{"Ctrl+D/U", "Scroll half page down/up"},
		{"", ""},
		{"SPACE", "Toggle between live and paused modes"},
		{"R", "Manual refresh (fetch latest data)"},
		{"Ctrl+R", "Retry last failed request"},
		{"", ""},
		{"C", "Copy current tab content to clipboard"},
		{"E", "Export current tab content to file"},
		{"", ""},
		{"+/=", "Increase refresh interval (slower)"},
		{"-/_", "Decrease refresh interval (faster)"},
		{"T", "Toggle theme (dark/light)"},
		{"", ""},
		{"[", "Toggle sidebar collapse"},
		{"]", "Toggle bottom panel collapse"},
		{"B", "Cycle bottom panel mode (Logs/Metrics/Events/Info)"},
		{"", ""},
		{"/", "Filter workloads in sidebar"},
		{"Enter", "Select workload (in sidebar)"},
		{"Esc", "Clear error / Close help / Exit filter"},
		{"", ""},
		{"H", "Toggle this help screen"},
		{"Q or Ctrl+C", "Quit the application"},
	}

	for _, sc := range shortcuts {
		s.WriteString("  ")
		s.WriteString(keyStyle.Render(sc.key))
		s.WriteString(strings.Repeat(" ", 12-len(sc.key)))
		s.WriteString(descStyle.Render(sc.desc))
		s.WriteString("\n")
	}

	s.WriteString("\n")
	s.WriteString(titleStyle.Render("DISPLAY PANELS"))
	s.WriteString("\n\n")

	panels := []struct {
		name string
		desc string
	}{
		{"Sidebar", "Interactive workload list with health indicators"},
		{"Main Panel - Charts", "Quantile forecast visualization (P10/P50/P90)"},
		{"Main Panel - Tables", "Replica scaling decisions with lead time"},
		{"Main Panel - Config", "Workload and scaler configuration details"},
		{"Main Panel - Logs", "Forecast and scaler event logs"},
		{"Bottom Panel", "Logs, metrics, events, and system info (press B to cycle)"},
	}

	for _, p := range panels {
		s.WriteString("  ")
		s.WriteString(keyStyle.Render(p.name))
		s.WriteString("\n    ")
		s.WriteString(descStyle.Render(p.desc))
		s.WriteString("\n\n")
	}

	s.WriteString(titleStyle.Render("CONFIGURATION"))
	s.WriteString("\n\n")
	s.WriteString(descStyle.Render("Config file: ~/.config/kedastral-tui/config.json"))
	s.WriteString("\n")
	s.WriteString(descStyle.Render("Override with flags: --forecaster-url, --scaler-url, --workload"))
	s.WriteString("\n\n")

	s.WriteString(lipgloss.NewStyle().Foreground(lipgloss.Color("241")).Render("Press H or any key to close this help screen"))

	return s.String()
}
