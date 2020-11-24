package main

import (
	"github.com/mgbaozi/spinet/pkg/models"
	"github.com/urfave/cli"
)

func core(c *cli.Context) {
	task := &models.Task{
		Inputs: []models.Input{
			models.SimpleInput{},
		},
		Outputs: []models.Output{
			models.SimpleOutput{
				Content: "$.content",
			},
		},
		Context: models.NewContext(),
	}
	task.Execute()
}