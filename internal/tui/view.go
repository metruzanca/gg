package tui

import (
	"strings"
)

func (m *Model) View() string {
	if m.quitting {
		return ""
	}

	switch m.modal {
	case ModalConfirm:
		return m.confirmView()
	case ModalCommit:
		return m.commitView()
	case ModalHelp:
		return m.helpView()
	}

	return m.mainView()
}

func (m *Model) mainView() string {
	var b strings.Builder
	width := m.width
	if width == 0 {
		width = 80
	}

	if m.err != nil {
		b.WriteString(ErrorStyle.Render(" error: "+m.err.Error()) + "\n")
		m.err = nil
	}

	b.WriteString(m.staged.View(m.cursor, m.expandedDirs, m.allExpanded, width))
	b.WriteString("\n")
	b.WriteString(m.unstaged.View(m.cursor-m.staged.Len(), m.expandedDirs, m.allExpanded, width))
	b.WriteString("\n")
	b.WriteString(m.footerView(width))

	return b.String()
}

func (m *Model) footerView(width int) string {
	keys := []string{
		"jk/↑↓:move",
		"hl/←→:expand",
		"s:stage",
		"a:stage-all",
		"e:expand-all",
		"1:staged",
		"2:worktree",
		"enter:commit",
		"?:help",
		"q:quit",
	}
	return FooterStyle.Render(strings.Repeat("─", width) + "\n " + strings.Join(keys, "  "))
}
