package ui

import (
	"strings"

	"github.com/charmbracelet/bubbles/cursor"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

func Input() *inputModel {
	im := newInputModel()
	return &im
}

type inputModel struct {
	title                 string
	titleStyle            lipgloss.Style
	description           string
	rotationDescription   string
	descriptionStyle      lipgloss.Style
	inner                 textinput.Model
	prompt                string
	promptIndex           int
	promptList            []string
	promptStyle           lipgloss.Style
	focusOnStart          bool
	disablePromptRotation bool

	value *string

	validation func(string) error
	hide       func() bool
}

func newInputModel() inputModel {
	im := inputModel{
		titleStyle:       TitleStyle,
		descriptionStyle: HelpStyle,
		promptStyle:      focusedStyle,
		inner:            textinput.New(),
		promptList:       make([]string, 0),
	}

	im.SetInnerCursorStyle(cursorStyle)
	im.CharLimit(256)
	im.AppendPrompts("")
	im.inner.Prompt = ""
	im.inner.Cursor.SetMode(cursor.CursorBlink)
	return im
}

func (im *inputModel) DisablePromptRotation() *inputModel {
	im.disablePromptRotation = true
	return im
}

func (im *inputModel) Title(s string) *inputModel {
	im.title = s
	return im
}

func (im *inputModel) getTitle() string {
	return im.title
}

func (im *inputModel) Description(desc string) *inputModel {
	im.description = desc
	return im
}

func (im *inputModel) RotationDescription(desc string) *inputModel {
	im.rotationDescription = desc
	return im
}

func (im *inputModel) FocusOnStart() *inputModel {
	im.focus()
	im.SetInnerCursorMode(cursor.CursorHide)
	return im
}

func (im *inputModel) Prompt(prompts ...string) *inputModel {
	im.promptList = make([]string, 0)
	im.promptList = append(im.promptList, prompts...)
	if len(prompts) > 0 {
		im.prompt = im.promptList[0]
	}
	return im
}

func (im *inputModel) Value(v *string) *inputModel {
	im.value = v
	return im
}

func (im *inputModel) CharLimit(n int) {
	im.inner.CharLimit = n
}

func (im *inputModel) SetInnerCursorMode(mode cursor.Mode) tea.Cmd {
	return im.inner.Cursor.SetMode(mode)
}

func (im *inputModel) SetInnerTextStyle(s lipgloss.Style) {
	im.inner.TextStyle = s
}

func (im *inputModel) focus() tea.Cmd {
	im.SetInnerPromptStyle(focusedStyle)
	im.SetInnerTextStyle(focusedStyle)
	return im.inner.Focus()
}

func (im *inputModel) Placeholder(p string) *inputModel {
	im.inner.Placeholder = p
	return im
}

func (im *inputModel) Validate(f func(string) error) *inputModel {
	im.inner.Validate = f
	return im
}

func (im *inputModel) rotatePrompt() {
	if len(im.promptList) == 0 {
		return
	}
	im.promptIndex++

	if im.promptIndex > len(im.promptList)-1 {
		im.promptIndex = 0
	}
	im.prompt = im.promptList[im.promptIndex]
}

func (im *inputModel) AppendPrompts(prompts ...string) {
	im.promptList = append(im.promptList, prompts...)
	im.prompt = im.promptList[im.promptIndex]
}

func (im *inputModel) SetInnerCursorStyle(s lipgloss.Style) {
	im.inner.Cursor.Style = s
}

func (im *inputModel) update(msg tea.Msg) tea.Cmd {
	updated, cmd := im.inner.Update(msg)
	*im.value = im.prompt + updated.Value()
	im.inner = updated
	return cmd
}

func (im *inputModel) SetInnerPromptStyle(s lipgloss.Style) {
	im.inner.PromptStyle = s
}

func (im *inputModel) blur() {
	im.SetInnerPromptStyle(noStyle)
	im.inner.Blur()
}

func (im *inputModel) render() string {
	b := strings.Builder{}
	// title
	b.WriteString(im.titleStyle.Render(im.title))
	b.WriteRune('\n')
	// the actual input
	im.inner.TextStyle = noStyle
	b.WriteString("> ")
	b.WriteString(im.promptStyle.Render(im.prompt))
	b.WriteString(im.inner.View())
	b.WriteRune('\n')
	if len(im.description) > 0 {
		b.WriteString(im.descriptionStyle.Render(im.description))
		b.WriteRune('\n')
	}
	b.WriteRune('\n')

	return b.String()
}

func (im *inputModel) SetHide(hide func() bool) {
	im.hide = hide
}

func (im *inputModel) IsHidden() bool {
	if im.hide == nil {
		return false
	}
	return im.hide()
}

// summaryEntry returns the summary title and value for this input field.
func (im *inputModel) summaryEntry() (title, value string) {
	if im.title == "" {
		return "", ""
	}
	val := ""
	if im.value != nil {
		val = *im.value
	}
	return im.title, val
}

func (im *inputModel) isFocused() bool {
	return im.inner.Focused()
}
