package main

import (
	"github.com/mgbaozi/spinet/pkg/apis"
	_ "github.com/mgbaozi/spinet/pkg/apps"
	"github.com/mgbaozi/spinet/pkg/logging"
	_ "github.com/mgbaozi/spinet/pkg/variables"
	"k8s.io/klog/v2"
	"os"

	"github.com/urfave/cli/v2"
)

var debug bool
var dryRun bool
var port int

func registerCustomApps(c *cli.Context) error {
	files := c.StringSlice("custom-app")
	for _, file := range files {
		if appSpec, err := apis.CustomAppFromYamlFile(file); err != nil {
			return err
		} else {
			if app, err := appSpec.Parse(); err != nil {
				return err
			} else {
				app.Register()
			}
		}
	}
	return nil
}

func globalConfig(c *cli.Context) error {
	dryRun = c.Bool("dry-run")
	debug = c.Bool("debug")
	port = c.Int("port")

	if err := logging.Init(c, debug); err != nil {
		klog.V(2).Infof("Init klog failed with error: %v", err)
		return err
	}
	return registerCustomApps(c)
}

func main() {
	cli.VersionFlag = &cli.BoolFlag{
		Name:  "version",
		Value: false,
	}
	app := cli.NewApp()
	app.Name = "spinet-cli"
	app.Usage = "Spinet command line tools"
	app.Commands = []*cli.Command{
		taskCli,
	}
	app.Flags = []cli.Flag{
		&cli.BoolFlag{
			Name:    "debug",
			Aliases: []string{"D"},
			Value:   false,
		},
		&cli.BoolFlag{
			Name:  "dry-run",
			Value: false,
		},
		&cli.IntFlag{
			Name:    "port",
			Aliases: []string{"p"},
			Value:   8080,
			Usage:   "Port for http service",
		},
		&cli.StringSliceFlag{
			Name:    "custom-app",
			Aliases: []string{"a"},
			Usage:   "Custom app yaml file",
		},
	}
	app.Flags = append(app.Flags, logging.CliFlags...)
	app.Before = globalConfig
	app.Action = core
	app.Version = "0.0.1"
	app.Run(os.Args)
}
