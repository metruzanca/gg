package tui

import (
	tea "github.com/charmbracelet/bubbletea"
	ggit "github.com/metru/gg/internal/git"
)

func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return m, nil

	case tea.KeyMsg:
		return m.handleKey(msg)
	}

	return m, nil
}

func (m *Model) handleKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	k := msg.String()

	if k == "q" || k == "ctrl+c" {
		if m.modal != ModalNone {
			m.modal = ModalNone
			m.refreshStatus()
			return m, nil
		}
		m.quitting = true
		return m, tea.Quit
	}

	if k == "?" {
		if m.modal == ModalHelp {
			m.modal = ModalNone
			m.refreshStatus()
			return m, nil
		}
		m.modal = ModalHelp
		return m, nil
	}

	if m.modal == ModalHelp {
		m.modal = ModalNone
		m.refreshStatus()
		return m, nil
	}

	switch m.modal {
	case ModalConfirm:
		return m.handleConfirmKeys(msg)
	case ModalCommit:
		return m.handleCommitKeys(msg)
	case ModalNone:
		return m.handleMainKeys(msg)
	}

	return m, nil
}

func (m *Model) handleMainKeys(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	k := msg.String()
	total := m.totalVisible()

	switch k {
	case "up", "k":
		if m.cursor > 0 {
			m.cursor--
		}
	case "down", "j":
		if m.cursor < total-1 {
			m.cursor++
		}
	case "1":
		if m.staged.Len() > 0 {
			m.cursor = 0
		}
	case "2":
		stagedLen := m.staged.Len()
		if m.unstaged.Len() > 0 {
			m.cursor = stagedLen
		}
	case "right", "l":
		e, ok := m.entryAtCursor()
		if ok && e.IsDir && !(m.allExpanded || m.expandedDirs[e.Path]) {
			m.expandDir(e.Path)
		}
	case "left", "h":
		e, ok := m.entryAtCursor()
		if ok && e.IsDir && (m.allExpanded || m.expandedDirs[e.Path]) {
			m.collapseDir(e.Path)
		}
	case "e":
		m.toggleExpandAll()
	case "a":
		m.stageAll()
	case "s":
		m.stageItem()
	case "enter":
		stagedPaths := ggit.CollectFiles(m.staged.tree)
		if len(stagedPaths) > 0 {
			m.openCommitModal()
			return m, m.typeInput.Focus()
		}
	}
	return m, nil
}

func (m *Model) handleConfirmKeys(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "esc":
		m.confirmCancel()
	case "enter":
		if m.confirmChoice == ConfirmYes {
			m.confirmStage()
		} else {
			m.confirmCancel()
		}
	case "left", "h":
		m.confirmChoice = ConfirmYes
	case "right", "l":
		m.confirmChoice = ConfirmNo
	}
	return m, nil
}

func (m *Model) handleCommitKeys(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	k := msg.String()

	switch k {
	case "esc":
		m.modal = ModalNone
		m.commitErr = ""
		m.refreshStatus()
		return m, nil

	case "enter":
		switch m.commitField {
		case fieldType:
			m.commitField = fieldTitle
			return m, m.titleInput.Focus()
		case fieldTitle:
			if m.titleInput.Value() != "" {
				m.commit()
			}
			return m, nil
		case fieldDescription:
			m.commit()
			return m, nil
		}
		return m, nil

	case "tab":
		m.commitField = (m.commitField + 1) % 3
		switch m.commitField {
		case fieldType:
			m.descInput.Blur()
			return m, m.typeInput.Focus()
		case fieldTitle:
			return m, m.titleInput.Focus()
		case fieldDescription:
			m.titleInput.Blur()
			return m, m.descInput.Focus()
		}
		return m, nil

	case "up":
		if m.commitField == fieldType && len(m.filteredTypes) > 0 {
			m.typeSuggestIdx--
			if m.typeSuggestIdx < 0 {
				m.typeSuggestIdx = len(m.filteredTypes) - 1
			}
			m.typeInput.SetValue(m.filteredTypes[m.typeSuggestIdx].Name)
			m.updateTypeSuggestions()
		}
		return m, nil

	case "down":
		if m.commitField == fieldType && len(m.filteredTypes) > 0 {
			m.typeSuggestIdx++
			if m.typeSuggestIdx >= len(m.filteredTypes) {
				m.typeSuggestIdx = 0
			}
			m.typeInput.SetValue(m.filteredTypes[m.typeSuggestIdx].Name)
			m.updateTypeSuggestions()
		}
		return m, nil
	}

	switch m.commitField {
	case fieldType:
		var cmd tea.Cmd
		m.typeInput, cmd = m.typeInput.Update(msg)
		m.updateTypeSuggestions()
		return m, cmd
	case fieldTitle:
		var cmd tea.Cmd
		m.titleInput, cmd = m.titleInput.Update(msg)
		return m, cmd
	case fieldDescription:
		var cmd tea.Cmd
		m.descInput, cmd = m.descInput.Update(msg)
		return m, cmd
	}

	return m, nil
}
