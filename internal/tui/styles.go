package tui

import "github.com/charmbracelet/lipgloss"

var (
	TitleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("6"))

	NormalStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("7"))

	DimStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("8"))

	SelectedStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("0")).
			Background(lipgloss.Color("4"))

	StagedColor   = lipgloss.Color("2")
	ModifiedColor = lipgloss.Color("3")
	AddedColor    = lipgloss.Color("2")
	DeletedColor  = lipgloss.Color("1")
	UntrackedColor = lipgloss.Color("8")

	ModalBorderStyle = lipgloss.NewStyle().
				Border(lipgloss.RoundedBorder()).
				BorderForeground(lipgloss.Color("4"))

	FooterStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("8"))

	SectionHeaderStyle = lipgloss.NewStyle().
				Bold(true).
				Foreground(lipgloss.Color("4")).
				Background(lipgloss.Color("0"))

	HelpTitleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("6"))

	HelpKeyStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("2"))

	InputStyle = lipgloss.NewStyle().
			Border(lipgloss.NormalBorder()).
			BorderForeground(lipgloss.Color("8"))

	FocusedInputStyle = lipgloss.NewStyle().
				Border(lipgloss.NormalBorder()).
				BorderForeground(lipgloss.Color("4"))

	SuggestionStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("8"))

	SelectedSuggestionStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("0")).
				Background(lipgloss.Color("4"))

	ErrorStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("1")).
			Bold(true)
)

func StagedStatusStyle(b byte) lipgloss.Style {
	switch b {
	case 'M':
		return lipgloss.NewStyle().Foreground(ModifiedColor)
	case 'A':
		return lipgloss.NewStyle().Foreground(AddedColor)
	case 'D':
		return lipgloss.NewStyle().Foreground(DeletedColor)
	default:
		return NormalStyle
	}
}

func UnstagedStatusStyle(b byte) lipgloss.Style {
	switch b {
	case 'M':
		return lipgloss.NewStyle().Foreground(ModifiedColor)
	case 'D':
		return lipgloss.NewStyle().Foreground(DeletedColor)
	case '?':
		return lipgloss.NewStyle().Foreground(UntrackedColor)
	default:
		return NormalStyle
	}
}
