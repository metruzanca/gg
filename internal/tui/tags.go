package tui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	ggit "github.com/metruzanca/gg/internal/git"
)

type TagMode int

const (
	TagList TagMode = iota
	TagEdit
	TagConfirmPush
)

type TagModel struct {
	tags   []string
	cursor int
	width  int
	err    error
	notice string

	mode    TagMode
	input   textinput.Model
	latest  string
	tagErr  string
	pending string
}

func NewTagModel() *TagModel {
	ti := textinput.New()
	ti.Placeholder = "tag name (e.g. v0.1.0)"
	ti.Width = 40

	m := &TagModel{
		input: ti,
		mode:  TagList,
	}
	m.refresh()
	return m
}

func (m *TagModel) Init() tea.Cmd {
	return nil
}

func (m *TagModel) refresh() {
	tags, err := ggit.Tags()
	if err != nil {
		m.err = err
		m.tags = nil
		return
	}
	m.err = nil
	m.tags = tags
	if m.cursor >= len(m.tags) {
		m.cursor = max(0, len(m.tags)-1)
	}
}

func (m *TagModel) openEdit() {
	latest := ggit.LatestVersion(m.tags)
	m.latest = latest
	m.input = textinput.New()
	m.input.Placeholder = "tag name (e.g. v0.1.0)"
	m.input.Width = 40
	if latest == "" {
		latest = "v0.1.0"
	}
	m.input.SetValue(latest)
	m.input.SetCursor(len(latest))
	m.tagErr = ""
	m.mode = TagEdit
}

func (m *TagModel) createTag() {
	name := m.input.Value()
	if strings.TrimSpace(name) == "" {
		m.tagErr = "tag name is required"
		return
	}
	if err := ggit.CreateTag(name); err != nil {
		m.tagErr = fmt.Sprintf("tag failed: %v", err)
		return
	}
	m.pending = name
	m.mode = TagConfirmPush
	m.refresh()
}

func (m *TagModel) confirmPush(push bool) {
	name := m.pending
	m.pending = ""
	if push {
		if err := ggit.PushTag(name); err != nil {
			m.tagErr = fmt.Sprintf("push failed: %v", err)
			m.notice = fmt.Sprintf("created tag %q", name)
		} else {
			m.notice = fmt.Sprintf("pushed tag %q", name)
		}
	} else {
		m.notice = fmt.Sprintf("created tag %q (not pushed)", name)
	}
	m.mode = TagList
}

func (m *TagModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		return m, nil

	case tea.KeyMsg:
		return m.handleTagKey(msg)
	}

	return m, nil
}

func (m *TagModel) handleTagKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	k := msg.String()

	if k == "ctrl+c" {
		if m.mode == TagEdit {
			m.mode = TagList
			return m, nil
		}
		return m, tea.Quit
	}

	if m.mode == TagEdit {
		return m.handleEditKey(msg)
	}
	if m.mode == TagConfirmPush {
		return m.handleConfirmPushKey(msg)
	}
	return m.handleListKey(msg)
}

func (m *TagModel) handleConfirmPushKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "y":
		m.confirmPush(true)
		return m, nil
	case "n", "esc", "q":
		m.confirmPush(false)
		return m, nil
	}
	return m, nil
}

func (m *TagModel) handleListKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "up", "k":
		if m.cursor > 0 {
			m.cursor--
		}
	case "down", "j":
		if m.cursor < len(m.tags)-1 {
			m.cursor++
		}
	case "n":
		m.openEdit()
		return m, m.input.Focus()
	case "q", "esc":
		return m, tea.Quit
	}
	return m, nil
}

func (m *TagModel) handleEditKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "esc":
		m.mode = TagList
		return m, nil
	case "q":
		if m.input.Value() == "" {
			m.mode = TagList
			return m, nil
		}
	case "enter":
		m.createTag()
		return m, nil
	case "up":
		v, p, _ := ggit.IncrementVersionAt(m.input.Value(), m.input.Position(), 1)
		m.input.SetValue(v)
		m.input.SetCursor(p)
		return m, nil
	case "down":
		v, p, _ := ggit.IncrementVersionAt(m.input.Value(), m.input.Position(), -1)
		m.input.SetValue(v)
		m.input.SetCursor(p)
		return m, nil
	}

	var cmd tea.Cmd
	m.input, cmd = m.input.Update(msg)
	return m, cmd
}

func (m *TagModel) View() string {
	if m.mode == TagEdit {
		return m.editView()
	}
	if m.mode == TagConfirmPush {
		return m.confirmPushView()
	}
	return m.listView()
}

func (m *TagModel) listView() string {
	var b strings.Builder
	width := m.width
	if width == 0 {
		width = 80
	}

	b.WriteString(TitleStyle.Render("Tags") + "\n\n")

	if m.err != nil {
		b.WriteString(ErrorStyle.Render(" error: "+m.err.Error()) + "\n")
	} else if len(m.tags) == 0 {
		b.WriteString(DimStyle.Render("  (no tags)") + "\n")
	} else {
		for i, t := range m.tags {
			if i == m.cursor {
				b.WriteString(SelectedStyle.Render("  "+t) + "\n")
			} else {
				b.WriteString(NormalStyle.Render("  "+t) + "\n")
			}
		}
	}

	if m.notice != "" {
		b.WriteString("\n" + NormalStyle.Render(" "+m.notice) + "\n")
	}

	b.WriteString("\n" + FooterStyle.Render(strings.Repeat("─", width)+"\n "+
		"jk/↑↓:move   n:new tag   q:quit"))

	return b.String()
}

func (m *TagModel) editView() string {
	var b strings.Builder
	b.WriteString(TitleStyle.Render("New tag") + "\n\n")

	latest := m.latest
	if latest == "" {
		latest = "(none)"
	}
	b.WriteString(DimStyle.Render("Current latest: "+latest) + "\n\n")

	b.WriteString("Name:\n" + m.input.View() + "\n")

	if m.tagErr != "" {
		b.WriteString(ErrorStyle.Render(" "+m.tagErr) + "\n")
	}

	b.WriteString("\n" + DimStyle.Render(
		" ←/→:move   ↑/↓:bump number   type:edit   enter:create   esc:cancel"))

	return ModalBorderStyle.Width(max(40, min(60, m.width-4))).Render(b.String())
}

func (m *TagModel) confirmPushView() string {
	var b strings.Builder
	b.WriteString(TitleStyle.Render("Tag created") + "\n\n")
	b.WriteString("Tag " + SelectedStyle.Render(m.pending) + " created.\n")
	b.WriteString("Push it to the remote?\n\n")
	b.WriteString("y:push   n:skip\n")
	return ModalBorderStyle.Width(max(40, min(60, m.width-4))).Render(b.String())
}

func RunTags() error {
	m := NewTagModel()
	p := tea.NewProgram(m, tea.WithAltScreen())
	_, err := p.Run()
	return err
}
