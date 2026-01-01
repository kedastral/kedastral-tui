package ui

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/HatiCode/kedastral-tui/config"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type setupStep int

const (
	stepForecasterURL setupStep = iota
	stepScalerURL
	stepWorkload
	stepComplete
)

type SetupModel struct {
	step          setupStep
	forecasterURL string
	scalerURL     string
	workload      string
	input         string
	err           error
	width         int
	height        int
}

func NewSetupModel() SetupModel {
	return SetupModel{
		step:      stepForecasterURL,
		scalerURL: "http://localhost:8082",
	}
}

func (m SetupModel) Init() tea.Cmd {
	return nil
}

func (m SetupModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit

		case "enter":
			switch m.step {
			case stepForecasterURL:
				if m.input != "" {
					m.forecasterURL = strings.TrimSpace(m.input)
					m.input = m.scalerURL
					m.step = stepScalerURL
				}
			case stepScalerURL:
				if m.input != "" {
					m.scalerURL = strings.TrimSpace(m.input)
					m.input = ""
					m.step = stepWorkload
				}
			case stepWorkload:
				if m.input != "" {
					m.workload = strings.TrimSpace(m.input)
					if err := m.saveConfig(); err != nil {
						m.err = err
					} else {
						m.step = stepComplete
						return m, tea.Quit
					}
				}
			}

		case "backspace":
			if len(m.input) > 0 {
				m.input = m.input[:len(m.input)-1]
			}

		default:
			if len(msg.String()) == 1 {
				m.input += msg.String()
			}
		}

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
	}

	return m, nil
}

func (m SetupModel) View() string {
	var s strings.Builder

	titleStyle := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("39"))
	promptStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("42"))
	inputStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("220"))
	helpStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("241"))

	s.WriteString(titleStyle.Render("Kedastral TUI - First Time Setup"))
	s.WriteString("\n\n")

	if m.err != nil {
		errorStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("196")).Bold(true)
		s.WriteString(errorStyle.Render(fmt.Sprintf("Error: %v\n\n", m.err)))
	}

	switch m.step {
	case stepForecasterURL:
		s.WriteString(promptStyle.Render("Forecaster URL"))
		s.WriteString("\n")
		s.WriteString("Enter the HTTP URL for the kedastral forecaster service.\n")
		s.WriteString("Example: http://localhost:8081\n\n")
		s.WriteString("> ")
		s.WriteString(inputStyle.Render(m.input))
		s.WriteString("_")

	case stepScalerURL:
		s.WriteString(promptStyle.Render("Scaler URL"))
		s.WriteString("\n")
		s.WriteString("Enter the HTTP URL for the kedastral scaler service.\n")
		s.WriteString(helpStyle.Render("(Press ENTER to use default: http://localhost:8082)\n\n"))
		s.WriteString("> ")
		s.WriteString(inputStyle.Render(m.input))
		s.WriteString("_")

	case stepWorkload:
		s.WriteString(promptStyle.Render("Workload Name"))
		s.WriteString("\n")
		s.WriteString("Enter the name of the workload to monitor.\n")
		s.WriteString("Example: test-app\n\n")
		s.WriteString("> ")
		s.WriteString(inputStyle.Render(m.input))
		s.WriteString("_")
	}

	s.WriteString("\n\n")
	s.WriteString(helpStyle.Render("Press ENTER to continue, Ctrl+C to cancel"))

	return s.String()
}

func (m *SetupModel) saveConfig() error {
	cfg := &config.Config{
		ForecasterURL: m.forecasterURL,
		ScalerURL:     m.scalerURL,
		Workload:      m.workload,
	}

	homeDir, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("failed to get home directory: %w", err)
	}

	configDir := filepath.Join(homeDir, ".config", "kedastral-tui")
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}

	configPath := filepath.Join(configDir, "config.json")

	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}

	if err := os.WriteFile(configPath, data, 0644); err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}

	return nil
}
