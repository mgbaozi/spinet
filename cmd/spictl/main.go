package main

import (
	_ "github.com/mgbaozi/spinet/pkg/apps"
	"github.com/mgbaozi/spinet/pkg/logging"
	_ "github.com/mgbaozi/spinet/pkg/variables"
	"k8s.io/klog/v2"
	"os"

	"github.com/urfave/cli/v2"
)

var debug bool
var dryRun bool
var server string
var namespace string

func globalConfig(c *cli.Context) error {
	dryRun = c.Bool("dry-run")
	debug = c.Bool("debug")
	server = c.String("server")
	namespace = c.String("namespace")

	if err := logging.Init(c, debug); err != nil {
		klog.V(2).Infof("Init klog failed with error: %v", err)
		return err
	}
	return nil
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
		applyCli,
		getCli,
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
		&cli.StringFlag{
			Name:    "server",
			Aliases: []string{"s"},
			Value:   "http://localhost:8080",
			Usage:   "API server address",
		},
		&cli.StringFlag{
			Name:    "namespace",
			Aliases: []string{"n"},
			Value:   "default",
			Usage:   "Resource namespace",
		},
	}
	app.Flags = append(app.Flags, logging.CliFlags...)
	app.Before = globalConfig
	// app.Action = core
	app.Version = "0.0.1"
	app.Run(os.Args)
}
