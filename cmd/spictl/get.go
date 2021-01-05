package main

import (
	"fmt"
	"github.com/mgbaozi/spinet/pkg/client"
	"github.com/urfave/cli/v2"
	"k8s.io/klog/v2"
)

func getNamespaces(c *cli.Context) error {
	client := client.NewClient(server).Namespaces()
	namespaces, err := client.List()
	if err != nil {
		klog.Errorf("Get namespaces failed with error: %v", err)
		return err
	}
	fmt.Println("Name")
	for _, ns := range namespaces {
		fmt.Printf("%s\n", ns.Name)
	}
	return nil
}

func getTasks(c *cli.Context) error {
	client := client.NewClient(server).Tasks(namespace)
	tasks, err := client.List()
	if err != nil {
		klog.Errorf("Get tasks failed with error: %v", err)
		return err
	}
	fmt.Println("Name")
	for _, task := range tasks {
		fmt.Printf("%s\n", task.Name)
	}
	return nil
}

func getApps(c *cli.Context) error {
	client := client.NewClient(server).Apps()
	apps, err := client.List()
	if err != nil {
		klog.Errorf("Get apps failed with error: %v", err)
		return err
	}
	fmt.Println("Name")
	for _, app := range apps {
		fmt.Printf("%s\n", app.Name)
	}
	return nil
}

var getCli = &cli.Command{
	Name:  "get",
	Usage: "get resource from cluster",
	Subcommands: []*cli.Command{
		{
			Name:    "namespaces",
			Aliases: []string{"ns", "namespace"},
			Action:  getNamespaces,
		},
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
