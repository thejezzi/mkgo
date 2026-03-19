package ui

import tea "github.com/charmbracelet/bubbletea"

// Field is the minimal interface for a UI form field. It covers only the
// methods that the model needs at runtime (rendering, focus management,
// event handling, and visibility). Builder/configuration methods live on
// the concrete types instead.
type Field interface {
	Hidable

	getTitle() string
	focus() tea.Cmd
	blur()
	update(tea.Msg) tea.Cmd
	isFocused() bool
	render() string
}
