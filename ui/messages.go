// Package ui implements the Bubble Tea UI components.
package ui

import (
	"time"

	"github.com/HatiCode/kedastral-tui/client"
)

type tickMsg time.Time

type snapshotMsg struct {
	data *client.SnapshotData
	err  error
}

type scalerMetricsMsg struct {
	data *client.ScalerMetrics
	err  error
}

type healthMsg struct {
	forecasterHealthy bool
	scalerHealthy     bool
}
