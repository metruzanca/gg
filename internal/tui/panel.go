package tui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
	ggit "github.com/metru/gg/internal/git"
)

type Panel struct {
	label   string
	tree    []*ggit.Entry
	visible []*ggit.Entry
}

func NewPanel(label string) *Panel {
	return &Panel{label: label}
}

func (p *Panel) SetTree(tree []*ggit.Entry) {
	p.tree = tree
}

func (p *Panel) Rebuild(expandedDirs map[string]bool, allExpanded bool) {
	p.visible = ggit.FlattenVisibleForPanel(p.tree, expandedDirs, allExpanded)
}

func (p *Panel) Len() int {
	return len(p.visible)
}

func (p *Panel) EntryAt(idx int) (*ggit.Entry, bool) {
	if idx < 0 || idx >= len(p.visible) {
		return nil, false
	}
	return p.visible[idx], true
}

func (p *Panel) View(localCursor int, expandedDirs map[string]bool, allExpanded bool, width int) string {
	var b strings.Builder

	headerW := width - 2
	pad := strings.Repeat("─", max(0, headerW-len(p.label)-6))
	b.WriteString(SectionHeaderStyle.Render(fmt.Sprintf("── %s %s", p.label, pad)) + "\n")

	if len(p.visible) == 0 {
		b.WriteString(DimStyle.Render(fmt.Sprintf("  (no %s files)", strings.ToLower(p.label))) + "\n")
		return b.String()
	}

	for i, e := range p.visible {
		b.WriteString(p.renderEntry(e, i == localCursor, expandedDirs, allExpanded) + "\n")
	}

	return b.String()
}

func (p *Panel) renderEntry(e *ggit.Entry, selected bool, expandedDirs map[string]bool, allExpanded bool) string {
	var line string
	indent := strings.Repeat("  ", e.Depth)

	if e.IsDir {
		expanded := allExpanded || expandedDirs[e.Path]
		icon := "▸"
		if expanded {
			icon = "▾"
		}
		line = fmt.Sprintf("%s%s %s/", indent, icon, e.Name)
	} else {
		var statusStyle lipgloss.Style
		if e.Section == ggit.SectionStaged {
			statusStyle = StagedStatusStyle(e.Status)
		} else {
			statusStyle = UnstagedStatusStyle(e.Status)
		}
		label := ggit.StatusLabel(e.Status)
		line = fmt.Sprintf("%s  %s  %s", indent, statusStyle.Render(label), DimStyle.Render(e.Name))
	}

	if selected {
		return SelectedStyle.Render(line)
	}
	return line
}
