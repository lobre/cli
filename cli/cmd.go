package cli

import (
	"flag"
	"fmt"
	"os"
	"text/tabwriter"
)

// Cmd represents an actual command.
type Cmd struct {
	// Name is the name of the command.
	Name string
	// Desc is the description that will be
	// used in the usage.
	Desc string
	// Run is the callback function that will be executed
	// if this command is the one desired by the end user.
	Run func(app, group, cmd *flag.FlagSet) error

	group *Group
	fs    *flag.FlagSet
}

// Flags will return a FlagSet that can be used to add
// command level flags.
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
