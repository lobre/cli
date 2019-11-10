package cli

import (
	"flag"
	"fmt"
	"os"
	"sort"
	"text/tabwriter"
)

// Group represents a group of commands.
type Group struct {
	// Name is the group name.
	Name string
	// Desc is the description that will be printed
	// in the usage.
	Desc string
	// Default is a function that will be called if
	// this group has been choosen without any commands.
	// It will by default show the usage.
	Default func(app, group *flag.FlagSet) error

	cmds map[string]*Cmd
	fs   *flag.FlagSet
}

// Flags will return a FlagSet that can be used to add
// group level flags.
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

	fmt.Printf("\nRun '%s %s COMMAND --help' for more information on a command.\n", os.Args[0], group.Name)
}

// AddCmd will add or replace a cli command.
// If the command is nil or does not have any name it will panic
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
