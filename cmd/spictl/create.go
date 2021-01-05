package main

import (
	"errors"
	"github.com/mgbaozi/spinet/pkg/apis"
	"github.com/mgbaozi/spinet/pkg/client"
	"github.com/urfave/cli/v2"
	"k8s.io/klog/v2"
)

func createNamespace(c *cli.Context) error {
	client := client.NewClient(server).Namespaces()
	name := c.Args().First()
	if name == "" {
		klog.Errorf("Namespace name required")
		return errors.New("need namespace name")
	}
	namespace := apis.Namespace{Name: name}
	if _, err := client.Create(&namespace); err != nil {
		klog.Errorf("Create namespace failed with error: %v", err)
		return err
	}
	return nil
}

var createCli = &cli.Command{
	Name:  "create",
	Usage: "create resources on cluster",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:    "from-file",
			Aliases: []string{"f"},
			Value:   "",
			Usage:   "From a yaml file",
		},
	},
	Subcommands: []*cli.Command{

		{
			Name:    "namespaces",
			Aliases: []string{"ns", "namespace"},
			Action:  createNamespace,
		},
		{
			Name:    "tasks",
			Aliases: []string{"t", "task"},
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:    "from-file",
					Aliases: []string{"f"},
					Value:   "",
					Usage:   "From a yaml file",
				},
			},
			Action: apply,
		},
		{
			Name:    "apps",
			Aliases: []string{"a", "app"},
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:    "from-file",
					Aliases: []string{"f"},
					Value:   "",
					Usage:   "From a yaml file",
				},
			},
			Action: apply,
		},
	},
}
