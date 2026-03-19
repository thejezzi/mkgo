// Package ui handles everything ui related in mkgo
package ui

import (
	"errors"
	"fmt"
	"os"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

func NewForm() (Form, error) {
	// These variables are shared between the FieldDef closures and the
	// resulting model. The FieldDef Value/CheckboxValue pointers write
	// directly into these during TUI interaction, and we copy them into
	// the returned model afterwards.
	var (
		module         string
		path           string
		tmpl           string
		gitRepo        string
		createMakefile bool
		initGit        bool
	)

	modulePrefixes := os.Getenv("MKGO_MODULE_PREFIXES")
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
			Value:                 &module,
			DisablePromptRotation: modulePrefixes == "",
		},
		{
			Type:          InputType,
			Title:         "Path",
			Description:   "The directory where your project will be created",
			Placeholder:   "projects/my-go-app",
			Prompts:       []string{"", "~/tmp/"},
			Focus:         true,
			Value:         &path,
			RotationTitle: "path prefix",
		},
		{
			Type:        ListType,
			Title:       "Template",
			Description: "Choose a template to quickly set up your project structure",
			Value:       &tmpl,
		},
		{
			Type:          InputType,
			Title:         "Git Repository",
			Description:   "Specify a Git repository to clone from (only for 'Git' template)",
			RotationTitle: "git prefix",
			Placeholder:   "github.com/user/repo",
			Prompts:       []string{"https://", "git@"},
			Value:         &gitRepo,
			Hide: func() bool {
				return tmpl != "Git"
			},
		},
		{
			Type:  GroupType,
			Title: "Additional options",
			Fields: []FieldDef{
				{
					Type:          CheckboxType,
					Description:   "Create a Makefile",
					CheckboxValue: &createMakefile,
				},
				{
					Type:          CheckboxType,
					Description:   "Initialize a git repository",
					CheckboxValue: &initGit,
				},
			},
			Hide: func() bool {
				return tmpl == "Git"
			},
		},
	}

	result, err := CreateForm(fieldDefs)
	if err != nil {
		return nil, err
	}

	// Copy the values captured by the field pointers into the model.
	if m, ok := result.(*model); ok {
		m.module = module
		m.path = path
		m.template = tmpl
		m.gitRepo = gitRepo
		m.createMakefile = createMakefile
		m.initGit = initGit
		return m, nil
	}
	return result, nil
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
