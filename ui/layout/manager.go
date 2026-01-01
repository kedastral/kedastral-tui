// Package layout provides layout management for the TUI panels.
package layout

// Rect represents a rectangular area on the screen.
type Rect struct {
	X, Y int // Position
	W, H int // Width and Height
}

// LayoutManager computes panel dimensions based on terminal size.
type LayoutManager struct {
	termWidth  int
	termHeight int

	sidebarWidth int // Computed: 25% of width, min 20, max 40
	bottomHeight int // Computed: 20% of height, min 5, max 15

	sidebarCollapsed bool
	bottomCollapsed  bool
}

// Layout holds the computed dimensions for all panels.
type Layout struct {
	Sidebar Rect
	Main    Rect
	Bottom  Rect
}

// NewLayoutManager creates a new layout manager.
func NewLayoutManager() *LayoutManager {
	return &LayoutManager{
		sidebarCollapsed: false,
		bottomCollapsed:  false,
	}
}

// SetTerminalSize updates the terminal dimensions.
func (l *LayoutManager) SetTerminalSize(width, height int) {
	l.termWidth = width
	l.termHeight = height
}

// ToggleSidebar toggles the sidebar collapsed state.
func (l *LayoutManager) ToggleSidebar() {
	l.sidebarCollapsed = !l.sidebarCollapsed
}

// ToggleBottom toggles the bottom panel collapsed state.
func (l *LayoutManager) ToggleBottom() {
	l.bottomCollapsed = !l.bottomCollapsed
}

// Compute calculates panel dimensions based on current settings.
func (l *LayoutManager) Compute() Layout {
	layout := Layout{}

	// Sidebar: 25% of width, min 20 cols, max 40 cols
	sidebarWidth := clamp(l.termWidth/4, 20, 40)
	if l.sidebarCollapsed {
		sidebarWidth = 0
	}

	// Bottom: 20% of height, min 5 rows, max 15 rows
	bottomHeight := clamp(l.termHeight/5, 5, 15)
	if l.bottomCollapsed {
		bottomHeight = 0
	}

	// Main: remaining space
	// Reserve 2 columns for borders between panels (1 on each side)
	borderWidth := 0
	if sidebarWidth > 0 {
		borderWidth = 2
	}

	mainWidth := l.termWidth - sidebarWidth - borderWidth
	if mainWidth < 0 {
		mainWidth = 0
	}

	// Reserve 2 rows for borders between main and bottom
	borderHeight := 0
	if bottomHeight > 0 {
		borderHeight = 2
	}

	mainHeight := l.termHeight - bottomHeight - borderHeight
	if mainHeight < 0 {
		mainHeight = 0
	}

	// Sidebar: left side, full height
	layout.Sidebar = Rect{
		X: 0,
		Y: 0,
		W: sidebarWidth,
		H: l.termHeight,
	}

	// Main: right of sidebar, above bottom
	layout.Main = Rect{
		X: sidebarWidth + borderWidth,
		Y: 0,
		W: mainWidth,
		H: mainHeight,
	}

	// Bottom: right of sidebar, below main
	layout.Bottom = Rect{
		X: sidebarWidth + borderWidth,
		Y: mainHeight + borderHeight,
		W: mainWidth,
		H: bottomHeight,
	}

	return layout
}

// clamp returns val clamped to the range [min, max].
func clamp(val, min, max int) int {
	if val < min {
		return min
	}
	if val > max {
		return max
	}
	return val
}
