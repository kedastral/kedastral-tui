// Package ui implements the Bubble Tea UI components.
package ui

import (
	"time"

	"github.com/HatiCode/kedastral-tui/client"
	"github.com/HatiCode/kedastral-tui/config"
	tea "github.com/charmbracelet/bubbletea"
)

type Mode int

const (
	ModeLive Mode = iota
	ModePaused
)

// Model holds the UI state.
type Model struct {
	cfg               *config.Config
	client            *client.Client
	mode              Mode
	snapshot          *client.SnapshotData
	scalerMetrics     *client.ScalerMetrics
	forecasterHealthy bool
	scalerHealthy     bool
	lastUpdate        time.Time
	err               error
	width             int
	height            int
	ready             bool
	showHelp          bool
	statusBar         any
	forecastChart     any
	replicaTable      any
}

// NewModel creates a new TUI model.
func NewModel(cfg *config.Config, c *client.Client) Model {
	return Model{
		cfg:    cfg,
		client: c,
		mode:   ModeLive,
	}
}

// Init initializes the model.
func (m Model) Init() tea.Cmd {
	return tea.Batch(
		tick(m.cfg.RefreshInterval),
		fetchData(m.client, m.cfg.Workload, m.cfg.LeadTime),
	)
}

func tick(d time.Duration) tea.Cmd {
	return tea.Tick(d, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}
