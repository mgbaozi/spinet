package main

import (
	"os"

	"github.com/urfave/cli/v2"
)

var debug bool
var dryRun bool

func globalConfig(c *cli.Context) error {
	dryRun = c.Bool("dry-run")
	debug = c.Bool("debug")
	return klogInit(c)
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
			Name:    "verbose",
			Aliases: []string{"V"},
		},
		&cli.BoolFlag{
			Name:  "dry-run",
			Value: false,
		},
	}
	app.Flags = append(app.Flags, klogCliFlags...)
	app.Before = globalConfig
	app.Action = core
	app.Version = "0.0.1"
	app.Run(os.Args)
}
