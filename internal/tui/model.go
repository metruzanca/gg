package tui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/metruzanca/gg/internal/config"
	ggit "github.com/metruzanca/gg/internal/git"
)

type ModalType int

const (
	ModalNone ModalType = iota
	ModalConfirm
	ModalCommit
	ModalHelp
)

type ConfirmChoice int

const (
	ConfirmYes ConfirmChoice = iota
	ConfirmNo
)

type commitFormField int

const (
	fieldType commitFormField = iota
	fieldTitle
	fieldDescription
)

type Model struct {
	staged       *Panel
	unstaged     *Panel
	cursor       int
	expandedDirs map[string]bool
	allExpanded  bool

	commitTypes []config.CommitType
	err         error
	width       int
	height      int
	quitting    bool

	modal ModalType

	confirmChoice ConfirmChoice
	confirmMsg    string
	confirmPaths  []string

	pendingSectionStaged bool
	pendingIndex         int

	typeInput      textinput.Model
	titleInput     textinput.Model
	descInput      textarea.Model
	commitField    commitFormField
	filteredTypes  []config.CommitType
	typeSuggestIdx int
	commitErr      string
}

func New() *Model {
	ti := textinput.New()
	ti.Placeholder = "type (e.g. fix)"
	ti.Width = 40
	ti.SetValue("fix")

	txi := textinput.New()
	txi.Placeholder = "commit title"
	txi.Width = 60

	da := textarea.New()
	da.Placeholder = "optional description"
	da.SetWidth(60)
	da.SetHeight(3)

	return &Model{
		staged:       NewPanel("Staged"),
		unstaged:     NewPanel("Unstaged"),
		modal:        ModalNone,
		expandedDirs: make(map[string]bool),
		typeInput:    ti,
		titleInput:   txi,
		descInput:    da,
		commitField:  fieldType,
	}
}

func (m *Model) Init() tea.Cmd {
	m.refreshStatus()
	m.commitTypes = config.DefaultTypes
	cfg, _ := config.Load("")
	if cfg != nil && len(cfg.Types) > 0 {
		m.commitTypes = cfg.Types
	}
	m.filteredTypes = m.commitTypes
	m.updateTypeSuggestions()
	return textinput.Blink
}

func (m *Model) refreshStatus() {
	files, err := ggit.Status()
	if err != nil {
		m.err = err
		m.staged.SetTree(nil)
		m.unstaged.SetTree(nil)
		return
	}
	m.err = nil

	stagedTree, unstagedTree := ggit.BuildTree(files)
	m.staged.SetTree(stagedTree)
	m.unstaged.SetTree(unstagedTree)
	m.rebuildAll()
}

func (m *Model) rebuildAll() {
	oldPath := ""
	e, ok := m.entryAtCursor()
	if ok {
		oldPath = e.Path
	}

	m.staged.Rebuild(m.expandedDirs, m.allExpanded)
	m.unstaged.Rebuild(m.expandedDirs, m.allExpanded)

	if oldPath != "" {
		newIdx := m.findEntryPath(oldPath)
		if newIdx >= 0 {
			m.cursor = newIdx
		}
	} else {
		total := m.staged.Len() + m.unstaged.Len()
		if m.cursor >= total {
			m.cursor = max(0, total-1)
		}
	}
}

func (m *Model) findEntryPath(path string) int {
	for i := 0; i < m.staged.Len(); i++ {
		e, _ := m.staged.EntryAt(i)
		if e.Path == path {
			return i
		}
	}
	offs := m.staged.Len()
	for i := 0; i < m.unstaged.Len(); i++ {
		e, _ := m.unstaged.EntryAt(i)
		if e.Path == path {
			return offs + i
		}
	}
	for i := 0; i < m.staged.Len(); i++ {
		e, _ := m.staged.EntryAt(i)
		if e.IsDir && strings.HasPrefix(path, e.Path+"/") {
			return i
		}
	}
	for i := 0; i < m.unstaged.Len(); i++ {
		e, _ := m.unstaged.EntryAt(i)
		if e.IsDir && strings.HasPrefix(path, e.Path+"/") {
			return offs + i
		}
	}
	return -1
}

func (m *Model) totalVisible() int {
	return m.staged.Len() + m.unstaged.Len()
}

func (m *Model) cursorInStaged() bool {
	return m.cursor < m.staged.Len()
}

func (m *Model) panelIndex() (staged bool, idx int) {
	if m.cursorInStaged() {
		return true, m.cursor
	}
	return false, m.cursor - m.staged.Len()
}

func (m *Model) moveCursorAfterChange(staged bool, idx int) {
	if staged {
		if m.staged.Len() > 0 {
			m.cursor = min(idx, m.staged.Len()-1)
		} else {
			m.cursor = 0
		}
		return
	}
	offs := m.staged.Len()
	if m.unstaged.Len() > 0 {
		m.cursor = offs + min(idx, m.unstaged.Len()-1)
	} else {
		m.cursor = max(0, offs-1)
	}
}

func (m *Model) entryAtCursor() (*ggit.Entry, bool) {
	if m.cursorInStaged() {
		return m.staged.EntryAt(m.cursor)
	}
	return m.unstaged.EntryAt(m.cursor - m.staged.Len())
}

func (m *Model) updateTypeSuggestions() {
	val := strings.ToLower(m.typeInput.Value())
	var matches []config.CommitType
	for _, t := range m.commitTypes {
		if strings.HasPrefix(strings.ToLower(t.Name), val) {
			matches = append(matches, t)
		}
	}
	if len(matches) == 0 {
		matches = m.commitTypes
	}
	m.filteredTypes = matches
	if m.typeSuggestIdx >= len(m.filteredTypes) {
		m.typeSuggestIdx = 0
	}
}

