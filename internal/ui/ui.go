package ui

import (
	"errors"
	"fmt"
	"os"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

func NewForm() (Form, error) {
	m := &model{}

	modulePrefixes := os.Getenv("GOSPROUT_MODULE_PREFIXES")
	prefixes := []string{}
	if modulePrefixes != "" {
		for _, p := range strings.Split(modulePrefixes, ",") {
			if !strings.HasSuffix(p, "/") {
				p += "/"
			}
			prefixes = append(prefixes, p)
		}
	}

	fieldDefs := []FieldDef{
		{
			Type:          InputType,
			Title:         "Module",
			Description:   "The name of your Go module",
			RotationTitle: "module prefix",
			Placeholder:   "your-project",
			Prompts:       prefixes,
			Validate: func(s string) error {
				if len(s) == 0 {
					return errors.New("cannot be empty")
				}
				return nil
			},
			Value:                 &m.module,
			DisablePromptRotation: modulePrefixes == "",
		},
		{
			Type:          InputType,
			Title:         "Path",
			Description:   "The directory where your project will be created",
			Placeholder:   "projects/my-go-app",
			Prompts:       []string{"", "~/tmp/"},
			Focus:         true,
			Value:         &m.path,
			RotationTitle: "path prefix",
		},
		{
			Type:        ListType,
			Title:       "Template",
			Description: "Choose a template to quickly set up your project structure",
			Value:       &m.template,
		},
		{
			Type:          InputType,
			Title:         "Git Repository",
			Description:   "Specify a Git repository to clone from (only for 'Git' template)",
			RotationTitle: "git prefix",
			Placeholder:   "github.com/user/repo",
			Prompts:       []string{"https://", "git@"},
			Value:         &m.gitRepo,
			Hide: func() bool {
				return m.template != "Git"
			},
		},
		{
			Type:  GroupType,
			Title: "Additional options",
			Fields: []FieldDef{
				{
					Type:          CheckboxType,
					Description:   "Create a Makefile",
					CheckboxValue: &m.createMakefile,
				},
				{
					Type:          CheckboxType,
					Description:   "Initialize a git repository",
					CheckboxValue: &m.initGit,
				},
			},
			Hide: func() bool {
				return m.template == "Git"
			},
		},
	}

	ui, err := CreateForm(fieldDefs)
	if err != nil {
		return nil, err
	}

	// Copy the fields into the model for later use
	if mm, ok := ui.(*model); ok {
		mm.module = m.module
		mm.path = m.path
		mm.template = m.template
		mm.gitRepo = m.gitRepo
		mm.createMakefile = m.createMakefile
		mm.initGit = m.initGit
		return mm, nil
	}
	return ui, nil
}

func form(fields ...Field) (Form, error) {
	m, err := newModel(fields...)
	if err != nil {
		return nil, err
	}
	defer m.Cleanup()

	if _, err := tea.NewProgram(m).Run(); err != nil {
		return nil, fmt.Errorf("could not start program: %w", err)
	}
	if m.aborted {
		return m, errors.New("was aborted")
	}

	return m, nil
}

type summaryEntry interface {
	summaryEntry() (string, string)
}
