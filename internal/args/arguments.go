package args

import (
	"errors"
	"flag"
	"os"
)

var (
	errModuleNameCannotBeZero = errors.New("module name cannot be empty")
	errPathCannotBeZero       = errors.New("path cannot be empty")
	errUiMode                 = errors.New("no flags were provided")
)

const _defaultTemplate = "simple"

type Arguments struct {
	name           string
	path           string
	template       string
	gitRepo        string
	createMakefile bool
	initGit        bool
}

func NewArguments(moduleName, projectPath, template, gitRepo string, createMakefile, initGit bool) *Arguments {
	if len(projectPath) == 0 {
		projectPath = moduleName
	}
	if len(template) == 0 {
		template = _defaultTemplate
	}
	return &Arguments{
		name:           moduleName,
		path:           projectPath,
		template:       template,
		gitRepo:        gitRepo,
		createMakefile: createMakefile,
		initGit:        initGit,
	}
}

// Flags parses all flags and returns a struct with all possible arguments or
// and error that indicates to use the ui mode
func Flags() (*Arguments, error) {
	path := flag.String(
		"path",
		"",
		"the path to put all the files",
	)

	template := flag.String(
		"template",
		"Simple",
		"specify a template to avoid some boilerplate setup",
	)

	gitRepo := flag.String(
		"git",
		"",
		"specify a git repository to clone from",
	)
	createMakefile := flag.Bool(
		"makefile",
		false,
		"create a Makefile",
	)
	initGit := flag.Bool(
		"init-git",
		false,
		"initialize a new git repository (default: true)",
	)
	flag.Parse()

	name := flag.Arg(0)
	if len(name) == 0 {
		return nil, errors.New("the module needs a name")
	}

	return NewArguments(name, *path, *template, *gitRepo, *createMakefile, *initGit).validate()
}

// validate make sure that all arguments are set to create the project
func (args *Arguments) validate() (*Arguments, error) {
	if len(os.Args) < 2 {
		return nil, errUiMode
	}
	if len(args.name) == 0 {
		return nil, errModuleNameCannotBeZero
	}
	if len(args.path) == 0 {
		return nil, errPathCannotBeZero
	}

	return args, nil
}

func (a *Arguments) Name() string { return a.name }
func (a *Arguments) Path() string {
	return expandHome(a.path)
}

// expandHome expands ~ or ~/ to the user's home directory.
func expandHome(path string) string {
	if path == "~" {
		home, err := os.UserHomeDir()
		if err == nil {
			return home
		}
		return path
	}
	if len(path) > 2 && path[:2] == "~/" {
		home, err := os.UserHomeDir()
		if err == nil {
			return home + path[1:]
		}
		return path
	}
	return path
}
func (a *Arguments) Template() string     { return a.template }
func (a *Arguments) CreateMakefile() bool { return a.createMakefile }
func (a *Arguments) GitRepo() string      { return a.gitRepo }
func (a *Arguments) InitGit() bool        { return a.initGit }

func (args Arguments) IsEmpty() bool {
	return len(args.name) == 0 && len(args.path) == 0
}
