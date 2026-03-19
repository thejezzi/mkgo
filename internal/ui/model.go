package ui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/cursor"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

// Form is the interface for interactive forms, providing summary and line count access.
type Form interface {
	Summary() string
	LineCount() int
	GetModule() string
	GetPath() string
	GetTemplate() string
	GetGitRepo() string
	GetCreateMakefile() bool
	GetInitGit() bool
}

type model struct {
	focusIndex int
	fields     []Field
	cursorMode cursor.Mode
	aborted    bool
	lineCount  int // dynamically updated line count

	module         string
	path           string
	template       string
	gitRepo        string
	createMakefile bool
	initGit        bool
}

func newModel(fields ...Field) (*model, error) {
	m := model{
		fields: make([]Field, len(fields)),
	}

	for i, field := range fields {
		if i == 0 {
			field.focus()
		}
		m.fields[i] = field
	}

	return &m, nil
}

func (m *model) Cleanup() {
	lines := m.LineCount()
	if lines > 0 {
		for i := 1; i < lines; i++ {
			// Move cursor up and clear line
			fmt.Print("\033[2K\033[A")
		}
		// Clear the final line (where the cursor ends up)
		fmt.Print("\033[2K\r")
	}
}

func (m *model) Init() tea.Cmd {
	return textinput.Blink
}

func (m *model) setAllCursorsBlink() {
	for _, field := range m.fields {
		input, ok := field.(*inputModel)
		if !ok {
			continue
		}
		input.SetInnerCursorMode(cursor.CursorBlink)
	}
}

func (m *model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if _, ok := msg.(tea.KeyMsg); !ok {
		return m.handleButKeyMsg(msg)
	}

	keyMsg := msg.(tea.KeyMsg)
	switch keyMsg.String() {
	case "ctrl+c":
		m.aborted = true
		return m, tea.Quit

	case "ctrl+r":
		if m.focusIndex == len(m.fields) {
			// Focus is on submit, ignore rotation
			break
		}
		input, ok := m.fields[m.focusIndex].(*inputModel)
		if !ok {
			break
		}
		if input.disablePromptRotation {
			break
		}
		input.rotatePrompt()
		cmd := m.updateFields(msg)
		return m, cmd

	// Set focus to next input
	case "tab", "shift+tab", "up", "down":
		return m.focusNext(keyMsg)
	case "enter":
		if m.focusIndex == len(m.fields) {
			return m, tea.Quit
		}
	}

	cmd := m.updateFields(msg)
	return m, cmd
}

func (m *model) handleButKeyMsg(msg tea.Msg) (*model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		for _, field := range m.fields {
			if list, ok := field.(*listModel); ok {
				list.setWidth(msg.Width)
			}
		}
	default:
	}

	return m, m.updateFields(msg)
}

func (m *model) focusNext(msg tea.KeyMsg) (*model, tea.Cmd) {
	s := msg.String()
	move := func(forward bool) {
		for {
			if forward {
				m.focusIndex++
				if m.focusIndex > len(m.fields) {
					m.focusIndex = 0
				}
			} else {
				m.focusIndex--
				if m.focusIndex < 0 {
					m.focusIndex = len(m.fields)
				}
			}
			// If we\'re at the submit button, stop
			if m.focusIndex == len(m.fields) {
				break
			}
			// Skip hidden fields and headers
			if m.fields[m.focusIndex].IsHidden() {
				continue
			}
			if _, ok := m.fields[m.focusIndex].(*headerModel); ok {
				continue
			}
			break
		}
	}

	if s == "up" || s == "shift+tab" {
		move(false)
	} else {
		m.setAllCursorsBlink()
		move(true)
	}

	return m, tea.Batch(m.evaluateFocusStyles()...)
}

func (m *model) evaluateFocusStyles() []tea.Cmd {
	cmds := make([]tea.Cmd, len(m.fields))
	for i, field := range m.fields {
		if i == m.focusIndex {
			cmds[i] = field.focus()
			continue
		}
		// Remove focused state
		field.blur()
	}
	return cmds
}

// updateFields updates the inner textinput elements and nothing more.
func (m *model) updateFields(msg tea.Msg) tea.Cmd {
	cmds := make([]tea.Cmd, len(m.fields))
	for i, field := range m.fields {
		if field.isFocused() {
			cmds[i] = field.update(msg)
		}
	}
	return tea.Batch(cmds...)
}

