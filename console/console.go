package console

import "github.com/urfave/cli/v2"

type Console interface {
}

type console struct {
}

func NewApp(name string, version string) *cli.App {
	app := cli.NewApp()
	app.Name = name
	app.Version = version
	return app
}
