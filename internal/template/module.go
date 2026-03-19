package template

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"golang.org/x/mod/modfile"
)

const (
	_gomodFileName    = "go.mod"
	_mainGoFileName   = "main.go"
	_mainTestFileName = "main_test.go"
	_makefileName     = "Makefile"
	_mainGoContent    = `package main

import "fmt"

func main() {
	fmt.Println("Hello, World!")
}`
	_mainTestGoContent = `package main

import "testing"

func TestMain(t *testing.T) {
	t.Log("Hello, World!")
}`
)

// dynamicMakefileContent returns a Makefile string for the given main.go relative path.
func dynamicMakefileContent(mainGoRelPath string) string {
	return strings.Join(
		[]string{
			".PHONY: build test",
			"build:\n\tgo build " + mainGoRelPath,
			"test:\n\tgo test ./...",
			"run:\n\tgo run " + mainGoRelPath,
		},
		"\n\n",
	)
}

// writeFileAt writes content to dir/filename, creating the directory if needed.
// If dir already ends with filename, it writes directly to dir instead.
func writeFileAt(dir, filename string, content []byte) error {
	dir = filepath.Clean(dir)
	filePath := filepath.Join(dir, filename)
	if filepath.Base(dir) == filename {
		filePath = dir
	}

	if err := ensureDir(filepath.Dir(filePath)); err != nil {
		return fmt.Errorf("could not create directory for %s: %w", filename, err)
	}
	if err := os.WriteFile(filePath, content, 0o644); err != nil {
		return fmt.Errorf("failed to write %s: %w", filename, err)
	}
	return nil
}

type gomodData []byte // go.mod file contents

func newGoMod(moduleName string) (gomodData, error) {
	modFile := &modfile.File{}

	if err := modFile.AddModuleStmt(moduleName); err != nil {
		return nil, fmt.Errorf("could not add module statement: %w", err)
	}

	if err := modFile.AddGoStmt(goVersion()); err != nil {
		return nil, fmt.Errorf("failed to add Go version directive: %w", err)
	}

	modData, err := modFile.Format()
	if err != nil {
		return nil, fmt.Errorf("failed to format modfile: %w", err)
	}

	return modData, nil
}

func (gmd gomodData) WriteToFile(dir string) error {
	return writeFileAt(dir, _gomodFileName, gmd)
}

// Options provides the required project creation arguments.
type Options interface {
	Name() string
	Path() string
	Template() string
	GitRepo() string
	CreateMakefile() bool
	InitGit() bool
}

func CreateNewModule(opts Options) error {
	if err := simple(opts); err != nil {
		return err
	}
	if opts.CreateMakefile() {
		if err := newMakefile(opts.Path()); err != nil {
			return err
		}
	}
	return nil
}

func CreateNewModuleWithTest(opts Options) error {
	if err := ensureDir(opts.Path()); err != nil {
		return err
	}

	gomod, err := newGoMod(opts.Name())
	if err != nil {
		return err
	}
	if err := gomod.WriteToFile(opts.Path()); err != nil {
		return err
	}

	basename := filepath.Base(opts.Name())
	cmdPath := filepath.Join(opts.Path(), "cmd", basename)
	mainGoPath := filepath.Join(cmdPath, _mainGoFileName)
	if err := newMainGoAt(mainGoPath); err != nil {
		return err
	}
	if err := writeFileAt(cmdPath, _mainTestFileName, []byte(_mainTestGoContent)); err != nil {
		return err
	}

	if opts.CreateMakefile() {
		if err := newMakefile(opts.Path()); err != nil {
			return err
		}
	}

	return nil
}

// simple creates a new Go module with a main.go file at the project root.
func simple(opts Options) error {
	if err := ensureDir(opts.Path()); err != nil {
		return err
	}

	gomod, err := newGoMod(opts.Name())
	if err != nil {
		return err
	}
	if err := gomod.WriteToFile(opts.Path()); err != nil {
		return err
	}

	mainGoPath := filepath.Join(opts.Path(), _mainGoFileName)
	return newMainGoAt(mainGoPath)
}

// newMainGoAt creates a main.go file at the specified absolute path.
func newMainGoAt(absPath string) error {
	return writeFileAt(filepath.Dir(absPath), filepath.Base(absPath), []byte(_mainGoContent))
}

// ensureDir creates the directory if it does not exist.
func ensureDir(path string) error {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return os.MkdirAll(path, 0o755)
	}
	return nil
}

// goVersion returns the current Go version string (e.g., 1.21.0).
func goVersion() string {
	version := runtime.Version()
	if len(version) > 2 && version[:2] == "go" {
		version = version[2:] // Strip "go" prefix
	}
	return version
}

// ReplaceModuleName updates the module name in go.mod at the given path.
func ReplaceModuleName(dir, newName string) error {
	goModPath := filepath.Join(dir, _gomodFileName)
	goModBytes, err := os.ReadFile(goModPath)
	if err != nil {
		return fmt.Errorf("could not read go.mod file: %w", err)
	}

	modFile, err := modfile.Parse(goModPath, goModBytes, nil)
	if err != nil {
		return fmt.Errorf("could not parse go.mod file: %w", err)
	}

	if err := modFile.AddModuleStmt(newName); err != nil {
		return fmt.Errorf("could not set module name: %w", err)
	}

	modData, err := modFile.Format()
	if err != nil {
		return fmt.Errorf("failed to format modfile: %w", err)
	}

	if err := os.WriteFile(goModPath, modData, 0o644); err != nil {
		return fmt.Errorf("failed to write go.mod file: %w", err)
	}

	return nil
}

func newMakefile(projectPath string) error {
	// Determine main.go location for Makefile
	mainGoRelPath := _mainGoFileName // default for simple
	cmdDir := filepath.Join(projectPath, "cmd")
	if _, err := os.Stat(cmdDir); err == nil {
		entries, _ := os.ReadDir(cmdDir)
		if len(entries) > 0 {
			modDir := entries[0].Name()
			candidate := filepath.Join("cmd", modDir, _mainGoFileName)
			if _, err := os.Stat(filepath.Join(projectPath, candidate)); err == nil {
				mainGoRelPath = candidate
			}
		}
	}
	makefileContent := dynamicMakefileContent(mainGoRelPath)

	return writeFileAt(projectPath, _makefileName, []byte(makefileContent))
}
