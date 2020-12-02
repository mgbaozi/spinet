package main

import (
	"fmt"
	"github.com/mgbaozi/spinet/pkg/apis"
	"github.com/urfave/cli/v2"
)

func standAloneTask(c *cli.Context) error {
	file := c.String("from-file")
	taskSpec, err := apis.FromYamlFile(file)
	if err != nil {
		fmt.Println("Parse yaml file failed:", err)
		return err
	}
	task := taskSpec.Parse()
	if !dryRun {
		task.Start()
	}
	return nil
}

var taskCli = &cli.Command{
	Name:  "task",
	Usage: "Running task as stand alone mode",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:    "from-file",
			Aliases: []string{"f"},
			Value:   "",
			Usage:   "From a yaml file",
		},
	},
	Action: standAloneTask,
}
