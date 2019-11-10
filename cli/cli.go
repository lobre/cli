package cli

import (
	"flag"
	"fmt"
	"os"
	"sort"
	"text/tabwriter"
)

var errorHandling = flag.ExitOnError

type Cmd struct {
	Name string
	Desc string
	Run  func(app, group, cmd *flag.FlagSet) error

	group *Group
	fs    *flag.FlagSet
}

func (cmd *Cmd) Flags() *flag.FlagSet {
	if cmd.fs == nil {
		cmd.fs = flag.NewFlagSet(cmd.Name, errorHandling)
		cmd.fs.Usage = cmd.defaultUsage
	}
	return cmd.fs
}

func (cmd *Cmd) defaultUsage() {
	var opt string
	if cmd.fs.NFlag() > 0 {
		opt = " [OPTIONS]"
	}
	if cmd.group != nil {
		fmt.Printf("\nUsage:	%s %s %s%s [PARAMS...]\n", os.Args[0], cmd.group.Name, cmd.Name, opt)
	} else {
		fmt.Printf("\nUsage:	%s %s%s [PARAMS...]\n", os.Args[0], cmd.Name, opt)
	}

	if cmd.Desc != "" {
		fmt.Printf("\n%s\n", cmd.Desc)
	}

	// options
	if cmd.fs.NFlag() > 0 {
		fmt.Print("\nOptions:\n")
		writer := tabwriter.NewWriter(os.Stdout, 0, 8, 2, '\t', 0)
		cmd.fs.VisitAll(func(f *flag.Flag) {
			fmt.Fprintf(writer, "  -%s\t%s\n", f.Name, f.Usage)
		})
		writer.Flush()
	}
}

type Group struct {
	Name    string
	Desc    string
	Default func(app, group *flag.FlagSet) error

	cmds map[string]*Cmd
	fs   *flag.FlagSet
}

func (group *Group) Flags() *flag.FlagSet {
	if group.fs == nil {
		group.fs = flag.NewFlagSet(group.Name, errorHandling)
		group.fs.Usage = group.defaultUsage
	}
	return group.fs
}

func (group *Group) defaultUsage() {
	var opt string
	if group.fs.NFlag() > 0 {
		opt = " [OPTIONS]"
	}
	fmt.Printf("\nUsage:	%s %s%s COMMAND\n", os.Args[0], group.Name, opt)

	if group.Desc != "" {
		fmt.Printf("\n%s\n", group.Desc)
	}

	// options
	if group.fs.NFlag() > 0 {
		fmt.Print("\nOptions:\n")
		writer := tabwriter.NewWriter(os.Stdout, 0, 8, 2, '\t', 0)
		group.fs.VisitAll(func(f *flag.Flag) {
			fmt.Fprintf(writer, "  -%s\t%s\n", f.Name, f.Usage)
		})
		writer.Flush()
	}

	// commands
	if len(group.cmds) > 0 {
		fmt.Print("\nCommands:\n")

		keys := make([]string, 0, len(group.cmds))
		for k := range group.cmds {
			keys = append(keys, k)
		}
		sort.Strings(keys)

		writer := tabwriter.NewWriter(os.Stdout, 0, 8, 2, '\t', 0)
		for _, k := range keys {
			fmt.Fprintf(writer, "  %s\t%s\n", group.cmds[k].Name, group.cmds[k].Desc)
		}
		writer.Flush()
	}
}

// AddCmd will add or replace a cli command.
// If the command does not have any name, it will replace the default command.
// If the command is nil, it will panic
func (group *Group) AddCmd(cmd *Cmd) {
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

	if group.cmds == nil {
		group.cmds = make(map[string]*Cmd)
	}

	cmd.group = group
	group.cmds[cmd.Name] = cmd
}

type App struct {
	Desc string

	groups map[string]*Group
	cmds   map[string]*Cmd

	fs *flag.FlagSet

	Default func(app *flag.FlagSet) error
}

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

func (app *App) Flags() *flag.FlagSet {
	if app.fs == nil {
		app.fs = flag.NewFlagSet(os.Args[0], errorHandling)
		app.fs.Usage = app.defaultUsage
	}
	return app.fs
}

// AddGroup will add or replace an cli group.
// If the group does not have any name, it will replace the default group.
// If the group is nil, it will panic
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

func (app *App) Run() {
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
			fmt.Println(err)
			os.Exit(1)
		}

	// with root command
	case groupIdx == 0 && cmdIdx != 0:
		cmd := app.cmds[os.Args[cmdIdx]]

		app.fs.Parse(os.Args[1:cmdIdx])
		if len(os.Args) > cmdIdx+1 {
			cmd.fs.Parse(os.Args[cmdIdx+1:])
		}

		if err := cmd.Run(app.fs, nil, cmd.fs); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

	// with group default action
	case groupIdx != 0 && cmdIdx == 0:
		group := app.groups[os.Args[groupIdx]]

		app.fs.Parse(os.Args[1:groupIdx])
		if len(os.Args) > groupIdx+1 {
			group.fs.Parse(os.Args[groupIdx+1:])
		}

		if err := group.Default(app.fs, group.fs); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

	// with app default action
	default:
		if len(os.Args) > 1 {
			app.fs.Parse(os.Args[1:])
		}

		if err := app.Default(app.fs); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	}

	os.Exit(0)
}

func (app *App) defaultUsage() {
	var opt string
	if app.fs.NFlag() > 0 {
		opt = " [OPTIONS]"
	}
	fmt.Printf("\nUsage:	%s%s COMMAND\n", os.Args[0], opt)

	if app.Desc != "" {
		fmt.Printf("\n%s\n", app.Desc)
	}

	// options
	if app.fs.NFlag() > 0 {
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
}

func isFlag(s string) bool {
	if len(s) < 2 || s[0] != '-' {
		return false
	}
	return true
}
