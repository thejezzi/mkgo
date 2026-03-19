package template

import (
	"fmt"
	"os/exec"

	"github.com/thejezzi/mkgo/internal/git"
)

// Template struct and methods

type Template struct {
	Name        string
	description string
	Create      func(opts Options) error
}

func New(
	name, description string,
	create func(opts Options) error,
) Template {
	return Template{
		Name:        name,
		description: description,
		Create:      create,
	}
}

func (t Template) Title() string       { return t.Name }
func (t Template) Description() string { return t.description }
func (t Template) FilterValue() string { return t.Name }

// Template creation logic

func simpleCreate(opts Options) error {
	return CreateNewModule(opts)
}

func testCreate(opts Options) error {
	if err := CreateNewModuleWithTest(opts); err != nil {
		return err
	}

	if opts.InitGit() {
		cmd := exec.Command("git", "init")
		cmd.Dir = opts.Path()
		if err := cmd.Run(); err != nil {
			return err
		}
	}

	return nil
}

func gitCreate(opts Options) error {
	if opts.GitRepo() == "" {
		return fmt.Errorf("git repository URL cannot be empty")
	}
	if err := git.Clone(opts.GitRepo(), opts.Path()); err != nil {
		return err
	}
	if err := git.Reinit(opts.Path()); err != nil {
		return err
	}
	return ReplaceModuleName(opts.Path(), opts.Name())
}

// Template definitions

var (
	Simple = New(
		"Simple",
		"A simple template with a cmd folder",
		simpleCreate,
	)
	Test = New(
		"Test",
		"A cmd folder with a main_test.go file",
		testCreate,
	)
	Git = New(
		"Git",
		"Create a project from a git repository",
		gitCreate,
	)
)

// All templates
var All = []Template{
	Simple,
	Test,
	Git,
}
