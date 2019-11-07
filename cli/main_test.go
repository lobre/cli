package cli_test

import (
	"fmt"
	"testing"

	"github.com/lobre/kits/cli"
)

func TestParse(t *testing.T) {
	app := cli.New()

	app.AddAction(&cli.Action{
		Name: "run",
		Run: func() {
			fmt.Println("hello")
		},
	})
}
