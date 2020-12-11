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

func (task *Task) Validate() *Task {
	if len(task.Namespace) == 0 {
		task.Namespace = "default"
	}
	return task
}

func (task Task) Parse() (res models.Task, err error) {
	task.Validate()
	res.Meta = models.Meta{
		Name:      task.Name,
		Namespace: task.Namespace,
	}
	res.Context = models.NewContextWithDictionary(task.Dictionary)
	for _, trigger := range task.Triggers {
		res.Triggers = append(res.Triggers, trigger.Parse())
	}
	for _, input := range task.Inputs {
		if item, err := input.Parse(); err != nil {
			return res, err
		} else {
			res.Inputs = append(res.Inputs, item)
		}
	}
	for _, condition := range task.Conditions {
		res.Conditions = append(res.Conditions, condition.Parse())
	}
	for _, output := range task.Outputs {
		if item, err := output.Parse(); err != nil {
			return res, err
		} else {
			res.Outputs = append(res.Outputs, item)
		}
	}
	return res, nil
}

func (trigger Trigger) Parse() models.Trigger {
	name := trigger.Type
	options := trigger.Options
	return models.NewTrigger(name, options)
}

func (input Input) Parse() (res models.Input, err error) {
	app := input.App
	options := input.Options
	for _, condition := range input.Conditions {
		res.Conditions = append(res.Conditions, condition.Parse())
	}
	for _, item := range input.Dependencies {
		if dependency, err := item.Parse(); err != nil {
			return res, err
		} else {
			res.Dependencies = append(res.Dependencies, dependency)
		}
	}
	res.App, err = models.NewApp(app, models.AppModeInput, options)
	return res, err
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

func (output Output) Parse() (res models.Output, err error) {
	app := output.App
	options := output.Options
	for _, condition := range output.Conditions {
		res.Conditions = append(res.Conditions, condition.Parse())
	}
	res.App, err = models.NewApp(app, models.AppModeOutPut, options)
	return res, err
}