const mainTitle = `
╔═════════════════════════════════════╗
║███╗   ███╗██╗  ██╗ ██████╗  ██████╗ ║
║████╗ ████║██║ ██╔╝██╔════╝ ██╔═══██╗║
║██╔████╔██║█████╔╝ ██║  ███╗██║   ██║║
║██║╚██╔╝██║██╔═██╗ ██║   ██║██║   ██║║
║██║ ╚═╝ ██║██║  ██╗╚██████╔╝╚██████╔╝║
║╚═╝     ╚═╝╚═╝  ╚═╝ ╚═════╝  ╚═════╝ ║
╚═════════════════════════════════════╝
`

const mainTitle2 = `
▗▖  ▗▖▗▖ ▗▖ ▗▄▄▖ ▄▄▄  
▐▛▚▞▜▌▐▌▗▞▘▐▌   █   █ 
▐▌  ▐▌▐▛▚▖ ▐▌▝▜▌▀▄▄▄▀ 
▐▌  ▐▌▐▌ ▐▌▝▚▄▞▘      
`

func (m *model) View() string {
	var b strings.Builder
	b.WriteString(TitleStyle.Render(strings.Trim(mainTitle2, "\n")) + "\n\n")

	for _, field := range m.fields {
		if field.IsHidden() {
			continue
		}
		b.WriteString(field.render())
	}

	button := &blurredButton
	if m.focusIndex == len(m.fields) {
		button = &focusedButton
	}
	fmt.Fprintf(&b, "\n%s\n\n", *button)
	if m.focusIndex < len(m.fields) {
		// check if any of the fields has prompt rotation enabled
		focusedInput, ok := m.fields[m.focusIndex].(*inputModel)
		if ok && !focusedInput.disablePromptRotation && focusedInput.rotationDescription != "" {
			b.WriteString(HelpStyle.Render("ctrl+r to change the " + focusedInput.rotationDescription))
		}
	}

	rendered := AppStyle.Render(b.String())
	// Count lines: split on '\n', count non-empty lines
	lines := strings.Count(rendered, "\n")
	if !strings.HasSuffix(rendered, "\n") && len(rendered) > 0 {
		lines++
	}
	m.lineCount = lines
	return rendered
}

// Fields returns all fields in the form.
func (m *model) Fields() []Field {
	return m.fields
}

// GetModule returns the value of the Module field.
func (m *model) GetModule() string {
	return m.module
}

// GetPath returns the value of the Path field.
func (m *model) GetPath() string {
	return m.path
}

// GetTemplate returns the value of the Template field.
func (m *model) GetTemplate() string {
	return m.template
}

// GetGitRepo returns the value of the Git Repository field.
func (m *model) GetGitRepo() string {
	return m.gitRepo
}

// GetCreateMakefile returns the value of the Create a Makefile checkbox.
func (m *model) GetCreateMakefile() bool {
	return m.createMakefile
}

// GetInitGit returns the value of the Initialize git repository checkbox.
func (m *model) GetInitGit() bool {
	return m.initGit
}

// LineCount returns the number of lines currently rendered by the UI.
func (m *model) LineCount() int {
	return m.lineCount
}

// Summary returns a styled summary string for the model fields.
func (m *model) Summary() string {
	type entry struct{ title, value string }
	var entries []entry
	maxTitleLen := 0

	for _, f := range m.fields {
		if f.IsHidden() {
			continue
		}
		var sumTitle, sumValue string
		if s, ok := f.(interface{ summaryEntry() (string, string) }); ok {
			sumTitle, sumValue = s.summaryEntry()
		}
		if sumTitle != "" {
			if len(sumTitle) > maxTitleLen {
				maxTitleLen = len(sumTitle)
			}
			entries = append(entries, entry{sumTitle, sumValue})
		}
	}

	var b strings.Builder
	b.WriteString(TitleStyle.Render("Summary:") + "\n\n")
	for _, entry := range entries {
		paddedTitle := fmt.Sprintf("%-*s", maxTitleLen, entry.title)
		b.WriteString(HelpStyle.Render(paddedTitle+": ") + entry.value + "\n")
	}
	return b.String()
}
