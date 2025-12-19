// Package config provides configuration parsing and management for the TUI.
package config

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// Config holds all TUI configuration.
type Config struct {
	ForecasterURL   string        `json:"forecaster_url,omitempty"`
	ScalerURL       string        `json:"scaler_url,omitempty"`
	Workload        string        `json:"workload,omitempty"`
	RefreshInterval time.Duration `json:"refresh_interval,omitempty"`
	LeadTime        time.Duration `json:"lead_time,omitempty"`
	LogLevel        string        `json:"log_level,omitempty"`
}

// ParseFlags parses configuration from file, environment variables, and command-line flags.
func ParseFlags() (*Config, bool) {
	fileConfig := loadConfigFile()

	cfg := &Config{}

	forecasterDefault := fileConfig.ForecasterURL
	if forecasterDefault == "" {
		forecasterDefault = getEnv("FORECASTER_URL", "")
	}

	scalerDefault := fileConfig.ScalerURL
	if scalerDefault == "" {
		scalerDefault = getEnv("SCALER_URL", "http://localhost:8082")
	}

	workloadDefault := fileConfig.Workload
	if workloadDefault == "" {
		workloadDefault = getEnv("WORKLOAD", "")
	}

	refreshDefault := fileConfig.RefreshInterval
	if refreshDefault == 0 {
		refreshDefault = getEnvDuration("REFRESH_INTERVAL", 5*time.Second)
	}

	leadTimeDefault := fileConfig.LeadTime
	if leadTimeDefault == 0 {
		leadTimeDefault = getEnvDuration("LEAD_TIME", 5*time.Minute)
	}

	logLevelDefault := fileConfig.LogLevel
	if logLevelDefault == "" {
		logLevelDefault = getEnv("LOG_LEVEL", "error")
	}

	flag.StringVar(&cfg.ForecasterURL, "forecaster-url", forecasterDefault, "Forecaster HTTP URL (required)")
	flag.StringVar(&cfg.ScalerURL, "scaler-url", scalerDefault, "Scaler HTTP URL")
	flag.StringVar(&cfg.Workload, "workload", workloadDefault, "Workload name to monitor (required)")
	flag.DurationVar(&cfg.RefreshInterval, "refresh-interval", refreshDefault, "Refresh interval in live mode")
	flag.DurationVar(&cfg.LeadTime, "lead-time", leadTimeDefault, "Lead time for replica selection highlighting")
	flag.StringVar(&cfg.LogLevel, "log-level", logLevelDefault, "Log level: debug, info, warn, error")

	flag.Parse()

	needsSetup := cfg.ForecasterURL == "" || cfg.Workload == ""

	if cfg.RefreshInterval > 0 && cfg.RefreshInterval < 1*time.Second {
		fmt.Fprintln(os.Stderr, "Error: --refresh-interval must be at least 1 second")
		os.Exit(1)
	}

	return cfg, needsSetup
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvDuration(key string, defaultValue time.Duration) time.Duration {
	if value := os.Getenv(key); value != "" {
		if d, err := time.ParseDuration(value); err == nil {
			return d
		}
	}
	return defaultValue
}

// getConfigPath returns the path to the configuration file.
func getConfigPath() string {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return ""
	}
	return filepath.Join(homeDir, ".config", "kedastral-tui", "config.json")
}

// LoadConfigFile loads configuration from the config file if it exists.
func LoadConfigFile() (*Config, error) {
	cfg := &Config{}

	configPath := getConfigPath()
	if configPath == "" {
		return nil, fmt.Errorf("unable to determine config path")
	}

	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	if err := json.Unmarshal(data, cfg); err != nil {
		return nil, fmt.Errorf("failed to parse config file: %w", err)
	}

	return cfg, nil
}

func loadConfigFile() *Config {
	cfg, err := LoadConfigFile()
	if err != nil {
		return &Config{}
	}
	return cfg
}
