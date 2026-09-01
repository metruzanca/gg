package tui

import (
	"fmt"
	"strings"

	ggit "github.com/metruzanca/gg/internal/git"
)

func (m *Model) confirmView() string {
	var b strings.Builder
	b.WriteString("\n")
	msg := m.confirmMsg
	b.WriteString(msg + "\n\n")

	yesLabel := "[ Yes ]"
	noLabel := "[ No ]"

	if m.confirmChoice == ConfirmYes {
		yesLabel = SelectedStyle.Render("[ Yes ]")
	} else {
		yesLabel = NormalStyle.Render("[ Yes ]")
	}
	if m.confirmChoice == ConfirmNo {
		noLabel = SelectedStyle.Render("[ No ]")
	} else {
		noLabel = NormalStyle.Render("[ No ]")
	}

	b.WriteString("  " + yesLabel + "  " + noLabel + "\n\n")
	b.WriteString(DimStyle.Render(" h/l ←/→: choose   enter: confirm   esc: cancel"))
	return ModalBorderStyle.Render(b.String())
}

func (m *Model) commitView() string {
	var b strings.Builder

	files, _ := ggit.StagedFiles()
	if len(files) > 0 {
		b.WriteString(DimStyle.Render("Files to commit:") + "\n")
		for _, f := range files {
			b.WriteString(DimStyle.Render("  "+f) + "\n")
		}
		b.WriteString("\n")
	} else {
		b.WriteString(DimStyle.Render("(no staged files)") + "\n\n")
	}

	typeLabel := "Type:"
	if m.commitField == fieldType {
		typeLabel = SelectedStyle.Render("Type:")
	}
	b.WriteString(typeLabel + " " + m.typeInput.View() + "\n")

	if m.commitField == fieldType && len(m.filteredTypes) > 0 {
		b.WriteString("  ")
		for i, t := range m.filteredTypes {
			s := t.Name
			if i == m.typeSuggestIdx {
				s = SelectedSuggestionStyle.Render(t.Name)
			} else {
				s = SuggestionStyle.Render(t.Name)
			}
			b.WriteString(s + "  ")
		}
		b.WriteString("\n")
	}

	titleLabel := "Title:"
	if m.commitField == fieldTitle {
		titleLabel = SelectedStyle.Render("Title:")
	}
	b.WriteString("\n" + titleLabel + "\n" + m.titleInput.View() + "\n")

	descLabel := "Description:"
	if m.commitField == fieldDescription {
		descLabel = SelectedStyle.Render("Description:")
	}
	b.WriteString("\n" + descLabel + "\n" + m.descInput.View())

	if m.commitErr != "" {
		b.WriteString("\n" + ErrorStyle.Render(" "+m.commitErr))
	}

	footer := "\n\n " + DimStyle.Render("tab:next   enter:commit   esc:cancel   ↑↓:type suggestion")
	b.WriteString(footer)

	return ModalBorderStyle.Width(min(70, m.width-4)).Render(b.String())
}

func (m *Model) helpView() string {
	sections := []struct {
		title string
		items [][2]string
	}{
		{
			title: "Navigation",
			items: [][2]string{
				{"j / ↓", "move cursor down"},
				{"k / ↑", "move cursor up"},
				{"h / ←", "collapse directory"},
				{"l / →", "expand directory"},
				{"1", "jump to top of staged files"},
				{"2", "jump to top of worktree"},
			},
		},
		{
			title: "Actions",
			items: [][2]string{
				{"s", "stage / unstage file"},
				{"a", "stage all unstaged files"},
				{"e", "expand / collapse all directories"},
				{"enter", "open commit modal"},
				{"q", "quit from anywhere"},
				{"?", "toggle this help"},
			},
		},
		{
			title: "Modals",
			items: [][2]string{
				{"esc", "close any modal"},
				{"enter", "confirm / commit (when form is valid)"},
				{"tab", "next field (commit modal)"},
				{"←/h  →/l", "choose option (confirm modal)"},
			},
		},
		{
			title: "Commit Modal",
			items: [][2]string{
				{"type", "type to filter conventional commit types, ↓↑ to select"},
				{"title", "required commit title"},
				{"description", "optional body text"},
			},
		},
	}

	var b strings.Builder
	b.WriteString(HelpTitleStyle.Render("Help") + "\n\n")

	for _, sec := range sections {
		b.WriteString(SectionHeaderStyle.Render(fmt.Sprintf(" %s ", sec.title)) + "\n")
		for _, item := range sec.items {
			b.WriteString(fmt.Sprintf("  %s  %s\n",
				HelpKeyStyle.Render(item[0]),
				DimStyle.Render(item[1]),
			))
		}
		b.WriteString("\n")
	}

	b.WriteString(DimStyle.Render(" press any key to close"))

	return ModalBorderStyle.Render(b.String())
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
