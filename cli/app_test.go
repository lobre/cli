package cli_test

import (
	"flag"
	"fmt"
	"testing"

	"github.com/lobre/kits/cli"
)

func TestRun(t *testing.T) {
	app := cli.New()

	noteGroup := cli.Group{
		Name: "notes",
	}

	echoCmd := cli.Cmd{
		Name: "echo",
		Run:  echo,
	}

	noteGroup.AddCmd(&echoCmd)

	lsCmd := cli.Cmd{
		Name: "ls",
		Run:  ls,
	}

	app.AddGroup(&noteGroup)
	app.RootGroup().AddCmd(&lsCmd)
	app.Run()
}

func echo(fs *flag.FlagSet) error {
	fmt.Println("echo specific file")
	return nil
}

func ls(fs *flag.FlagSet) error {
	fmt.Println("ls all files")
	return nil
}
