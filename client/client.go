// Package client provides HTTP client for kedastral forecaster and scaler APIs.
package client

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// Client handles HTTP communication with forecaster and scaler services.
type Client struct {
	forecasterURL string
	scalerURL     string
	httpClient    *http.Client
}

// Snapshot represents a forecast snapshot from the forecaster.
type Snapshot struct {
	Workload        string    `json:"workload"`
	Metric          string    `json:"metric"`
	GeneratedAt     time.Time `json:"generatedAt"`
	StepSeconds     int       `json:"stepSeconds"`
	HorizonSeconds  int       `json:"horizonSeconds"`
	Values          []float64 `json:"values"`
	DesiredReplicas []int     `json:"desiredReplicas"`
}

// SnapshotData contains enriched snapshot information for display.
type SnapshotData struct {
	Snapshot      Snapshot
	Stale         bool
	ForecastAge   time.Duration
	LeadTimeIndex int
}

// ScalerMetrics contains parsed metrics from the scaler.
type ScalerMetrics struct {
	Active            bool
	ForecastAgeSeen   float64
	DesiredReplicas   int
	ConnectionHealthy bool
}

// New creates a new Client instance.
func New(forecasterURL, scalerURL string) *Client {
	return &Client{
		forecasterURL: forecasterURL,
		scalerURL:     scalerURL,
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

// GetSnapshot fetches the current forecast snapshot for the given workload.
func (c *Client) GetSnapshot(ctx context.Context, workload string, leadTime time.Duration) (*SnapshotData, error) {
	url := fmt.Sprintf("%s/forecast/current?workload=%s", c.forecasterURL, workload)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch snapshot: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("forecaster returned status %d: %s", resp.StatusCode, string(body))
	}

	var snapshot Snapshot
	if err := json.NewDecoder(resp.Body).Decode(&snapshot); err != nil {
		return nil, fmt.Errorf("failed to decode snapshot: %w", err)
	}

	stale := resp.Header.Get("X-Kedastral-Stale") == "true"
	forecastAge := time.Since(snapshot.GeneratedAt)

	leadTimeIndex := 0
	if snapshot.StepSeconds > 0 {
		stepDuration := time.Duration(snapshot.StepSeconds) * time.Second
		leadSteps := int(leadTime / stepDuration)
		if leadSteps >= len(snapshot.DesiredReplicas) {
			leadSteps = len(snapshot.DesiredReplicas) - 1
		}
		if leadSteps > 0 {
			leadTimeIndex = leadSteps
		}
	}

	return &SnapshotData{
		Snapshot:      snapshot,
		Stale:         stale,
		ForecastAge:   forecastAge,
		LeadTimeIndex: leadTimeIndex,
	}, nil
}

// GetScalerMetrics fetches and parses metrics from the scaler.
func (c *Client) GetScalerMetrics(ctx context.Context) (*ScalerMetrics, error) {
	url := fmt.Sprintf("%s/metrics", c.scalerURL)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch metrics: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("scaler returned status %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read metrics: %w", err)
	}

	metrics := &ScalerMetrics{
		ConnectionHealthy: true,
	}

	metrics.parsePrometheusMetrics(string(body))

	return metrics, nil
}

// GetHealthStatus checks the health of both forecaster and scaler.
func (c *Client) GetHealthStatus(ctx context.Context) (forecasterHealthy, scalerHealthy bool) {
	forecasterHealthy = c.checkHealth(ctx, c.forecasterURL)
	scalerHealthy = c.checkHealth(ctx, c.scalerURL)
	return
}

func (c *Client) checkHealth(ctx context.Context, baseURL string) bool {
	url := fmt.Sprintf("%s/healthz", baseURL)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return false
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return false
	}
	defer resp.Body.Close()

	return resp.StatusCode == http.StatusOK
}

// parsePrometheusMetrics parses Prometheus text format metrics.
func (m *ScalerMetrics) parsePrometheusMetrics(body string) {
	lines := splitLines(body)

	for _, line := range lines {
		if len(line) == 0 || line[0] == '#' {
			continue
		}

		if contains(line, "kedastral_scaler_desired_replicas_returned") {
			if val := extractValue(line); val != 0 {
				m.DesiredReplicas = int(val)
			}
		} else if contains(line, "kedastral_scaler_forecast_age_seen_seconds") {
			m.ForecastAgeSeen = extractValue(line)
		} else if contains(line, "kedastral_scaler_grpc_requests_total") && contains(line, `status="active"`) {
			if extractValue(line) > 0 {
				m.Active = true
			}
		}
	}
}

func splitLines(s string) []string {
	var lines []string
	var current string
	for _, ch := range s {
		if ch == '\n' {
			lines = append(lines, current)
			current = ""
		} else {
			current += string(ch)
		}
	}
	if current != "" {
		lines = append(lines, current)
	}
	return lines
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && findSubstring(s, substr)
}

func findSubstring(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}

func extractValue(line string) float64 {
	for i := len(line) - 1; i >= 0; i-- {
		if line[i] == ' ' || line[i] == '\t' {
			var val float64
			if _, err := fmt.Sscanf(line[i+1:], "%f", &val); err == nil {
				return val
			}
			break
		}
	}
	return 0
}
