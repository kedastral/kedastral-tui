// Package components provides UI components for the TUI.
package components

import (
	"fmt"
	"strings"
	"time"

	"github.com/charmbracelet/lipgloss"
	"github.com/HatiCode/kedastral-tui/client"
)

// ReplicaTable renders a table of replica scaling decisions.
type ReplicaTable struct {
	width int
}

// NewReplicaTable creates a new replica table component.
func NewReplicaTable(width int) *ReplicaTable {
	return &ReplicaTable{width: width}
}

// Render renders the replica table.
func (r *ReplicaTable) Render(snapshot *client.SnapshotData) string {
	if snapshot == nil || len(snapshot.Snapshot.DesiredReplicas) == 0 {
		return "No replica data available"
	}

	var s strings.Builder

	headerStyle := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("39"))
	selectedStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("42")).Bold(true)
	normalStyle := lipgloss.NewStyle()

	s.WriteString(headerStyle.Render("REPLICA SCALING DECISIONS"))
	s.WriteString("\n\n")

	s.WriteString(fmt.Sprintf("%-10s  %-12s  %-8s  %-10s\n",
		"Time", "Forecast", "Desired", "Selected"))
	s.WriteString(strings.Repeat("─", r.width-4))
	s.WriteString("\n")

	snap := snapshot.Snapshot
	stepDuration := time.Duration(snap.StepSeconds) * time.Second

	maxRows := 10
	if len(snap.DesiredReplicas) < maxRows {
		maxRows = len(snap.DesiredReplicas)
	}

	for i := 0; i < maxRows; i++ {
		timeOffset := stepDuration * time.Duration(i)
		timeStr := formatTimeOffset(timeOffset)

		forecastVal := 0.0
		if i < len(snap.Values) {
			forecastVal = snap.Values[i]
		}

		replicas := snap.DesiredReplicas[i]

		isSelected := i == snapshot.LeadTimeIndex

		line := fmt.Sprintf("%-10s  %-12s  %-8d  ",
			timeStr,
			fmt.Sprintf("%.1f %s", forecastVal, snap.Metric),
			replicas,
		)

		if isSelected {
			line += "← SELECTED"
			s.WriteString(selectedStyle.Render(line))
		} else {
			s.WriteString(normalStyle.Render(line))
		}
		s.WriteString("\n")
	}

	if len(snap.DesiredReplicas) > maxRows {
		s.WriteString(fmt.Sprintf("\n... and %d more steps", len(snap.DesiredReplicas)-maxRows))
	}

	return s.String()
}

func formatTimeOffset(d time.Duration) string {
	if d == 0 {
		return "Now"
	}

	minutes := int(d.Minutes())
	if minutes < 1 {
		return fmt.Sprintf("+%ds", int(d.Seconds()))
	}
	return fmt.Sprintf("+%dm", minutes)
}
