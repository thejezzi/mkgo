package template

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"strings"

	"golang.org/x/mod/modfile"
)

const (
	_gomodFileName    = "go.mod"
	_mainGoFileName   = "main.go"
	_mainTestFileName = "main_test.go"
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

type gomodData []byte // go.mod file contents

func newGoMod(moduleName string) (gomodData, error) {
	modFile := &modfile.File{}

	if err := modFile.AddModuleStmt(moduleName); err != nil {
		return nil, fmt.Errorf("could not add module statement: %v", err)
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

func (gmd gomodData) WriteToFile(path string) error {
	cleanedPath := filepath.Clean(path)
	splitted := strings.Split(cleanedPath, "/")
	if splitted[len(splitted)-1] != _gomodFileName {
		splitted = append(splitted, _gomodFileName)
	}
	gomodPath := strings.Join(splitted, "/")

	dirPath := filepath.Dir(gomodPath)
	if err := ensureDir(dirPath); err != nil {
		return fmt.Errorf("could not create main.go file: %v", err)
	}
	if err := os.WriteFile(gomodPath, gmd, 0o644); err != nil {
		return fmt.Errorf("failed to write go.mod file: %w", err)
	}
	return nil
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

	basename := path.Base(opts.Name())
	cmdPath := filepath.Join(opts.Path(), "cmd", basename)
	mainGoPath := filepath.Join(cmdPath, _mainGoFileName)
	if err := newMainGoAt(mainGoPath); err != nil {
		return err
	}
	if err := newMainTestGo(cmdPath); err != nil {
		return err
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
	if err := newMainGoAt(mainGoPath); err != nil {
		return err
	}

	return nil
}

// newMainGoAt creates a main.go file at the specified absolute path.
func newMainGoAt(path string) error {
	dirPath := filepath.Dir(path)
	if err := ensureDir(dirPath); err != nil {
		return fmt.Errorf("could not create main.go file: %v", err)
	}
	if err := os.WriteFile(path, []byte(_mainGoContent), 0o644); err != nil {
		return fmt.Errorf("failed to write main.go file: %w", err)
	}
	return nil
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

// newMainTestGo creates a main_test.go file in the specified path.
func newMainTestGo(path string) error {
	cleanedPath := filepath.Clean(path)
	splitted := strings.Split(cleanedPath, "/")
	if splitted[len(splitted)-1] != _mainTestFileName {
		splitted = append(splitted, _mainTestFileName)
	}
	mainTestGoPath := strings.Join(splitted, "/")

	dirPath := filepath.Dir(mainTestGoPath)
	if err := ensureDir(dirPath); err != nil {
		return fmt.Errorf("could not create main_test.go file: %v", err)
	}

	if err := os.WriteFile(mainTestGoPath, []byte(_mainTestGoContent), 0o644); err != nil {
		return fmt.Errorf("failed to write main_test.go file: %w", err)
	}

	return nil
}

// ReplaceModuleName updates the module name in go.mod at the given path.
func ReplaceModuleName(path, newName string) error {
	goModPath := filepath.Join(path, _gomodFileName)
	goModBytes, err := os.ReadFile(goModPath)
	if err != nil {
		return fmt.Errorf("could not read go.mod file: %w", err)
	}

	modFile, err := modfile.Parse(goModPath, goModBytes, nil)
	if err != nil {
		return fmt.Errorf("could not parse go.mod file: %w", err)
	}

	if err := modFile.AddModuleStmt(newName); err != nil {
		return fmt.Errorf("could not add module statement: %v", err)
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

// newMainGo creates a main.go file in the specified path.
// Deprecated: use newMainGoAt for root main.go creation.
func newMainGo(path string) error {
	return newMainGoAt(path)
}

func newMakefile(path string) error {
	cleanedPath := filepath.Clean(path)
	splitted := strings.Split(cleanedPath, "/")
	if splitted[len(splitted)-1] != "Makefile" {
		splitted = append(splitted, "Makefile")
	}
	makefilePath := strings.Join(splitted, "/")

	dirPath := filepath.Dir(makefilePath)
	if err := ensureDir(dirPath); err != nil {
		return fmt.Errorf("could not create Makefile: %v", err)
	}

	// Determine main.go location for Makefile
	mainGoRelPath := "main.go" // default for simple
	if _, err := os.Stat(filepath.Join(path, "cmd")); err == nil {
		// If cmd/<modulename>/main.go exists, use that
		entries, _ := os.ReadDir(filepath.Join(path, "cmd"))
		if len(entries) > 0 {
			modDir := entries[0].Name()
			candidate := filepath.Join("cmd", modDir, "main.go")
			if _, err := os.Stat(filepath.Join(path, candidate)); err == nil {
				mainGoRelPath = candidate
			}
		}
	}
	makefileContent := dynamicMakefileContent(mainGoRelPath)

	if err := os.WriteFile(makefilePath, []byte(makefileContent), 0o644); err != nil {
		return fmt.Errorf("failed to write Makefile: %w", err)
	}

	return nil
}
