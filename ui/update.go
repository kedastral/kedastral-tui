package ui

import (
	"context"
	"fmt"
	"time"

	"github.com/HatiCode/kedastral-tui/client"
	"github.com/HatiCode/kedastral-tui/components"
	"github.com/HatiCode/kedastral-tui/config"
	"github.com/HatiCode/kedastral-tui/ui/panels"
	"github.com/HatiCode/kedastral-tui/ui/theme"
	tea "github.com/charmbracelet/bubbletea"
)

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "tab", "shift+tab", "w", "m", "[", "]":
			return m.handleFocusSwitch(msg)
		}

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
				m.loading = true
				return m, fetchData(m.client, m.currentWorkload, m.cfg.LeadTime)
			}
		case "escape", "esc":
			if m.showHelp {
				m.showHelp = false
			} else if m.err != nil {
				m.err = nil
			}
		case "c":
			if !m.showHelp {
				if err := m.copyCurrentTab(); err != nil {
					m.toastManager.Add(fmt.Sprintf("Copy failed: %v", err), components.ToastError, 3*time.Second)
				}
			}
		case "e":
			if !m.showHelp {
				if err := m.exportCurrentTab(); err != nil {
					m.toastManager.Add(fmt.Sprintf("Export failed: %v", err), components.ToastError, 3*time.Second)
				}
			}
		case "+", "=":
			if !m.showHelp {
				m.cfg.RefreshInterval = m.cfg.RefreshInterval + time.Second
				if m.cfg.RefreshInterval > 60*time.Second {
					m.cfg.RefreshInterval = 60 * time.Second
				}
				if err := config.SaveConfig(m.cfg); err == nil {
					m.toastManager.Add(fmt.Sprintf("Refresh interval: %s", m.cfg.RefreshInterval), components.ToastInfo, 2*time.Second)
				}
			}
		case "-", "_":
			if !m.showHelp {
				m.cfg.RefreshInterval = m.cfg.RefreshInterval - time.Second
				if m.cfg.RefreshInterval < time.Second {
					m.cfg.RefreshInterval = time.Second
				}
				if err := config.SaveConfig(m.cfg); err == nil {
					m.toastManager.Add(fmt.Sprintf("Refresh interval: %s", m.cfg.RefreshInterval), components.ToastInfo, 2*time.Second)
				}
			}
		case "ctrl+r":
			if !m.showHelp {
				m.loading = true
				m.err = nil
				m.toastManager.Add("Retrying...", components.ToastInfo, 1*time.Second)
				return m, fetchData(m.client, m.currentWorkload, m.cfg.LeadTime)
			}
		case "t":
			if !m.showHelp {
				if m.cfg.Theme == "dark" {
					m.cfg.Theme = "light"
					m.theme = theme.Light
				} else {
					m.cfg.Theme = "dark"
					m.theme = theme.Dark
				}
				if err := config.SaveConfig(m.cfg); err == nil {
					m.toastManager.Add(fmt.Sprintf("Theme: %s", m.cfg.Theme), components.ToastInfo, 2*time.Second)
				}
			}
		}

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.layoutMgr.SetTerminalSize(msg.Width, msg.Height)
		m.ready = true

		layout := m.layoutMgr.Compute()
		contentHeight := layout.Main.H - 10
		contentWidth := layout.Main.W - 4
		if contentHeight < 5 {
			contentHeight = 5
		}
		if contentWidth < 10 {
			contentWidth = 10
		}

		for tabID, vp := range m.tabViewports {
			vp.Width = contentWidth
			vp.Height = contentHeight
			m.tabViewports[tabID] = vp
		}

	case tickMsg:
		if m.mode == ModeLive {
			m.loading = true
			return m, tea.Batch(
				tick(m.cfg.RefreshInterval),
				fetchData(m.client, m.cfg.Workload, m.cfg.LeadTime),
			)
		}

	case snapshotMsg:
		m.loading = false
		m.lastUpdate = time.Now()
		if msg.err != nil {
			m.err = msg.err
		} else {
			m.snapshot = msg.data
			m.err = nil
		}

	case quantileSnapshotMsg:
		m.loading = false
		m.lastUpdate = time.Now()
		if msg.err != nil {
			m.err = msg.err
		} else {
			m.quantileSnapshot = msg.data
			m.apiVersion = msg.data.APIVersion
			m.err = nil

			if m.bottomPanel != nil {
				m.bottomPanel.UpdateAPIVersion(msg.data.APIVersion)
			}

			cmds = append(cmds, func() tea.Msg {
				return panels.NewLogMsg{
					Log: fmt.Sprintf("Forecast received, age: %.1fs", msg.data.ForecastAge.Seconds()),
				}
			})
		}

	case scalerMetricsMsg:
		if msg.err == nil {
			m.scalerMetrics = msg.data
			// Update bottom panel metrics
			if m.bottomPanel != nil {
				m.bottomPanel.UpdateMetrics(msg.data)
			}
		}

	case healthMsg:
		m.forecasterHealthy = msg.forecasterHealthy
		m.scalerHealthy = msg.scalerHealthy

	case workloadListMsg:
		if msg.err == nil && len(msg.workloads) > 0 {
			m.workloads = msg.workloads
			layout := m.layoutMgr.Compute()
			sidebar := panels.NewSidebar(msg.workloads, layout.Sidebar.W, layout.Sidebar.H)
			m.sidebar = &sidebar
		}

	case panels.WorkloadSelectedMsg:
		m.currentWorkload = msg.Workload
		m.loading = true
		return m, fetchData(m.client, msg.Workload, m.cfg.LeadTime)

	case panels.TabSwitchMsg:
		m.activeTab = TabID(msg.TabID)
	}

	var spinnerCmd tea.Cmd
	m.spinner, spinnerCmd = m.spinner.Update(msg)
	if spinnerCmd != nil {
		cmds = append(cmds, spinnerCmd)
	}

	if m.sidebar != nil && m.focusedPanel == PanelSidebar {
		var cmd tea.Cmd
		*m.sidebar, cmd = m.sidebar.Update(msg)
		if cmd != nil {
			cmds = append(cmds, cmd)
		}
	}

	if m.mainTabs != nil && m.focusedPanel == PanelMain {
		var cmd tea.Cmd
		*m.mainTabs, cmd = m.mainTabs.Update(msg)
		if cmd != nil {
			cmds = append(cmds, cmd)
		}

		if vp, ok := m.tabViewports[m.activeTab]; ok {
			if keyMsg, ok := msg.(tea.KeyMsg); ok {
				switch keyMsg.String() {
				case "j", "down", "k", "up", "g", "G", "ctrl+d", "ctrl+u", "pgdown", "pgup":
					vp, cmd = vp.Update(msg)
					m.tabViewports[m.activeTab] = vp
					if cmd != nil {
						cmds = append(cmds, cmd)
					}
				}
			}
		}
	}

	if m.bottomPanel != nil && m.focusedPanel == PanelBottom {
		var cmd tea.Cmd
		*m.bottomPanel, cmd = m.bottomPanel.Update(msg)
		if cmd != nil {
			cmds = append(cmds, cmd)
		}
	}

	if m.bottomPanel != nil {
		if _, ok := msg.(panels.NewLogMsg); ok {
			var cmd tea.Cmd
			*m.bottomPanel, cmd = m.bottomPanel.Update(msg)
			if cmd != nil {
				cmds = append(cmds, cmd)
			}
		}
	}

	if len(cmds) > 0 {
		return m, tea.Batch(cmds...)
	}
	return m, nil
}

func fetchData(c *client.Client, workload string, leadTime time.Duration) tea.Cmd {
	return func() tea.Msg {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		quantileSnapshotCh := make(chan quantileSnapshotMsg, 1)
		metricsCh := make(chan scalerMetricsMsg, 1)

		go func() {
			data, err := c.GetQuantileSnapshot(ctx, workload, leadTime)
			quantileSnapshotCh <- quantileSnapshotMsg{data: data, err: err}
		}()

		go func() {
			data, err := c.GetScalerMetrics(ctx)
			metricsCh <- scalerMetricsMsg{data: data, err: err}
		}()

		quantileSnapshot := <-quantileSnapshotCh
		metrics := <-metricsCh

		forecasterHealthy, scalerHealthy := c.GetHealthStatus(ctx)

		tea.Batch(
			func() tea.Msg { return quantileSnapshot },
			func() tea.Msg { return metrics },
			func() tea.Msg { return healthMsg{forecasterHealthy, scalerHealthy} },
		)

		return quantileSnapshot
	}
}
