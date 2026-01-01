package panels

import (
	"fmt"
	"io"
	"time"

	"github.com/HatiCode/kedastral-tui/client"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type WorkloadSelectedMsg struct {
	Workload string
}

type workloadItem struct {
	info client.WorkloadInfo
}

func (w workloadItem) FilterValue() string { return w.info.Name }

type SidebarModel struct {
	list      list.Model
	workloads []client.WorkloadInfo
	width     int
	height    int
	focused   bool
}

func NewSidebar(workloads []client.WorkloadInfo, width, height int) SidebarModel {
	items := make([]list.Item, len(workloads))
	for i, w := range workloads {
		items[i] = workloadItem{info: w}
	}

	delegate := list.NewDefaultDelegate()
	delegate.Styles.SelectedTitle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("42")).
		Bold(true)
	delegate.Styles.SelectedDesc = lipgloss.NewStyle().
		Foreground(lipgloss.Color("42"))

	l := list.New(items, newWorkloadDelegate(), width-4, height-4)
	l.Title = "Workloads"
	l.SetShowStatusBar(false)
	l.SetFilteringEnabled(true)
	l.Styles.Title = lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("39")).
		Padding(0, 1)

	return SidebarModel{
		list:      l,
		workloads: workloads,
		width:     width,
		height:    height,
	}
}

func (s SidebarModel) Update(msg tea.Msg) (SidebarModel, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.Type == tea.KeyEnter && s.list.FilterState() != list.Filtering {
			if selected, ok := s.list.SelectedItem().(workloadItem); ok {
				return s, func() tea.Msg {
					return WorkloadSelectedMsg{Workload: selected.info.Name}
				}
			}
		}
	}

	var cmd tea.Cmd
	s.list, cmd = s.list.Update(msg)
	return s, cmd
}

func (s SidebarModel) View() string {
	return s.list.View()
}

func (s *SidebarModel) SetSize(width, height int) {
	s.width = width
	s.height = height
	s.list.SetSize(width-4, height-4)
}

func (s *SidebarModel) SetFocused(focused bool) {
	s.focused = focused
}

func (s *SidebarModel) UpdateWorkloads(workloads []client.WorkloadInfo) {
	s.workloads = workloads

	items := make([]list.Item, len(workloads))
	for i, w := range workloads {
		items[i] = workloadItem{info: w}
	}

	s.list.SetItems(items)
}

type workloadDelegate struct{}

func newWorkloadDelegate() workloadDelegate {
	return workloadDelegate{}
}

func (d workloadDelegate) Height() int                               { return 1 }
func (d workloadDelegate) Spacing() int                              { return 0 }
func (d workloadDelegate) Update(msg tea.Msg, m *list.Model) tea.Cmd { return nil }

func (d workloadDelegate) Render(w io.Writer, m list.Model, index int, item list.Item) {
	workload, ok := item.(workloadItem)
	if !ok {
		return
	}

	// Check if this item is selected
	isSelected := index == m.Index()

	// Base style
	nameStyle := lipgloss.NewStyle()
	if isSelected {
		nameStyle = nameStyle.Foreground(lipgloss.Color("42")).Bold(true)
	}

	// Health indicator
	healthStyle := lipgloss.NewStyle()
	health := "[✓]"
	if !workload.info.Healthy {
		health = "[!]"
		healthStyle = healthStyle.Foreground(lipgloss.Color("196"))
	} else {
		healthStyle = healthStyle.Foreground(lipgloss.Color("42"))
	}

	// Age indicator
	age := formatAge(workload.info.LastForecast)
	ageStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("241"))

	// Render the item
	prefix := "  "
	if isSelected {
		prefix = "> "
	}

	fmt.Fprintf(w, "%s%s %s %s\n",
		prefix,
		nameStyle.Render(workload.info.Name),
		ageStyle.Render(age),
		healthStyle.Render(health),
	)
}

func formatAge(t time.Time) string {
	if t.IsZero() {
		return "---"
	}

	age := time.Since(t)

	if age < time.Minute {
		return fmt.Sprintf("%ds", int(age.Seconds()))
	} else if age < time.Hour {
		return fmt.Sprintf("%dm", int(age.Minutes()))
	} else if age < 24*time.Hour {
		return fmt.Sprintf("%dh", int(age.Hours()))
	} else {
		return fmt.Sprintf("%dd", int(age.Hours()/24))
	}
}
