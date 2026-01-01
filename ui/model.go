// Package ui implements the Bubble Tea UI components.
package ui

import (
	"context"
	"time"

	"github.com/HatiCode/kedastral-tui/client"
	"github.com/HatiCode/kedastral-tui/components"
	"github.com/HatiCode/kedastral-tui/config"
	"github.com/HatiCode/kedastral-tui/ui/layout"
	"github.com/HatiCode/kedastral-tui/ui/panels"
	"github.com/HatiCode/kedastral-tui/ui/theme"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
)

type Mode int

const (
	ModeLive Mode = iota
	ModePaused
)

// PanelID identifies which panel has focus.
type PanelID int

const (
	PanelSidebar PanelID = iota
	PanelMain
	PanelBottom
)

// TabID identifies which tab is active in the main panel.
type TabID int

const (
	TabCharts TabID = iota
	TabTables
	TabConfig
	TabLogs
)

type Model struct {
	cfg               *config.Config
	client            *client.Client
	mode              Mode
	snapshot          *client.SnapshotData
	quantileSnapshot  *client.QuantileSnapshotData
	scalerMetrics     *client.ScalerMetrics
	forecasterHealthy bool
	scalerHealthy     bool
	lastUpdate        time.Time
	err               error
	width             int
	height            int
	ready             bool
	showHelp          bool
	loading           bool
	spinner           components.LoadingSpinner
	toastManager      *components.ToastManager
	statusBar         any
	forecastChart     any
	replicaTable      any

	focusedPanel PanelID
	layoutMgr    *layout.LayoutManager

	workloads       []client.WorkloadInfo
	currentWorkload string
	apiVersion      int

	activeTab TabID

	sidebar     *panels.SidebarModel
	mainTabs    *panels.TabBarModel
	bottomPanel *panels.BottomPanelModel

	tabViewports map[TabID]viewport.Model
	theme        *theme.Theme
}

func NewModel(cfg *config.Config, c *client.Client) Model {
	tabBar := panels.NewTabBar(100)
	bottomPanel := panels.NewBottomPanel(100, 10, cfg)
	spinner := components.NewLoadingSpinner()
	toastManager := components.NewToastManager(100)
	currentTheme := theme.Get(cfg.Theme)

	tabViewports := make(map[TabID]viewport.Model)
	for _, tabID := range []TabID{TabCharts, TabTables, TabConfig, TabLogs} {
		vp := viewport.New(100, 20)
		tabViewports[tabID] = vp
	}

	return Model{
		cfg:             cfg,
		client:          c,
		mode:            ModeLive,
		focusedPanel:    PanelMain,
		layoutMgr:       layout.NewLayoutManager(),
		currentWorkload: cfg.Workload,
		activeTab:       TabCharts,
		mainTabs:        &tabBar,
		bottomPanel:     &bottomPanel,
		spinner:         spinner,
		toastManager:    toastManager,
		loading:         false,
		tabViewports:    tabViewports,
		theme:           currentTheme,
	}
}

func (m Model) Init() tea.Cmd {
	return tea.Batch(
		m.spinner.Init(),
		tick(m.cfg.RefreshInterval),
		fetchWorkloadList(m.client),
		fetchData(m.client, m.cfg.Workload, m.cfg.LeadTime),
	)
}

func fetchWorkloadList(c *client.Client) tea.Cmd {
	return func() tea.Msg {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		workloads, err := c.GetWorkloads(ctx)
		return workloadListMsg{workloads: workloads, err: err}
	}
}

func tick(d time.Duration) tea.Cmd {
	return tea.Tick(d, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}
