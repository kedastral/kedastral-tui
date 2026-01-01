package components

import (
	"time"

	"github.com/charmbracelet/lipgloss"
)

type ToastType int

const (
	ToastInfo ToastType = iota
	ToastSuccess
	ToastWarning
	ToastError
)

type Toast struct {
	Message   string
	Type      ToastType
	ExpiresAt time.Time
}

type ToastManager struct {
	toasts []Toast
	width  int
}

func NewToastManager(width int) *ToastManager {
	return &ToastManager{
		toasts: make([]Toast, 0),
		width:  width,
	}
}

func (t *ToastManager) Add(message string, toastType ToastType, duration time.Duration) {
	toast := Toast{
		Message:   message,
		Type:      toastType,
		ExpiresAt: time.Now().Add(duration),
	}
	t.toasts = append(t.toasts, toast)
}

func (t *ToastManager) Update() {
	now := time.Now()
	var active []Toast
	for _, toast := range t.toasts {
		if toast.ExpiresAt.After(now) {
			active = append(active, toast)
		}
	}
	t.toasts = active
}

func (t *ToastManager) Render() string {
	t.Update()

	if len(t.toasts) == 0 {
		return ""
	}

	var result string
	for _, toast := range t.toasts {
		style := t.getStyle(toast.Type)
		result += style.Render(toast.Message) + "\n"
	}

	return result
}

func (t *ToastManager) getStyle(toastType ToastType) lipgloss.Style {
	base := lipgloss.NewStyle().
		Padding(0, 1).
		MarginBottom(1)

	switch toastType {
	case ToastSuccess:
		return base.
			Foreground(lipgloss.Color("42")).
			Background(lipgloss.Color("235"))
	case ToastWarning:
		return base.
			Foreground(lipgloss.Color("220")).
			Background(lipgloss.Color("235"))
	case ToastError:
		return base.
			Foreground(lipgloss.Color("196")).
			Background(lipgloss.Color("235"))
	default:
		return base.
			Foreground(lipgloss.Color("39")).
			Background(lipgloss.Color("235"))
	}
}

func (t *ToastManager) HasToasts() bool {
	t.Update()
	return len(t.toasts) > 0
}
