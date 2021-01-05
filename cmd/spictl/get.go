package main

import (
	"github.com/urfave/cli/v2"
)

func getTasks(c *cli.Context) error {
	return nil
}

func getApps(c *cli.Context) error {
	return nil
}

var getCli = &cli.Command{
	Name:  "get",
	Usage: "get resource from cluster",
	Subcommands: []*cli.Command{
		{
			Name:    "tasks",
			Aliases: []string{"t", "task"},
			Action:  getTasks,
		},
		{
			Name:    "apps",
			Aliases: []string{"a", "app"},
			Action:  getApps,
		},
	},
}
