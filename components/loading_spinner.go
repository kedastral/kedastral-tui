package components

import (
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type LoadingSpinner struct {
	spinner spinner.Model
}

func NewLoadingSpinner() LoadingSpinner {
	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))
	return LoadingSpinner{spinner: s}
}

func (s LoadingSpinner) Init() tea.Cmd {
	return s.spinner.Tick
}

func (s LoadingSpinner) Update(msg tea.Msg) (LoadingSpinner, tea.Cmd) {
	var cmd tea.Cmd
	s.spinner, cmd = s.spinner.Update(msg)
	return s, cmd
}

func (s LoadingSpinner) View() string {
	return s.spinner.View()
}
