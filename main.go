package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/HatiCode/kedastral-tui/client"
	"github.com/HatiCode/kedastral-tui/config"
	"github.com/HatiCode/kedastral-tui/ui"
)

var version = "dev"

func main() {
	for _, arg := range os.Args[1:] {
		if arg == "--version" || arg == "-version" {
			fmt.Printf("kedastral-tui %s\n", version)
			os.Exit(0)
		}
	}

	cfg, needsSetup := config.ParseFlags()

	if needsSetup {
		setupModel := ui.NewSetupModel()
		p := tea.NewProgram(setupModel)

		if _, err := p.Run(); err != nil {
			fmt.Fprintf(os.Stderr, "Setup error: %v\n", err)
			os.Exit(1)
		}

		newCfg, err := config.LoadConfigFile()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to load config after setup: %v\n", err)
			os.Exit(1)
		}

		if newCfg.ForecasterURL == "" || newCfg.Workload == "" {
			fmt.Fprintln(os.Stderr, "Setup incomplete. Please try again.")
			os.Exit(1)
		}

		if newCfg.RefreshInterval == 0 {
			newCfg.RefreshInterval = cfg.RefreshInterval
		}
		if newCfg.LeadTime == 0 {
			newCfg.LeadTime = cfg.LeadTime
		}
		if newCfg.LogLevel == "" {
			newCfg.LogLevel = cfg.LogLevel
		}

		cfg = newCfg
	}

	c := client.New(cfg.ForecasterURL, cfg.ScalerURL)

	model := ui.NewModel(cfg, c)

	p := tea.NewProgram(model, tea.WithAltScreen())

	if _, err := p.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
