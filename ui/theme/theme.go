package theme

import "github.com/charmbracelet/lipgloss"

type Theme struct {
	Name string

	Primary       lipgloss.Color
	Secondary     lipgloss.Color
	Success       lipgloss.Color
	Warning       lipgloss.Color
	Error         lipgloss.Color
	Info          lipgloss.Color
	Muted         lipgloss.Color
	Background    lipgloss.Color
	Foreground    lipgloss.Color
	Border        lipgloss.Color
	BorderFocused lipgloss.Color
}

var Dark = &Theme{
	Name:          "dark",
	Primary:       lipgloss.Color("39"),
	Secondary:     lipgloss.Color("205"),
	Success:       lipgloss.Color("42"),
	Warning:       lipgloss.Color("220"),
	Error:         lipgloss.Color("196"),
	Info:          lipgloss.Color("39"),
	Muted:         lipgloss.Color("241"),
	Background:    lipgloss.Color("235"),
	Foreground:    lipgloss.Color("252"),
	Border:        lipgloss.Color("241"),
	BorderFocused: lipgloss.Color("39"),
}

var Light = &Theme{
	Name:          "light",
	Primary:       lipgloss.Color("27"),
	Secondary:     lipgloss.Color("162"),
	Success:       lipgloss.Color("34"),
	Warning:       lipgloss.Color("172"),
	Error:         lipgloss.Color("160"),
	Info:          lipgloss.Color("27"),
	Muted:         lipgloss.Color("240"),
	Background:    lipgloss.Color("255"),
	Foreground:    lipgloss.Color("235"),
	Border:        lipgloss.Color("240"),
	BorderFocused: lipgloss.Color("27"),
}

var Available = []*Theme{Dark, Light}

func Get(name string) *Theme {
	for _, t := range Available {
		if t.Name == name {
			return t
		}
	}
	return Dark
}

func (t *Theme) Title() lipgloss.Style {
	return lipgloss.NewStyle().Bold(true).Foreground(t.Primary)
}

func (t *Theme) Mode() lipgloss.Style {
	return lipgloss.NewStyle().Bold(true).Foreground(t.Warning)
}

func (t *Theme) SuccessText() lipgloss.Style {
	return lipgloss.NewStyle().Foreground(t.Success)
}

func (t *Theme) ErrorText() lipgloss.Style {
	return lipgloss.NewStyle().Foreground(t.Error)
}

func (t *Theme) MutedText() lipgloss.Style {
	return lipgloss.NewStyle().Foreground(t.Muted)
}

func (t *Theme) GetBorderColor(focused bool) lipgloss.Color {
	if focused {
		return t.BorderFocused
	}
	return t.Border
}

func (t *Theme) ActiveTab() lipgloss.Style {
	return lipgloss.NewStyle().
		Bold(true).
		Foreground(t.Success).
		Background(t.Background).
		Padding(0, 2)
}

func (t *Theme) InactiveTab() lipgloss.Style {
	return lipgloss.NewStyle().
		Foreground(t.Muted).
		Padding(0, 2)
}