func (m *Model) commit() {
	if m.typeInput.Value() == "" || m.titleInput.Value() == "" {
		m.commitErr = "type and title are required"
		return
	}
	spec := ggit.CommitSpec{
		Type:        m.typeInput.Value(),
		Title:       m.titleInput.Value(),
		Description: m.descInput.Value(),
	}
	if err := ggit.Commit(spec); err != nil {
		m.commitErr = fmt.Sprintf("commit failed: %v", err)
		return
	}
	m.modal = ModalNone
	m.commitErr = ""
	m.typeInput.SetValue("fix")
	m.titleInput.SetValue("")
	m.descInput.Reset()
	m.commitField = fieldType
	m.refreshStatus()
}

func (m *Model) stageItem() {
	e, ok := m.entryAtCursor()
	if !ok {
		return
	}
	inStaged, idx := m.panelIndex()

	var paths []string
	if e.IsDir {
		paths = ggit.CollectFiles([]*ggit.Entry{e})
	} else {
		paths = []string{e.Path}
	}

	if inStaged {
		for _, p := range paths {
			if err := ggit.Unstage(p); err != nil {
				m.err = err
				return
			}
		}
	} else {
		if m.stagedHasAnyFile(paths) {
			msg := fmt.Sprintf("File %q already has staged changes.\nStage remaining hunks?", e.Path)
			if e.IsDir {
				msg = fmt.Sprintf("Folder %q contains files with staged changes.\nStage remaining hunks?", e.Path)
			}
			m.showConfirm(msg, paths)
			return
		}
		for _, p := range paths {
			if err := ggit.Stage(p); err != nil {
				m.err = err
				return
			}
		}
	}

	m.refreshStatus()
	m.moveCursorAfterChange(inStaged, idx)
}

func (m *Model) stagedHasFile(path string) bool {
	n := m.staged.Len()
	for i := 0; i < n; i++ {
		e, _ := m.staged.EntryAt(i)
		if !e.IsDir && e.Path == path {
			return true
		}
	}
	return false
}

func (m *Model) stagedHasAnyFile(paths []string) bool {
	for _, p := range paths {
		if m.stagedHasFile(p) {
			return true
		}
	}
	return false
}

func (m *Model) stageAll() {
	paths := ggit.CollectFiles(m.unstaged.tree)
	if len(paths) == 0 {
		return
	}
	inStaged, idx := m.panelIndex()
	stagedPaths := ggit.CollectFiles(m.staged.tree)
	if len(stagedPaths) > 0 {
		m.showConfirm("Some files are already staged.\nStage all remaining files?", paths)
		return
	}
	for _, p := range paths {
		ggit.Stage(p)
	}
	m.refreshStatus()
	m.moveCursorAfterChange(inStaged, idx)
}

func (m *Model) showConfirm(msg string, paths []string) {
	m.confirmMsg = msg
	m.confirmPaths = paths
	m.confirmChoice = ConfirmNo
	m.pendingSectionStaged, m.pendingIndex = m.panelIndex()
	m.modal = ModalConfirm
}

func (m *Model) confirmStage() {
	for _, p := range m.confirmPaths {
		ggit.Stage(p)
	}
	m.modal = ModalNone
	m.confirmPaths = nil
	m.refreshStatus()
	m.moveCursorAfterChange(m.pendingSectionStaged, m.pendingIndex)
}

func (m *Model) confirmCancel() {
	m.modal = ModalNone
	m.confirmPaths = nil
}

func (m *Model) expandDir(path string) {
	if m.allExpanded {
		m.allExpanded = false
		m.markAllDirsExpanded()
	}
	m.expandedDirs[path] = true
	m.rebuildAll()
}

func (m *Model) collapseDir(path string) {
	if m.allExpanded {
		m.allExpanded = false
		m.expandedDirs = make(map[string]bool)
		m.markAllDirsExpanded()
	}
	delete(m.expandedDirs, path)
	m.rebuildAll()
}

func (m *Model) toggleExpandAll() {
	if m.allExpanded {
		m.allExpanded = false
		m.expandedDirs = make(map[string]bool)
	} else {
		m.allExpanded = true
		m.expandedDirs = make(map[string]bool)
	}
	m.rebuildAll()
}

func (m *Model) markAllDirsExpanded() {
	for _, p := range ggit.AllDirPaths(m.staged.tree) {
		m.expandedDirs[p] = true
	}
	for _, p := range ggit.AllDirPaths(m.unstaged.tree) {
		m.expandedDirs[p] = true
	}
}

func (m *Model) openCommitModal() {
	ti := textinput.New()
	ti.Placeholder = "type (e.g. fix)"
	ti.Width = 40
	ti.SetValue("fix")
	ti.Focus()

	txi := textinput.New()
	txi.Placeholder = "commit title"
	txi.Width = 60

	da := textarea.New()
	da.Placeholder = "optional description"
	da.SetWidth(60)
	da.SetHeight(3)

	m.typeInput = ti
	m.titleInput = txi
	m.descInput = da
	m.commitField = fieldType
	m.filteredTypes = m.commitTypes
	m.typeSuggestIdx = 0
	m.commitErr = ""
	m.updateTypeSuggestions()
	m.modal = ModalCommit
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func Run() error {
	m := New()
	p := tea.NewProgram(m, tea.WithAltScreen())
	_, err := p.Run()
	return err
}
