package cli

import (
	"flag"
	"fmt"
	"os"
	"sort"
	"text/tabwriter"
)

var errorHandling = flag.ExitOnError

// App represents an application.
type App struct {
	// Desc is the application description.
	// It will be printed in the usage.
	Desc string

	// Default is a function that will be called
	// if the cli is called without any groups and commands.
	// By default, it will print the usage.
	Default func(app *flag.FlagSet) error

	groups map[string]*Group
	cmds   map[string]*Cmd

	fs *flag.FlagSet
}

// New will create a properly configured application.
func New() *App {
	app := App{
		groups: make(map[string]*Group),
		cmds:   make(map[string]*Cmd),

		fs: flag.NewFlagSet(os.Args[0], errorHandling),

		Default: func(app *flag.FlagSet) error {
			app.Usage()
			return nil
		},
	}

	app.fs.Usage = app.defaultUsage

	return &app
}

// Flags will return a FlagSet that can be used to add
// application level flags.
func (app *App) Flags() *flag.FlagSet {
	if app.fs == nil {
		app.fs = flag.NewFlagSet(os.Args[0], errorHandling)
		app.fs.Usage = app.defaultUsage
	}
	return app.fs
}

// AddGroup will add or replace a cli group.
// If the group is nil or if it does have a name it will panic.
func (app *App) AddGroup(group *Group) {
	if group == nil {
		panic("cannot add nil group")
	}

	if group.Name == "" {
		panic("cannot add group with empty name")
	}

	if group.cmds == nil {
		group.cmds = make(map[string]*Cmd)
	}

	if group.fs == nil {
		group.fs = flag.NewFlagSet(group.Name, errorHandling)
		group.fs.Usage = group.defaultUsage
	}

	// define group default cmd
	if group.Default == nil {
		group.Default = func(app *flag.FlagSet, group *flag.FlagSet) error {
			group.Usage()
			return nil
		}
	}

	app.groups[group.Name] = group
}

// AddCmd will add or replace a root command.
// The command will be directly attached to the application and
// is not part of any groups.
// If the command is nil or if it does have a name it will panic.
func (app *App) AddCmd(cmd *Cmd) {
	if cmd == nil {
		panic("cannot add nil command")
	}

	if cmd.Name == "" {
		panic("cannot add command with empty name")
	}

	if cmd.Run == nil {
		panic("cannot add command without run function")
	}

	if cmd.fs == nil {
		cmd.fs = flag.NewFlagSet(cmd.Name, errorHandling)
		cmd.fs.Usage = cmd.defaultUsage
	}

	if app.cmds == nil {
		app.cmds = make(map[string]*Cmd)
	}

	app.cmds[cmd.Name] = cmd
}

// Run should be called after having finished configuring the
// commands and groups, and after having defined flags.
// It is usually the last command that you will called in your main.
// It will parse command line flags and arguments and call the correct
// command.
func (app *App) Run() error {
	var groupIdx, cmdIdx int

	for i, f := range os.Args {
		if i == 0 {
			continue
		}

		if !isFlag(f) {
			_, isGroup := app.groups[f]
			_, isCmd := app.cmds[f]

			if groupIdx != 0 {
				_, isCmd = app.groups[os.Args[groupIdx]].cmds[f]
			}

			switch {
			case groupIdx == 0 && isGroup:
				groupIdx = i
			case cmdIdx == 0 && isCmd:
				cmdIdx = i
			}

			if cmdIdx != 0 {
				break
			}
		}
	}

	switch {

	// with command from group
	case groupIdx != 0 && cmdIdx != 0:
		group := app.groups[os.Args[groupIdx]]
		cmd := group.cmds[os.Args[cmdIdx]]

		app.fs.Parse(os.Args[1:groupIdx])
		group.fs.Parse(os.Args[groupIdx+1 : cmdIdx])
		if len(os.Args) > cmdIdx+1 {
			cmd.fs.Parse(os.Args[cmdIdx+1:])
		}

		if err := cmd.Run(app.fs, group.fs, cmd.fs); err != nil {
			return err
		}

	// with root command
	case groupIdx == 0 && cmdIdx != 0:
		cmd := app.cmds[os.Args[cmdIdx]]

		app.fs.Parse(os.Args[1:cmdIdx])
		if len(os.Args) > cmdIdx+1 {
			cmd.fs.Parse(os.Args[cmdIdx+1:])
		}

		if err := cmd.Run(app.fs, nil, cmd.fs); err != nil {
			return err
		}

	// with group default action
	case groupIdx != 0 && cmdIdx == 0:
		group := app.groups[os.Args[groupIdx]]

		app.fs.Parse(os.Args[1:groupIdx])
		if len(os.Args) > groupIdx+1 {
			group.fs.Parse(os.Args[groupIdx+1:])
		}

		if err := group.Default(app.fs, group.fs); err != nil {
			return err
		}

	// with app default action
	default:
		if len(os.Args) > 1 {
			app.fs.Parse(os.Args[1:])
		}

		if err := app.Default(app.fs); err != nil {
			return err
		}
	}

	return nil
}

func (app *App) defaultUsage() {
	var opt string
	if nbFlags(app.fs) > 0 {
		opt = " [OPTIONS]"
	}
	fmt.Printf("\nUsage:	%s%s COMMAND\n", os.Args[0], opt)

	if app.Desc != "" {
		fmt.Printf("\n%s\n", app.Desc)
	}

	// options
	if nbFlags(app.fs) > 0 {
		fmt.Print("\nOptions:\n")
		writer := tabwriter.NewWriter(os.Stdout, 0, 8, 2, '\t', 0)
		app.fs.VisitAll(func(f *flag.Flag) {
			fmt.Fprintf(writer, "  -%s\t%s\n", f.Name, f.Usage)
		})
		writer.Flush()
	}

	// groups
	if len(app.groups) > 0 {
		fmt.Print("\nManagement Commands:\n")

		keys := make([]string, 0, len(app.groups))
		for k := range app.groups {
			keys = append(keys, k)
		}
		sort.Strings(keys)

		writer := tabwriter.NewWriter(os.Stdout, 0, 8, 2, '\t', 0)
		for _, k := range keys {
			fmt.Fprintf(writer, "  %s\t%s\n", app.groups[k].Name, app.groups[k].Desc)
		}
		writer.Flush()
	}

	// commands
	if len(app.cmds) > 0 {
		fmt.Print("\nCommands:\n")

		keys := make([]string, 0, len(app.cmds))
		for k := range app.cmds {
			keys = append(keys, k)
		}
		sort.Strings(keys)

		writer := tabwriter.NewWriter(os.Stdout, 0, 8, 2, '\t', 0)
		for _, k := range keys {
			fmt.Fprintf(writer, "  %s\t%s\n", app.cmds[k].Name, app.cmds[k].Desc)
		}
		writer.Flush()
	}

	fmt.Printf("\nRun '%s COMMAND --help' for more information on a command.\n", os.Args[0])
}

func isFlag(s string) bool {
	if len(s) < 2 || s[0] != '-' {
		return false
	}
	return true
}

func nbFlags(fs *flag.FlagSet) int {
	count := 0
	fs.VisitAll(func(f *flag.Flag) {
		count++
	})
	return count
}
