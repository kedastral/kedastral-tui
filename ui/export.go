package ui

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/HatiCode/kedastral-tui/components"
)

func (m *Model) exportCurrentTab() error {
	timestamp := time.Now().Format("20060102-150405")
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("failed to get home directory: %w", err)
	}

	downloadsDir := filepath.Join(homeDir, "Downloads")
	if _, err := os.Stat(downloadsDir); os.IsNotExist(err) {
		downloadsDir = homeDir
	}

	var filename string
	var exportErr error

	switch m.activeTab {
	case TabCharts:
		if m.quantileSnapshot != nil {
			filename = filepath.Join(downloadsDir, fmt.Sprintf("kedastral-forecast-%s.json", timestamp))
			data, err := json.MarshalIndent(m.quantileSnapshot, "", "  ")
			if err != nil {
				return fmt.Errorf("failed to marshal forecast data: %w", err)
			}
			exportErr = os.WriteFile(filename, data, 0644)
		} else if m.snapshot != nil {
			filename = filepath.Join(downloadsDir, fmt.Sprintf("kedastral-forecast-%s.json", timestamp))
			data, err := json.MarshalIndent(m.snapshot, "", "  ")
			if err != nil {
				return fmt.Errorf("failed to marshal forecast data: %w", err)
			}
			exportErr = os.WriteFile(filename, data, 0644)
		} else {
			return fmt.Errorf("no forecast data to export")
		}

	case TabTables:
		if m.snapshot == nil {
			return fmt.Errorf("no table data to export")
		}

		filename = filepath.Join(downloadsDir, fmt.Sprintf("kedastral-replicas-%s.csv", timestamp))
		file, err := os.Create(filename)
		if err != nil {
			return fmt.Errorf("failed to create file: %w", err)
		}
		defer file.Close()

		writer := csv.NewWriter(file)
		defer writer.Flush()

		writer.Write([]string{"Time Offset (seconds)", "Forecast Value", "Desired Replicas"})

		snap := m.snapshot.Snapshot
		stepDuration := time.Duration(snap.StepSeconds) * time.Second
		for i := 0; i < len(snap.DesiredReplicas); i++ {
			timeOffset := stepDuration * time.Duration(i)
			forecastVal := 0.0
			if i < len(snap.Values) {
				forecastVal = snap.Values[i]
			}
			record := []string{
				fmt.Sprintf("%.0f", timeOffset.Seconds()),
				fmt.Sprintf("%.2f", forecastVal),
				fmt.Sprintf("%d", snap.DesiredReplicas[i]),
			}
			if err := writer.Write(record); err != nil {
				return fmt.Errorf("failed to write CSV record: %w", err)
			}
		}

	case TabConfig:
		filename = filepath.Join(downloadsDir, fmt.Sprintf("kedastral-config-%s.json", timestamp))
		configData := map[string]interface{}{
			"workload":         m.cfg.Workload,
			"forecaster_url":   m.cfg.ForecasterURL,
			"scaler_url":       m.cfg.ScalerURL,
			"refresh_interval": m.cfg.RefreshInterval.String(),
			"lead_time":        m.cfg.LeadTime.String(),
		}

		if m.snapshot != nil {
			configData["metric"] = m.snapshot.Snapshot.Metric
			configData["step_seconds"] = m.snapshot.Snapshot.StepSeconds
			configData["horizon_seconds"] = m.snapshot.Snapshot.HorizonSeconds
		}

		data, err := json.MarshalIndent(configData, "", "  ")
		if err != nil {
			return fmt.Errorf("failed to marshal config: %w", err)
		}
		exportErr = os.WriteFile(filename, data, 0644)

	case TabLogs:
		filename = filepath.Join(downloadsDir, fmt.Sprintf("kedastral-logs-%s.txt", timestamp))
		logsContent := "Kedastral TUI Logs\n"
		logsContent += "==================\n\n"
		logsContent += "(Log content from bottom panel would go here)\n"
		exportErr = os.WriteFile(filename, []byte(logsContent), 0644)

	default:
		return fmt.Errorf("cannot export from this tab")
	}

	if exportErr != nil {
		return fmt.Errorf("failed to write file: %w", exportErr)
	}

	m.toastManager.Add(fmt.Sprintf("✓ Exported to %s", filepath.Base(filename)), components.ToastSuccess, 3*time.Second)
	return nil
}
