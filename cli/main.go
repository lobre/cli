package cli

import (
	"errors"
	"flag"
	"os"
)

var ErrorHandling = flag.ExitOnError

type Action struct {
	Name string
	Run  func()

	flagset *flag.FlagSet
}

func (a *Action) FlagSet() *flag.FlagSet {
	if a.flagset == nil {
		a.flagset = flag.NewFlagSet(a.Name, ErrorHandling)
	}

	return a.flagset
}

type Object struct {
	Name string

	actions map[string]*Action
	flagset *flag.FlagSet
}

func (o *Object) AddAction(a ...*Action) {
	o.actions = append(o.actions, a...)
}

func (o *Object) FlagSet() *flag.FlagSet {
	if o.flagset == nil {
		o.flagset = flag.NewFlagSet(o.Name, ErrorHandling)
	}

	return o.flagset
}

type App struct {
	Name string

	root Object
	objs []*Object
}

func New() *App {
	var app App
	app.root = Object{}
	return &app
}

func (app *App) AddObject(obj ...*Object) {
	app.objs = append(app.objs, obj...)
}

func (app *App) AddAction(a ...*Action) {
	app.root.actions = append(app.root.actions, a...)
}

func (app *App) FlagSet() *flag.FlagSet {
	return app.root.FlagSet()
}

func (app *App) Parse() error {
	if len(os.Args) < 2 {
		return errors.New("expected at least action in command")
	}

	var actionIdx, objectIdx int
	for i := 1; i < len(os.Args); i++ {
		arg := os.Args[i]

		if !isFlag(arg) {
			if actionIdx == 0 {
				actionIdx = i
				break
			}
			if objectIdx == 0 {
				objectIdx = actionIdx
				actionIdx = i
				break
			}
		}
	}

	if actionIdx == 0 {
		return errors.New("expected at least action in command")
	}

    // parse object flags
    switch 

	if objectIdx == 0 {
	} else {
        switch os.Args[1] {
            case:
        }
    }

	return nil
}

func isFlag(s string) bool {
	if len(s) < 2 || s[0] != '-' {
		return false
	}
	return true
}
