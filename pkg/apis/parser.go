package apis

import (
	_ "github.com/mgbaozi/spinet/pkg/apps"
	"github.com/mgbaozi/spinet/pkg/models"
	_ "github.com/mgbaozi/spinet/pkg/operators"
	_ "github.com/mgbaozi/spinet/pkg/triggers"
	"gopkg.in/yaml.v3"
	"io/ioutil"
)

func FromYaml(content []byte) (Task, error) {
	var task Task
	err := yaml.Unmarshal(content, &task)
	return task, err
}

func FromYamlFile(filename string) (Task, error) {
	content, err := ioutil.ReadFile(filename)
	if err != nil {
		return Task{}, err
	}
	return FromYaml(content)
}

func (task Task) Parse() models.Task {
	context := models.NewContextWithDictionary(task.Dictionary)
	var triggers []models.Trigger
	var inputs []models.Input
	var conditions []models.Condition
	var outputs []models.Output
	for _, trigger := range task.Triggers {
		triggers = append(triggers, trigger.Parse())
	}
	for _, input := range task.Inputs {
		inputs = append(inputs, input.Parse())
	}
	for _, condition := range task.Conditions {
		conditions = append(conditions, condition.Parse())
	}
	for _, output := range task.Outputs {
		outputs = append(outputs, output.Parse())
	}
	return models.Task{
		Name:       task.Name,
		Triggers:   triggers,
		Inputs:     inputs,
		Conditions: conditions,
		Outputs:    outputs,
		Context:    context,
	}
}

func (trigger Trigger) Parse() models.Trigger {
	name := trigger.Type
	options := trigger.Options
	return models.NewTrigger(name, options)
}

func (input Input) Parse() models.Input {
	name := input.App
	options := input.Options
	var conditions []models.Condition
	for _, condition := range input.Conditions {
		conditions = append(conditions, condition.Parse())
	}
	return models.Input{
		App: models.NewApp(name, options),
		// TODO: Dictionary
		Conditions: conditions,
	}
}

func (condition Condition) Parse() models.Condition {
	name := condition.Operator
	var conditions []models.Condition
	for _, item := range condition.Conditions {
		conditions = append(conditions, item.Parse())
	}
	var values []models.Value
	for _, value := range condition.Values {
		values = append(values, models.ParseValue(value))
	}
	return models.Condition{
		Operator:   models.NewOperator(name),
		Conditions: conditions,
		Values:     values,
	}
}

func (output Output) Parse() models.Output {
	name := output.App
	options := output.Options
	var conditions []models.Condition
	for _, condition := range output.Conditions {
		conditions = append(conditions, condition.Parse())
	}
	return models.Output{
		App: models.NewApp(name, options),
		// TODO: Dictionary
		Conditions: conditions,
	}
}
