package ui

import (
	"fmt"
	"strings"
	"time"

	"github.com/HatiCode/kedastral-tui/components"
	"github.com/atotto/clipboard"
)

func (m *Model) copyCurrentTab() error {
	var content string

	switch m.activeTab {
	case TabCharts:
		if m.quantileSnapshot != nil {
			content = fmt.Sprintf("Workload: %s\n", m.currentWorkload)
			content += fmt.Sprintf("Forecast Age: %s\n", m.quantileSnapshot.ForecastAge.Round(time.Second))
			content += fmt.Sprintf("Generated At: %s\n\n", m.quantileSnapshot.Snapshot.GeneratedAt.Format(time.RFC3339))
			content += "Quantile Forecast Data:\n"
			content += fmt.Sprintf("Metric: %s\n", m.quantileSnapshot.Snapshot.Metric)
			content += fmt.Sprintf("Step: %ds\n", m.quantileSnapshot.Snapshot.StepSeconds)
			content += fmt.Sprintf("Horizon: %ds\n", m.quantileSnapshot.Snapshot.HorizonSeconds)
		} else if m.snapshot != nil {
			content = fmt.Sprintf("Workload: %s\n", m.currentWorkload)
			content += fmt.Sprintf("Forecast Age: %s\n", m.snapshot.ForecastAge.Round(time.Second))
			content += fmt.Sprintf("Generated At: %s\n", m.snapshot.Snapshot.GeneratedAt.Format(time.RFC3339))
		} else {
			content = "No forecast data available"
		}

	case TabTables:
		if m.snapshot != nil {
			snap := m.snapshot.Snapshot
			stepDuration := time.Duration(snap.StepSeconds) * time.Second
			content = "Replica Scaling Data:\n\n"
			content += "Time Offset\tForecast\tDesired Replicas\n"
			for i := 0; i < len(snap.DesiredReplicas); i++ {
				timeOffset := stepDuration * time.Duration(i)
				forecastVal := 0.0
				if i < len(snap.Values) {
					forecastVal = snap.Values[i]
				}
				content += fmt.Sprintf("+%s\t%.2f\t%d\n",
					timeOffset,
					forecastVal,
					snap.DesiredReplicas[i],
				)
			}
		} else {
			content = "No table data available"
		}

	case TabConfig:
		var lines []string
		lines = append(lines, "Workload Configuration:")
		lines = append(lines, fmt.Sprintf("  Workload: %s", m.cfg.Workload))
		if m.snapshot != nil {
			lines = append(lines, fmt.Sprintf("  Metric: %s", m.snapshot.Snapshot.Metric))
			lines = append(lines, fmt.Sprintf("  Step Duration: %ds", m.snapshot.Snapshot.StepSeconds))
			lines = append(lines, fmt.Sprintf("  Horizon: %ds", m.snapshot.Snapshot.HorizonSeconds))
		}
		lines = append(lines, "")
		lines = append(lines, "TUI Configuration:")
		lines = append(lines, fmt.Sprintf("  Forecaster URL: %s", m.cfg.ForecasterURL))
		lines = append(lines, fmt.Sprintf("  Scaler URL: %s", m.cfg.ScalerURL))
		lines = append(lines, fmt.Sprintf("  Refresh Interval: %s", m.cfg.RefreshInterval))
		lines = append(lines, fmt.Sprintf("  Lead Time: %s", m.cfg.LeadTime))
		content = strings.Join(lines, "\n")

	case TabLogs:
		if m.bottomPanel != nil {
			content = "Logs:\n\n(Log data from bottom panel)"
		} else {
			content = "No logs available"
		}

	default:
		return fmt.Errorf("unknown tab")
	}

	if err := clipboard.WriteAll(content); err != nil {
		return fmt.Errorf("failed to copy to clipboard: %w", err)
	}

	m.toastManager.Add("✓ Copied to clipboard", components.ToastSuccess, 2*time.Second)
	return nil
}
