package ui

import (
	"testing"

	tea "github.com/charmbracelet/bubbletea"
)

type dummyField struct{}

// Hidable
func (d *dummyField) SetHide(f func() bool) {}
func (d *dummyField) IsHidden() bool        { return false }

// Field interface (runtime methods only)
func (d *dummyField) getTitle() string           { return "dummy" }
func (d *dummyField) focus() tea.Cmd             { return nil }
func (d *dummyField) blur()                      {}
func (d *dummyField) update(msg tea.Msg) tea.Cmd { return nil }
func (d *dummyField) isFocused() bool            { return true }
func (d *dummyField) render() string             { return "Hello\nWorld\n" }

func TestModelLineCount(t *testing.T) {
	m, err := newModel(&dummyField{})
	if err != nil {
		t.Fatalf("failed to create model: %v", err)
	}
	_ = m.View()
	var ui Form = m
	if ui.LineCount() < 2 {
		t.Errorf("expected at least 2 lines, got %d", ui.LineCount())
	}
}
