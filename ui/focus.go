package ui

import tea "github.com/charmbracelet/bubbletea"

func (m Model) handleFocusSwitch(msg tea.KeyMsg) (Model, tea.Cmd) {
	switch msg.String() {
	case "tab":
		// Move focus to next panel (right)
		m.focusedPanel = (m.focusedPanel + 1) % 3
		return m, nil

	case "shift+tab":
		// Move focus to previous panel (left)
		m.focusedPanel = (m.focusedPanel - 1 + 3) % 3
		return m, nil

	case "w":
		// Jump to sidebar
		m.focusedPanel = PanelSidebar
		return m, nil

	case "m":
		// Jump to main panel
		m.focusedPanel = PanelMain
		return m, nil

	case "[":
		// Toggle sidebar collapse
		m.layoutMgr.ToggleSidebar()
		return m, nil

	case "]":
		// Toggle bottom panel collapse
		m.layoutMgr.ToggleBottom()
		return m, nil
	}

	return m, nil
}
