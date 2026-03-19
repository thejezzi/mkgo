package ui

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

func Checkbox() *checkboxModel {
	return newCheckboxModel()
}

type checkboxModel struct {
	title       string
	value       *bool
	focused     bool
	hidden      func() bool
	description string
}

func newCheckboxModel() *checkboxModel {
	return &checkboxModel{}
}

func (cm *checkboxModel) Title(s string) *checkboxModel {
	cm.title = s
	return cm
}

func (cm *checkboxModel) Description(s string) *checkboxModel {
	cm.description = s
	return cm
}

func (cm *checkboxModel) SetValue(v *bool) *checkboxModel {
	cm.value = v
	return cm
}

func (cm *checkboxModel) getTitle() string {
	return cm.title
}

func (cm *checkboxModel) focus() tea.Cmd {
	cm.focused = true
	return nil
}

func (cm *checkboxModel) blur() {
	cm.focused = false
}

func (cm *checkboxModel) update(msg tea.Msg) tea.Cmd {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case " ", "enter":
			if cm.value != nil {
				*cm.value = !*cm.value
			}
		}
	}
	return nil
}

// summaryEntry returns the summary title and value for this checkbox field.
func (cm *checkboxModel) summaryEntry() (title, value string) {
	t := cm.description
	if t == "" {
		t = cm.title
	}
	val := ""
	if cm.value != nil {
		if *cm.value {
			val = "Yes"
		} else {
			val = "No"
		}
	}
	return t, val
}

func (cm *checkboxModel) isFocused() bool {
	return cm.focused
}

func (cm *checkboxModel) render() string {
	var s strings.Builder

	// Checkmark
	check := "☐"
	if cm.value != nil && *cm.value {
		check = checkmarkChecked.Render("✓")
	}

	desc := cm.description

	if cm.focused {
		// Focused: checkmark stays green if checked, description is pink
		if cm.value != nil && *cm.value {
			s.WriteString(check)
			s.WriteString(focusedStyle.Render(" " + desc))
		} else {
			s.WriteString(focusedStyle.Render(check + " " + desc))
		}
	} else {
		// Not focused: checkmark green if checked, rest normal
		s.WriteString(check + " " + desc)
	}

	s.WriteString("\n")
	return s.String()
}

func (cm *checkboxModel) SetHide(hide func() bool) {
	cm.hidden = hide
}

func (cm *checkboxModel) IsHidden() bool {
	if cm.hidden == nil {
		return false
	}
	return cm.hidden()
}
