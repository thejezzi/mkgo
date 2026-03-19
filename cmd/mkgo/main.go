package main

import (
	"context"
	"fmt"
	"os"
	"strings"

	cli "github.com/urfave/cli/v3"

	argsPkg "github.com/thejezzi/mkgo/internal/args"
	"github.com/thejezzi/mkgo/internal/template"
	"github.com/thejezzi/mkgo/internal/ui"
)

// newArgumentsFromUI converts a ui.Form to *argsPkg.Arguments.
func newArgumentsFromUI(f ui.Form) *argsPkg.Arguments {
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

// createProject finds the matching template and executes it.
func createProject(args *argsPkg.Arguments) error {
	var templ template.Template
	found := false
	for _, t := range template.All {
		if strings.EqualFold(t.Name, args.Template()) {
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
	cmd := &cli.Command{
		Name:      "mkgo",
		Usage:     "Scaffold a new Go project",
		ArgsUsage: "[module-name]",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "path",
				Aliases: []string{"p"},
				Usage:   "the directory to create the project in",
			},
			&cli.StringFlag{
				Name:    "template",
				Aliases: []string{"t"},
				Value:   argsPkg.DefaultTemplate,
				Usage:   "project template (Simple, Test, Git)",
			},
			&cli.StringFlag{
				Name:    "git",
				Aliases: []string{"g"},
				Usage:   "git repository URL to clone from",
			},
			&cli.BoolFlag{
				Name:    "makefile",
				Aliases: []string{"m"},
				Usage:   "create a Makefile",
			},
			&cli.BoolFlag{
				Name:  "init-git",
				Usage: "initialize a new git repository",
			},
		},
		Action: func(ctx context.Context, cmd *cli.Command) error {
			moduleName := cmd.Args().First()

			// No positional argument → launch interactive TUI
			if moduleName == "" {
				form, err := ui.NewForm()
				if err != nil {
					return err
				}
				args := newArgumentsFromUI(form)
				return createProject(args)
			}

			// CLI mode: build Arguments from flags
			args := argsPkg.NewArguments(
				moduleName,
				cmd.String("path"),
				cmd.String("template"),
				cmd.String("git"),
				cmd.Bool("makefile"),
				cmd.Bool("init-git"),
			)

			return createProject(args)
		},
	}

	if err := cmd.Run(context.Background(), os.Args); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}
