package console

import "github.com/urfave/cli/v2"

type Console interface {
}

type console struct {
}

func NewApp(name string, version string, basePath string) *cli.App {
	app := cli.NewApp()
	app.Name = name
	app.Version = version
	app.Flags = []cli.Flag{
		&cli.StringFlag{
			Name:    "env",
			EnvVars: []string{"env"},
			Value:   "dev",
			Usage:   "specify runtime environment: dev, qa, prod",
		},
		&cli.StringFlag{
			Name:    "config",
			EnvVars: []string{"config"},
			Value:   basePath + "/configs/",
			Usage:   "config file directory",
		},
	}
	return app
}
