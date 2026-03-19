package args

import (
	"os"
)

const DefaultTemplate = "Simple"

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
		template = DefaultTemplate
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
