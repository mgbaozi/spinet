package main

import (
	"fmt"
	"github.com/mgbaozi/spinet/pkg/apis"
	"github.com/urfave/cli"
)

func standAloneTask(c *cli.Context) {
	file := c.String("from-file")
	dry := c.Bool("dry-run")
	taskSpec, err := apis.FromYamlFile(file)
	if err != nil {
		fmt.Println("Parse yaml file failed:", err)
	}
	fmt.Println(taskSpec)
	task := taskSpec.Parse()
	fmt.Println(task)
	if !dry {
		task.Start()
	}
}

var taskCli = cli.Command{
	Name:  "task",
	Usage: "Running task as stand alone mode",
	Flags: []cli.Flag{
		cli.StringFlag{
			Name:  "from-file, f",
			Value: "",
			Usage: "From a yaml file",
		},
		cli.BoolFlag{
			Name:   "dry-run, d",
			Hidden: false,
			Usage:  "Print task information only",
		},
	},
	Action: standAloneTask,
}
