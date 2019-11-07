package cli

import (
	"flag"
	"fmt"
	"os"
)

var ErrorHandling = flag.ExitOnError

type Action struct {
	Name string
	Run  func(fs *flag.FlagSet) error
	fs   *flag.FlagSet
}

func (act *Action) Flags() *flag.FlagSet {
	if act.fs == nil {
		act.fs = flag.NewFlagSet(act.Name, ErrorHandling)
	}
	return act.fs
}

type Object struct {
	Name    string
	actions map[string]*Action
}

// AddAction will add or replace an cli action.
// The action name is taken as the unique ID so it should be provided.
// If the action does not have any name, it will replace the default action.
// If the action is nil, it will panic
func (obj *Object) AddAction(act *Action) {
	if act == nil {
		panic("cannot add nil action")
	}

	obj.actions[act.Name] = act
}

func (obj *Object) RootAction() *Action {
	return obj.actions[""]
}

type App struct {
	Name    string
	objects map[string]*Object
}

func New() *App {
	var app App

	app.objects = make(map[string]*Object)
	app.AddObject(&Object{})

	return &app
}

// TODO: to implement
func (app *App) Usage() string {
	return "usage"
}

func (app *App) RootObject() *Object {
	return app.objects[""]
}

// AddObject will add or replace an cli object.
// If the object does not have any name, it will replace the default object.
// If the object is nil, it will panic
func (app *App) AddObject(obj *Object) {
	if obj == nil {
		panic("cannot add nil object")
	}

	if obj.actions == nil {
		obj.actions = make(map[string]*Action)
	}

	if _, ok := obj.actions[""]; !ok || obj.actions[""].Run == nil {
		obj.actions[""].Run = app.defaultRun
	}

	app.objects[obj.Name] = obj
}

func (app *App) defaultRun(fs *flag.FlagSet) error {
	fmt.Println(app.Usage())
	return nil
}

func (app *App) Run() {
	var obj, act string
	var flagsIdx int = 1

	switch {
	case len(os.Args) >= 3 && !isFlag(os.Args[1]) && !isFlag(os.Args[2]):
		obj = os.Args[2]
		act = os.Args[1]
		flagsIdx = 3
	case len(os.Args) >= 2 && !isFlag(os.Args[1]):
		act = os.Args[1]
		flagsIdx = 2
	}

	var object *Object
	var action *Action
	var ok bool

	if object, ok = app.objects[obj]; !ok {
		fmt.Println("object does not exist")
		fmt.Println(app.Usage())
		os.Exit(0)
	}

	if action, ok = object.actions[act]; !ok {
		fmt.Println("action does not exist")
		fmt.Println(app.Usage())
		os.Exit(0)
	}

	if len(os.Args) > flagsIdx {
		action.fs.Parse(os.Args[flagsIdx:])
	}

	if err := action.Run(action.fs); err != nil {
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
