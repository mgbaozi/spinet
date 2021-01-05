package main

import (
	"errors"
	"fmt"
	"github.com/mgbaozi/spinet/pkg/apis"
	"github.com/mgbaozi/spinet/pkg/client"
	"github.com/urfave/cli/v2"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"k8s.io/klog/v2"
	"strings"
)

func applyTask(task apis.Task) error {
	c := client.NewClient(server).Tasks(task.Namespace)
	if _, err := c.Create(&task); err != nil {
		klog.Errorf("Apply task failed with error: %v", err)
		return err
	}
	fmt.Printf("task %s applied\n", task.Name)
	return nil
}

func applyApp(app apis.CustomApp) error {
	c := client.NewClient(server).Apps()
	if _, err := c.Create(&app); err != nil {
		klog.Errorf("Apply app failed with error: %v", err)
		return err
	}
	fmt.Printf("app %s applied\n", app.Name)
	return nil
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
