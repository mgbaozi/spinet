package main

import (
	"errors"
	"fmt"
	"github.com/mgbaozi/spinet/pkg/apis"
	"github.com/mgbaozi/spinet/pkg/client"
	"github.com/urfave/cli/v2"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"strings"
)

func applyTask(task apis.Task) error {
	c := client.NewClient(server).Tasks(task.Namespace)
	_, err := c.Create(&task)
	return err
}

func applyApp(app apis.CustomApp) error {
	c := client.NewClient(server).Apps()
	_, err := c.Create(&app)
	return err
}

func apply(c *cli.Context) error {
	file := c.String("from-file")
	content, err := ioutil.ReadFile(file)
	if err != nil {
		return err
	}
	var meta apis.Meta
	if err = yaml.Unmarshal(content, &meta); err != nil {
		return err
	}
	if meta.Kind == "" {
		meta.Kind = "task"
	}
	switch strings.ToLower(meta.Kind) {
	case "task":
		if task, err := apis.TaskFromYaml(content); err != nil {
			return err
		} else {
			return applyTask(task)
		}
	case "app":
		if app, err := apis.CustomAppFromYaml(content); err != nil {
			return err
		} else {
			return applyApp(app)
		}
	default:
		return errors.New(fmt.Sprintf("unsupported resource type %s", meta.Kind))
	}
	//TODO: return error as default
	return nil
}

var applyCli = &cli.Command{
	Name:  "apply",
	Usage: "apply resources to cluster",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:    "from-file",
			Aliases: []string{"f"},
			Value:   "",
			Usage:   "From a yaml file",
		},
	},
	Action: apply,
}
