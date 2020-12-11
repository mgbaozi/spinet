package main

import (
	"github.com/mgbaozi/spinet/pkg/apis"
	"github.com/urfave/cli/v2"
	"k8s.io/klog/v2"
)

func standAloneTask(c *cli.Context) error {
	file := c.String("from-file")
	taskSpec, err := apis.FromYamlFile(file)
	if err != nil {
		klog.Errorf("Parse yaml file failed with error: %v", err)
		return err
	}
	task, err := taskSpec.Parse()
	if err != nil {
		klog.Errorf("Parse task failed with error: %v", err)
		return err
	}
	if !dryRun {
		go serveHTTP(port)
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
