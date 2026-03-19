package ui

import (
	"github.com/charmbracelet/bubbles/list"
	"github.com/thejezzi/mkgo/internal/template"
)

type FieldType int

const (
	InputType FieldType = iota
	ListType
	CheckboxType
	GroupType
)

type FieldDef struct {
	Type                  FieldType
	Title                 string
	Description           string
	RotationTitle         string
	Placeholder           string
	Prompts               []string
	Focus                 bool
	Validate              func(string) error
	Value                 *string
	DisablePromptRotation bool
	Hide                  func() bool
	CheckboxValue         *bool
	Fields                []FieldDef
}

func createInputField(fd FieldDef) Field {
	input := Input().
		Title(fd.Title).
		Description(fd.Description).
		RotationDescription(fd.RotationTitle).
		Placeholder(fd.Placeholder).
		Prompt(fd.Prompts...).
		Validate(fd.Validate).
		Value(fd.Value)
	if fd.DisablePromptRotation {
		input.DisablePromptRotation()
	}
	if fd.Focus {
		input.FocusOnStart()
	}
	input.SetHide(fd.Hide)
	return input
}

func createCheckboxField(fd FieldDef) Field {
	checkbox := Checkbox().
		SetValue(fd.CheckboxValue).
		Title(fd.Title).
		Description(fd.Description)
	checkbox.SetHide(fd.Hide)
	return checkbox
}

func createListField(fd FieldDef) Field {
	items := make([]list.Item, len(template.All))
	for i, t := range template.All {
		items[i] = t
	}
	l := List().SetItems(items...).Title(fd.Title).Value(fd.Value)
	return l
}

func buildFields(fieldDefs []FieldDef) []Field {
	var fields []Field
	for _, fd := range fieldDefs {
		switch fd.Type {
		case GroupType:
			fields = append(fields, createHeaderField(fd))

			childrenDefs := make([]FieldDef, len(fd.Fields))
			copy(childrenDefs, fd.Fields)
			if fd.Hide != nil {
				for i := range childrenDefs {
					childrenDefs[i].Hide = fd.Hide
				}
			}
			fields = append(fields, buildFields(childrenDefs)...)
		case ListType:
			fields = append(fields, createListField(fd))
		case CheckboxType:
			fields = append(fields, createCheckboxField(fd))
		default: // InputType
			fields = append(fields, createInputField(fd))
		}
	}
	return fields
}

func CreateForm(fieldDefs []FieldDef) (Form, error) {
	fields := buildFields(fieldDefs)
	return form(fields...)
}

func createHeaderField(fd FieldDef) Field {
	header := Header().Title(fd.Title)
	header.SetHide(fd.Hide)
	return header
}
