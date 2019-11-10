# go-cli

Package cli provides a simple to way to create a cli application with commands, groups of commands (subcommands) and flags.

I know cobra (https://github.com/spf13/cobra) is the most used open source package for this purpose but it has a wider scope and I usually don't want to use a bazooka for creating really simple cli applications. That's why I created this package.

cli tries to follow the same philosophy introduced by the Docker cli. It comes with two kinds of objects:
 - cmd
 - group

To better explain them, see an example with Docker cli. Groups are what Docker calls management commands. They are a subcommands related to a specific object that needs to be manipulated. It means:
 - networks
 - volumes
 - containers
 - images

And then, commands can be applied to these management commands.
See this example.

    docker container create nginx -f

Here "container" is the group, "create" is the command, "nginx" is a parameter and "-f" is a flag.

There can be only one level of groups. This package does not allow creating groups into groups.

You can also directly define commands that are not attached to any groups. For example:

    docker ps

The main component in this cli package in an app. An app is the root object that will contain either direct commands (as "ps" above"), or groups of commands.

Here is how to create a simple command.

    app := cli.New()
    app.Desc = "This is the app description"
    app.AddCmd(&cli.Cmd{
        Name: "ps",
        Run: func(app, group, cmd *flag.FlagSet) {
            fmt.Println("here is the result of ps")
        },
    })
    app.Run()

A command can also be assigned to a group that itself will be added to the app.

    app := cli.New()
    app.Desc = "This is the app description"
    ls := cli.Cmd{
        Name: "ls",
        Run: func(app, group, cmd *flag.FlagSet) {
            fmt.Println("here is the result of container ls")
        }
    }
    containers := cli.Group{
        Name: "container",
        Desc: "Containers management",
    })
    containers.AddCmd(&ls)
    app.AddGroup(&containers)
    app.Run()

Finally, you can define flags onto the app, group or command. You just have to call the `Flags()` method on one of them and it will return a regular `flag.FlagSet` that you can use to add flags.
Then, you can use the `Lookup` function on a `FlagSet` from the `Run` function of your command to retrieve a flag or use `Args` to retrieve a non flag argument.

    app := cli.New()
    app.Desc = "This is the app description"
    app.Flags().Bool("debug", false, "debug mode")
    app.AddCmd(&cli.Cmd{
        Name: "ps",
        Run: func(app, group, cmd *flag.FlagSet) {
            if debug := app.Lookup("debug"); debug != nil {
                fmt.Println("set the debug mode")
            }
            fmt.Println("here is the result of ps")
        },
    })
    app.Run()

Where running your cli app, a usage documentation will be printed if you add the `--help` flag, or if you don't provide a command or a group.
