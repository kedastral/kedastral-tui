package ui

import (
	"context"

	"github.com/HatiCode/kedastral-tui/client"
	tea "github.com/charmbracelet/bubbletea"
	"time"
)

// Update handles messages and updates the model.
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit
		case "h":
			m.showHelp = !m.showHelp
		case " ":
			if m.showHelp {
				m.showHelp = false
			} else if m.mode == ModeLive {
				m.mode = ModePaused
			} else {
				m.mode = ModeLive
				return m, tick(m.cfg.RefreshInterval)
			}
		case "r":
			if !m.showHelp {
				return m, fetchData(m.client, m.cfg.Workload, m.cfg.LeadTime)
			}
		}

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.ready = true

	case tickMsg:
		if m.mode == ModeLive {
			return m, tea.Batch(
				tick(m.cfg.RefreshInterval),
				fetchData(m.client, m.cfg.Workload, m.cfg.LeadTime),
			)
		}

	case snapshotMsg:
		m.lastUpdate = time.Now()
		if msg.err != nil {
			m.err = msg.err
		} else {
			m.snapshot = msg.data
			m.err = nil
		}

	case scalerMetricsMsg:
		if msg.err == nil {
			m.scalerMetrics = msg.data
		}

	case healthMsg:
		m.forecasterHealthy = msg.forecasterHealthy
		m.scalerHealthy = msg.scalerHealthy
	}

	return m, nil
}

func fetchData(c *client.Client, workload string, leadTime time.Duration) tea.Cmd {
	return func() tea.Msg {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		snapshotCh := make(chan snapshotMsg, 1)
		metricsCh := make(chan scalerMetricsMsg, 1)

		go func() {
			data, err := c.GetSnapshot(ctx, workload, leadTime)
			snapshotCh <- snapshotMsg{data: data, err: err}
		}()

		go func() {
			data, err := c.GetScalerMetrics(ctx)
			metricsCh <- scalerMetricsMsg{data: data, err: err}
		}()

		snapshot := <-snapshotCh
		metrics := <-metricsCh

		forecasterHealthy, scalerHealthy := c.GetHealthStatus(ctx)

		tea.Batch(
			func() tea.Msg { return snapshot },
			func() tea.Msg { return metrics },
			func() tea.Msg { return healthMsg{forecasterHealthy, scalerHealthy} },
		)

		return snapshot
	}
}
