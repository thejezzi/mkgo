package main

import (
	"fmt"
	"os"

	argsPkg "github.com/thejezzi/mkgo/internal/args"
	"github.com/thejezzi/mkgo/internal/template"
	"github.com/thejezzi/mkgo/internal/ui"
)

// newArgumentsFromUI converts a ui.UI to *argsPkg.Arguments
func newArgumentsFromUI(f ui.Form) *argsPkg.Arguments {
	// Converts a ui.Form to *argsPkg.Arguments
	if f == nil {
		return nil
	}
	return argsPkg.NewArguments(
		f.GetModule(),
		f.GetPath(),
		f.GetTemplate(),
		f.GetGitRepo(),
		f.GetCreateMakefile(),
		f.GetInitGit(),
	)
}

// getArguments returns Arguments from flags or UI
func getArguments() (*argsPkg.Arguments, error) {
	if len(os.Args) > 1 {
		return argsPkg.Flags()
	}
	form, err := ui.NewForm()
	if err != nil {
		return nil, err
	}
	return newArgumentsFromUI(form), nil
}

func run() error {
	args, err := getArguments()
	if err != nil {
		return err
	}

	var templ template.Template
	found := false
	for _, t := range template.All {
		if t.Name == args.Template() {
			found = true
			templ = t
			break
		}
	}

	if !found {
		return fmt.Errorf("template not found: %s", args.Template())
	}

	if err := templ.Create(args); err != nil {
		return err
	}

	fmt.Println(ui.SummaryFromArguments(args))
	return nil
}

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "could not create new project: %v\n", err)
		os.Exit(1)
	}
}
