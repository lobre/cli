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

// SetAction will add or replace an cli action.
// The action name is taken as the unique ID so it should be provided.
func (o *Object) SetAction(a *Action) error {
	if a == nil {
		return errors.New("cannot add nil action")
	}
	if a.Name == "" {
		return errors.New("cannot add action with empty name")
	}
	o.actions[a.Name] = a
}

func (o *Object) FlagSet() *flag.FlagSet {
	if o.flagset == nil {
		o.flagset = flag.NewFlagSet(o.Name, ErrorHandling)
	}

	return o.flagset
}

type App struct {
	Name string

	root    Object
	objects map[string]*Object
}

func New() *App {
	var app App
	app.root = Object{}
	return &app
}

// SetObject will add or replace an cli object.
// The object name is taken as the unique ID so it should be provided.
func (app *App) SetObject(obj *Object) error {
	if obj == nil {
		return errors.New("cannot add nil object")
	}
	if obj.Name == "" {
		return errors.New("cannot add object with empty name")
	}
	app.objects[obj.Name] = obj
}

// SetAction will add or replace an cli action of the default object.
// The action name is taken as the unique ID so it should be provided.
func (app *App) SetAction(a *Action) error {
	return app.root.SetAction(a)
}

func (app *App) FlagSet() *flag.FlagSet {
	return app.root.FlagSet()
}

func (app *App) Parse() error {
	if len(os.Args) < 2 || isFlag(os.Args[1]) {
		return errors.New("expected at least action in command")
	}

	// TODO: use these flag
	var hasObj bool
	var actionIdx int = 1

	// Rewrite in order to set the action as optional (because of myapp --version for instance)

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

	if objectIdx == 0 {
		//
	} else if _, ok := app.objects[os.Args[1]]; ok {
		// this feels wrong because objectIdx is useless
	}

	return nil
}

func isFlag(s string) bool {
	if len(s) < 2 || s[0] != '-' {
		return false
	}
	return true
}
