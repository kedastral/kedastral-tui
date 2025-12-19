# Kedastral TUI

A terminal user interface (TUI) for monitoring Kedastral forecaster and scaler services.

## Features

- 🚀 **Interactive Setup**: First-run wizard to configure connection settings
- 📊 **Live Monitoring**: Real-time updates of forecast and scaler status
- ⏸️ **Pause/Resume**: Toggle between live and paused modes
- 🔄 **Manual Refresh**: Force refresh on demand
- 📁 **Config File**: Persistent configuration in `~/.config/kedastral-tui/config.json`

## Installation

### Pre-built Binaries (Recommended)

Download the latest release for your platform from the [releases page](https://github.com/HatiCode/kedastral-tui/releases).

**Linux (amd64)**
```bash
wget https://github.com/HatiCode/kedastral-tui/releases/latest/download/kedastral-tui-<version>-linux-amd64.tar.gz
tar -xzf kedastral-tui-<version>-linux-amd64.tar.gz
chmod +x kedastral-tui
sudo mv kedastral-tui /usr/local/bin/
```

**macOS (Apple Silicon)**
```bash
wget https://github.com/HatiCode/kedastral-tui/releases/latest/download/kedastral-tui-<version>-darwin-arm64.tar.gz
tar -xzf kedastral-tui-<version>-darwin-arm64.tar.gz
chmod +x kedastral-tui
sudo mv kedastral-tui /usr/local/bin/
```

### From Source

```bash
git clone https://github.com/HatiCode/kedastral-tui.git
cd kedastral-tui
make build
```

The binary will be created at `bin/kedastral-tui`.

## Usage

### First Run

On first launch (or if no config exists), you'll be guided through an interactive setup:

```bash
./bin/kedastral-tui
```

The wizard will ask for:
1. **Forecaster URL** (e.g., `http://localhost:8081`)
2. **Scaler URL** (default: `http://localhost:8082`)
3. **Workload name** (e.g., `test-app`)

Configuration is saved to `~/.config/kedastral-tui/config.json`.

### Subsequent Runs

After setup, simply run:

```bash
./bin/kedastral-tui
```

Or with custom settings:

```bash
./bin/kedastral-tui --forecaster-url=http://kedastral-forecaster:8081 --workload=my-app
```

### Keyboard Controls

- **SPACE**: Toggle between live and paused modes
- **R**: Manual refresh (fetch latest data)
- **H**: Toggle help screen
- **Q** or **Ctrl+C**: Quit

## Configuration

Configuration priority (highest to lowest):
1. Command-line flags
2. Environment variables
3. Config file (`~/.config/kedastral-tui/config.json`)
4. Defaults

### Command-line Flags

```bash
--forecaster-url    Forecaster HTTP URL (required)
--scaler-url        Scaler HTTP URL (default: http://localhost:8082)
--workload          Workload name to monitor (required)
--refresh-interval  Refresh interval in live mode (default: 5s)
--lead-time         Lead time for replica selection (default: 5m)
--log-level         Log level: debug, info, warn, error (default: error)
--version           Print version and exit
```

### Environment Variables

```bash
export FORECASTER_URL=http://localhost:8081
export SCALER_URL=http://localhost:8082
export WORKLOAD=test-app
export REFRESH_INTERVAL=5s
export LEAD_TIME=5m
```

### Config File

Location: `~/.config/kedastral-tui/config.json`

Example:
```json
{
  "forecaster_url": "http://localhost:8081",
  "scaler_url": "http://localhost:8082",
  "workload": "test-app"
}
```

## Development

### Build

```bash
make build
```

### Run

```bash
make run
```

### Clean

```bash
make clean
```

### Format

```bash
make fmt
```

### Test

```bash
make test
```

## Features in Detail

### 📊 **Visual Components**
- **Status Bar**: Real-time connection health, forecast age, mode indicator
- **Forecast Chart**: ASCII line chart showing predicted metric values over the horizon
- **Replica Table**: Detailed scaling decisions with lead time highlighting
- **Scaler Status**: Active/inactive state and current replica count
- **Help Screen**: Press `H` for keyboard shortcuts and panel descriptions

### ⚙️ **Functionality**
- **Interactive Setup**: First-run wizard saves configuration automatically
- **Live Monitoring**: Auto-refresh every 5s (configurable)
- **Pause Mode**: Freeze updates to inspect current state
- **Manual Refresh**: Force data fetch with `R` key
- **Multi-source Config**: File → Env vars → Flags precedence

## Screenshots

```
┌──────────────────────────────────────────────────────────────┐
│ Kedastral Monitor - workload: test-app  [LIVE]  Last: 2s ago│
│ Status: Forecaster ✓  Scaler ✓  Forecast age: 3s            │
├──────────────────────────────────────────────────────────────┤
│ FORECAST TIMELINE (next 30min)                              │
│                                                              │
│ RPS                                                          │
│ 220 ┤     ●●●                                                │
│ 180 ┤   ●●   ●●      ●●●                                     │
│ 140 ┤ ●●       ●●●●●●   ●●                                   │
│ 100 ┤●                    ●●                                 │
│     └──────────────────────────────────────────────────────  │
│      Now                                            +30m     │
├──────────────────────────────────────────────────────────────┤
│ REPLICA SCALING DECISIONS                                    │
│                                                              │
│ Time        Forecast      Desired   Selected                │
│ ───────────────────────────────────────────────────────────  │
│ Now         36.2 rps      2                                  │
│ +5min       180.5 rps     3         ← SELECTED (lead=5m)    │
│ +10min      220.3 rps     3                                  │
│ ...                                                          │
├──────────────────────────────────────────────────────────────┤
│ SCALER STATUS                                                │
│ Active: ✓  Desired replicas: 3  Forecast age seen: 2.1s     │
│                                                              │
│ [SPACE] pause  [R] refresh  [H] help  [Q] quit              │
└──────────────────────────────────────────────────────────────┘
```

## Architecture

```
kedastral-tui/
├── main.go              # Entry point with version support
├── config/
│   └── config.go        # Multi-source configuration management
├── client/
│   └── client.go        # HTTP client for forecaster/scaler APIs
├── ui/
│   ├── model.go         # Bubble Tea model (state management)
│   ├── update.go        # Event handling and async updates
│   ├── view.go          # Main view rendering
│   ├── messages.go      # Custom Bubble Tea messages
│   └── setup.go         # Interactive setup wizard
└── components/          # Reusable UI components
    ├── status_bar.go    # Status bar with health indicators
    ├── forecast_chart.go # ASCII line chart renderer
    ├── replica_table.go  # Replica scaling decisions table
    └── help.go           # Help screen component
```

## CI/CD

This project uses GitHub Actions for:
- **Continuous Integration**: Automated testing and linting on every push
- **Releases**: Automatic binary builds for Linux, macOS, and Windows on tagged releases
- **Security**: Weekly security scans with Gosec

To create a release:
```bash
git tag v0.1.0
git push origin v0.1.0
```

## License

See LICENSE file.
