package main

import (
	"github.com/mgbaozi/spinet/pkg/models"
	"github.com/urfave/cli"
)

func core(c *cli.Context) {
	task := &models.Task{
		Triggers: []models.Trigger{
			models.NewTimer(map[string]interface{}{
				"period": 10,
			}),
		},
		Inputs: []models.Input{
			{
				App: &models.Simple{},
				Conditions: []models.Condition{
					{
						Operator: models.EQ{},
						Values: []models.Value{
							{Type: "variable", Value: "content"},
							{Type: "constant", Value: "apple"},
						},
					},
				},
			},
		},
		Outputs: []models.Output{
			{
				App: &models.Simple{
					Content: "ok",
				},
			},
		},
		Context: models.NewContext(),
	}
	task.Start()
}
