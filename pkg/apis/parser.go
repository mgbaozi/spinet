package apis

import (
	"github.com/mgbaozi/spinet/pkg/models"
	_ "github.com/mgbaozi/spinet/pkg/operators"
	_ "github.com/mgbaozi/spinet/pkg/triggers"
	"gopkg.in/yaml.v3"
	"io/ioutil"
)

func TaskFromYaml(content []byte) (Task, error) {
	var task Task
	err := yaml.Unmarshal(content, &task)
	return task, err
}

func TaskFromYamlFile(filename string) (Task, error) {
	content, err := ioutil.ReadFile(filename)
	if err != nil {
		return Task{}, err
	}
	return TaskFromYaml(content)
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
	res.SetDictionary(task.Dictionary)
	for _, trigger := range task.Triggers {
		res.Triggers = append(res.Triggers, trigger.Parse())
	}
	for _, input := range task.Inputs {
		if item, err := input.Parse(models.AppModeInput); err != nil {
			return res, err
		} else {
			res.Inputs = append(res.Inputs, item)
		}
	}
	for _, condition := range task.Conditions {
		res.Conditions = append(res.Conditions, condition.Parse())
	}
	for _, output := range task.Outputs {
		if item, err := output.Parse(models.AppModeOutPut); err != nil {
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

func (step Step) Parse(mode models.AppMode) (res models.Step, err error) {
	app := step.App
	options := step.Options
	for _, condition := range step.Conditions {
		res.Conditions = append(res.Conditions, condition.Parse())
	}
	for _, item := range step.Dependencies {
		if dependency, err := item.Parse(mode); err != nil {
			return res, err
		} else {
			res.Dependencies = append(res.Dependencies, dependency)
		}
	}
	res.Mapper = make(map[string]models.Value)
	for key, value := range step.Mapper {
		res.Mapper[key] = models.ParseValue(value)
	}
	res.Mode = mode
	res.App, err = models.NewApp(app, mode, options)
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

func CustomAppFromYaml(content []byte) (CustomApp, error) {
	var app CustomApp
	err := yaml.Unmarshal(content, &app)
	return app, err
}

func CustomAppFromYamlFile(filename string) (CustomApp, error) {
	content, err := ioutil.ReadFile(filename)
	if err != nil {
		return CustomApp{}, err
	}
	return CustomAppFromYaml(content)
}

func (app CustomApp) Parse() (res models.CustomApp, err error) {
	res.Modes = app.Modes
	res.Task, err = app.Task.Parse()
	return
}
