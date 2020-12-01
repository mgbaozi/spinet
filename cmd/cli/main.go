package main

import (
	"github.com/mgbaozi/spinet/pkg/logging"
	"os"

	"github.com/urfave/cli/v2"
)

var debug bool
var verbose bool
var dryRun bool

func globalConfig(c *cli.Context) error {
	logLevel := c.String("log-level")
	verbose = c.Bool("verbose")
	debug = c.Bool("debug")
	dryRun = c.Bool("dry-run")
	if verbose {
		logging.SetLevel(logging.TraceLevel)
	} else if debug {
		logging.SetLevel(logging.DebugLevel)
	} else {
		logging.SetLevelWithString(logLevel)
	}
	return nil
}

func main() {
	app := cli.NewApp()
	app.Name = "spinet-cli"
	app.Usage = "Spinet command line tools"
	app.Commands = []*cli.Command{
		taskCli,
	}
	app.Flags = []cli.Flag{
		&cli.StringFlag{
			Name:    "log-level",
			Aliases: []string{"log", "L"},
			Value:   "info",
		},
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
	app.Before = globalConfig
	app.Action = core
	app.Version = "0.0.1"
	app.Run(os.Args)
}
