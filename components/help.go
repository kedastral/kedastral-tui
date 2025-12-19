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
		{"SPACE", "Toggle between live and paused modes"},
		{"R", "Manual refresh (fetch latest data)"},
		{"H", "Toggle this help screen"},
		{"Q", "Quit the application"},
		{"Ctrl+C", "Force quit"},
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
		{"Status Bar", "Shows workload, mode (LIVE/PAUSED), connection health"},
		{"Forecast Chart", "ASCII visualization of predicted metric values"},
		{"Replica Table", "Scaling decisions with lead time selection"},
		{"Scaler Status", "Current scaler state and desired replicas"},
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
