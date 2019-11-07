package cli

import (
	"flag"
	"fmt"
	"os"
)

var ErrorHandling = flag.ExitOnError

type Cmd struct {
	Name string
	Run  func(fs *flag.FlagSet) error
	fs   *flag.FlagSet
}

func (cmd *Cmd) Flags() *flag.FlagSet {
	if cmd.fs == nil {
		cmd.fs = flag.NewFlagSet(cmd.Name, ErrorHandling)
	}

	return cmd.fs
}

type Group struct {
	Name string
	cmds map[string]*Cmd
}

// AddCmd will add or replace a cli command.
// If the command does not have any name, it will replace the default command.
// If the command is nil, it will panic
func (group *Group) AddCmd(cmd *Cmd) {
	if cmd == nil {
		panic("cannot add nil command")
	}

	if cmd.fs == nil {
		cmd.fs = flag.NewFlagSet(cmd.Name, ErrorHandling)
	}

	if group.cmds == nil {
		group.cmds = make(map[string]*Cmd)
	}

	group.cmds[cmd.Name] = cmd
}

func (group *Group) RootCmd() *Cmd {
	return group.cmds[""]
}

type App struct {
	Name   string
	groups map[string]*Group
}

func New() *App {
	var app App

	app.groups = make(map[string]*Group)
	app.AddGroup(&Group{})

	return &app
}

// TODO: to implement
func (app *App) Usage() string {
	return "usage"
}

func (app *App) RootGroup() *Group {
	return app.groups[""]
}

// AddGroup will add or replace an cli group.
// If the group does not have any name, it will replace the default group.
// If the group is nil, it will panic
func (app *App) AddGroup(group *Group) {
	if group == nil {
		panic("cannot add nil group")
	}

	if group.cmds == nil {
		group.cmds = make(map[string]*Cmd)
	}

	if _, ok := group.cmds[""]; !ok {
		group.AddCmd(&Cmd{})
	}

	if group.cmds[""].Run == nil {
		group.cmds[""].Run = app.defaultRun
	}

	app.groups[group.Name] = group
}

func (app *App) defaultRun(fs *flag.FlagSet) error {
	fmt.Println(app.Usage())
	return nil
}

func (app *App) Run() {
	var groupArg, cmdArg string
	var flagsIdx int = 1

	switch {
	case len(os.Args) >= 3 && !isFlag(os.Args[1]) && !isFlag(os.Args[2]):
		groupArg = os.Args[1]
		cmdArg = os.Args[2]
		flagsIdx = 3
	case len(os.Args) >= 2 && !isFlag(os.Args[1]):
		cmdArg = os.Args[1]
		flagsIdx = 2
	}

	var group *Group
	var cmd *Cmd
	var ok bool

	if group, ok = app.groups[groupArg]; !ok {
		fmt.Printf("group %s does not exist\n", groupArg)
		fmt.Println(app.Usage())
		os.Exit(0)
	}

	if cmd, ok = group.cmds[cmdArg]; !ok {
		if groupArg == "" {
			fmt.Printf("command %s does not exist\n", cmdArg)
		} else {
			fmt.Printf("command %s does not exist in group %s\n", cmdArg, groupArg)
		}
		fmt.Println(app.Usage())
		os.Exit(0)
	}

	if len(os.Args) > flagsIdx {
		cmd.fs.Parse(os.Args[flagsIdx:])
	}

	if err := cmd.Run(cmd.fs); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	os.Exit(0)
}

func isFlag(s string) bool {
	if len(s) < 2 || s[0] != '-' {
		return false
	}
	return true
}
