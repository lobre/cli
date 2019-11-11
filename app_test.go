package cli_test

import (
	"flag"
	"os"
	"strings"
	"testing"

	"github.com/lobre/cli"
)

func setArgs(args string) {
	program := os.Args[0]
	os.Args = []string{
		program,
	}

	os.Args = append(os.Args, strings.Split(args, " ")...)
}

func TestApplicationFlag(t *testing.T) {
	setArgs("-debug")

	app := cli.New()
	app.Flags().Bool("debug", false, "debug mode")
	app.Default = func(app *flag.FlagSet) error {
		if debug := app.Lookup("debug"); debug == nil {
			t.Fatal("application flag lost")
		}
		return nil
	}
	app.Run()
}

func TestApplicationArg(t *testing.T) {
	setArgs("myargument")

	app := cli.New()
	app.Default = func(app *flag.FlagSet) error {
		if app.Arg(0) != "myargument" {
			t.Fatal("application arg lost")
		}
		return nil
	}
	app.Run()
}

func TestApplicationCmd(t *testing.T) {
	setArgs("mycmd")

	jobDone := false

	app := cli.New()
	cmd := cli.Cmd{
		Name: "mycmd",
		Run: func(app, group, cmd *flag.FlagSet) error {
			jobDone = true
			return nil
		},
	}
	app.AddCmd(&cmd)
	app.Run()

	if !jobDone {
		t.Fatal("application command not called")
	}
}

func TestFull(t *testing.T) {
	setArgs("-p myproject container -clean recreate -f nginx")

	jobDone := false

	app := cli.New()
	app.Flags().String("p", "", "project")

	group := cli.Group{
		Name: "container",
	}
	group.Flags().Bool("clean", false, "cleaning")

	cmd := cli.Cmd{
		Name: "recreate",
		Run: func(app, group, cmd *flag.FlagSet) error {
			// check app flag
			project := app.Lookup("p")
			if project == nil {
				t.Fatal("project application flag lost")
			}

			if project.Value.String() != "myproject" {
				t.Fatal("project application value wrong")
			}

			if clean := group.Lookup("clean"); clean == nil {
				t.Fatal("group flag clean lost")
			}

			if force := cmd.Lookup("f"); force == nil {
				t.Fatal("command flag f lost")
			}

			if cmd.Arg(0) != "nginx" {
				t.Fatal("command arg nginx is wrong")
			}

			jobDone = true
			return nil
		},
	}
	cmd.Flags().Bool("f", false, "force")

	group.AddCmd(&cmd)
	app.AddGroup(&group)
	app.Run()

	if !jobDone {
		t.Fatal("command did not execute")
	}
}
