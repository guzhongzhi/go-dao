package console

import (
	"fmt"
	"github.com/guzhongzhi/gmicro/config"
	"github.com/guzhongzhi/gmicro/logger"
	"github.com/urfave/cli/v2"
	"os"
	"path"
	"reflect"
	"strings"
)

type Config struct {
	Name      string
	Version   string
	BasePath  string
	EnvPrefix string
	out       interface{} //application config struct
}

//name: application name
//version: application version
//basePath: the base path of the application
//out: the config struct of the application
func NewConfig(name, version string, basePath string, out interface{}) *Config {
	if name == "" {
		name = "test"
	}
	if version == "" {
		version = "1.0"
	}
	if basePath == "" {
		basePath = path.Dir(path.Dir(os.Args[0]))
	}

	return &Config{
		Name:      name,
		Version:   version,
		BasePath:  basePath,
		EnvPrefix: "",
		out:       out,
	}
}

type Console interface {
	App() *cli.App
}

type console struct {
	config *Config
	app    *cli.App
}

func (s *console) App() *cli.App {
	return s.app
}

func New(cfg *Config) Console {
	app := cli.NewApp()
	app.Name = cfg.Name
	app.Version = cfg.Version
	app.Flags = []cli.Flag{
		&cli.StringFlag{
			Name:    "env",
			EnvVars: []string{cfg.EnvPrefix + "ENV", cfg.EnvPrefix + "env"},
			Value:   "dev",
			Usage:   "specify runtime environment: dev, qa, prod",
		},
		&cli.StringFlag{
			Name:    "config",
			EnvVars: []string{cfg.EnvPrefix + "CONFIG", cfg.EnvPrefix + "config"},
			Value:   cfg.BasePath + "/configs/",
			Usage:   "config file directory",
		},
	}
	app.Commands = []*cli.Command{
		&cli.Command{
			Name:  "env:show",
			Usage: "display all the config variables which can be changed by ENV",
			Action: func(ctx *cli.Context) error {
				env := ctx.String("env")
				cfgPath := ctx.String("config")
				err := config.LoadConfigFiles(cfgPath, env, cfg.out, logger.Default(), cfg.EnvPrefix)
				if err != nil {
					panic(err)
				}

				t := reflect.TypeOf(cfg.out)
				keys := config.GenerateCfgKeys(t, "")

				for key, _ := range keys {
					key = strings.ToUpper(key)
					key = strings.Replace(key, "/", ".", -1)
					fmt.Println(fmt.Sprintf("%s%s", cfg.EnvPrefix, key))
				}

				return nil
			},
		},
	}
	return &console{
		config: cfg,
		app:    app,
	}
}
