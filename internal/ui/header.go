package ui

import (
	tea "github.com/charmbracelet/bubbletea"
)

func Header() *headerModel {
	return &headerModel{}
}

type headerModel struct {
	title  string
	hidden func() bool
}

func (hm *headerModel) Title(s string) *headerModel {
	hm.title = s
	return hm
}

func (hm *headerModel) getTitle() string {
	return hm.title
}

func (hm *headerModel) focus() tea.Cmd {
	return nil
}

func (hm *headerModel) blur() {
}

func (hm *headerModel) update(msg tea.Msg) tea.Cmd {
	return nil
}

func (hm *headerModel) isFocused() bool {
	return false
}

func (hm *headerModel) render() string {
	return TitleStyle.Render(hm.title) + "\n"
}

func (hm *headerModel) SetHide(hide func() bool) {
	hm.hidden = hide
}

func (hm *headerModel) IsHidden() bool {
	if hm.hidden == nil {
		return false
	}
	return hm.hidden()
}
