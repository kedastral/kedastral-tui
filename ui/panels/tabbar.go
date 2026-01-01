package panels

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type TabID int

const (
	TabCharts TabID = iota
	TabTables
	TabConfig
	TabLogs
)

type TabSwitchMsg struct {
	TabID TabID
}

type TabInfo struct {
	ID    TabID
	Title string
	Icon  string
}

type TabBarModel struct {
	tabs      []TabInfo
	activeIdx int
	width     int
}

var (
	activeTabStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("42")).
			Background(lipgloss.Color("235")).
			Padding(0, 2)

	inactiveTabStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("241")).
				Padding(0, 2)
)

func NewTabBar(width int) TabBarModel {
	return TabBarModel{
		tabs: []TabInfo{
			{ID: TabCharts, Title: "Charts", Icon: "■"},
			{ID: TabTables, Title: "Tables", Icon: "▤"},
			{ID: TabConfig, Title: "Config", Icon: "⚙"},
			{ID: TabLogs, Title: "Logs", Icon: "≡"},
		},
		width:     width,
		activeIdx: 0,
	}
}

func (t TabBarModel) Update(msg tea.Msg) (TabBarModel, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "1":
			t.activeIdx = 0
			return t, t.emitTabSwitch()
		case "2":
			t.activeIdx = 1
			return t, t.emitTabSwitch()
		case "3":
			t.activeIdx = 2
			return t, t.emitTabSwitch()
		case "4":
			t.activeIdx = 3
			return t, t.emitTabSwitch()
		case "h", "left":
			if t.activeIdx > 0 {
				t.activeIdx--
				return t, t.emitTabSwitch()
			}
		case "l", "right":
			if t.activeIdx < len(t.tabs)-1 {
				t.activeIdx++
				return t, t.emitTabSwitch()
			}
		}
	}
	return t, nil
}

func (t TabBarModel) View() string {
	var tabs []string
	for i, tab := range t.tabs {
		style := inactiveTabStyle
		if i == t.activeIdx {
			style = activeTabStyle
		}
		tabs = append(tabs, style.Render(fmt.Sprintf("%s %s", tab.Icon, tab.Title)))
	}

	return lipgloss.JoinHorizontal(lipgloss.Top, tabs...)
}

func (t TabBarModel) ActiveTab() TabID {
	if t.activeIdx >= 0 && t.activeIdx < len(t.tabs) {
		return t.tabs[t.activeIdx].ID
	}
	return TabCharts
}

func (t *TabBarModel) SetActiveTab(id TabID) {
	for i, tab := range t.tabs {
		if tab.ID == id {
			t.activeIdx = i
			return
		}
	}
}

func (t TabBarModel) emitTabSwitch() tea.Cmd {
	return func() tea.Msg {
		return TabSwitchMsg{TabID: t.ActiveTab()}
	}
}
