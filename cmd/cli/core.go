package main

import (
	"github.com/mgbaozi/spinet/pkg/apps"
	"github.com/mgbaozi/spinet/pkg/models"
	"github.com/mgbaozi/spinet/pkg/operators"
	"github.com/mgbaozi/spinet/pkg/triggers"
	"github.com/urfave/cli/v2"
)

func core(c *cli.Context) error {
	task := &models.Task{
		Triggers: []models.Trigger{
			triggers.NewTimer(map[string]interface{}{
				"period": 10,
			}),
		},
		Inputs: []models.Step{
			{
				App: &apps.Simple{},
				Conditions: []models.Condition{
					{
						Operator: operators.EQ{},
						Values: []models.Value{
							{Type: "variable", Value: "content"},
							{Type: "constant", Value: "apple"},
						},
					},
				},
			},
		},
		Outputs: []models.Step{
			{
				App: &apps.Simple{
					Content: "ok",
				},
			},
		},
		Context: models.NewContext(),
	}
	task.Start()
	return nil
}
