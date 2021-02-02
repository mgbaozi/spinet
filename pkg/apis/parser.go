package apis

import (
	"github.com/mgbaozi/spinet/pkg/models"
	_ "github.com/mgbaozi/spinet/pkg/operators"
	_ "github.com/mgbaozi/spinet/pkg/triggers"
	"github.com/mgbaozi/spinet/pkg/values"
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
		task.Namespace = models.DefaultNamespace
	}
	return task
}

func (task Task) Parse() (res models.Task, err error) {
	task.Validate()
	res.Meta = models.Meta{
		Name:      task.Name,
		Namespace: task.Namespace,
	}
	res.Context = models.NewContext()
	res.Aggregator = make(models.Mapper)
	res.SetDictionary(task.Dictionary)
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
	res.Aggregator = models.ParseMapper(task.Aggregator)
	return res, nil
}

func (trigger Trigger) Parse() models.Trigger {
	name := trigger.Type
	options := trigger.Options
	return models.NewTrigger(name, options)
}

func (step Step) Parse() (res models.Step, err error) {
	app := step.App
	options := step.Options
	for _, condition := range step.Conditions {
		res.Conditions = append(res.Conditions, condition.Parse())
	}
	for _, item := range step.Dependencies {
		if dependency, err := item.Parse(); err != nil {
			return res, err
		} else {
			res.Dependencies = append(res.Dependencies, dependency)
		}
	}
	res.Mapper = make(map[string]values.Value)
	for key, value := range step.Mapper {
		res.Mapper[key] = values.Parse(value)
	}
	res.App, err = models.NewApp(app, options)
	return res, err
}

func (condition Condition) Parse() models.Condition {
	name := condition.Operator
	var conditions []models.Condition
	for _, item := range condition.Conditions {
		conditions = append(conditions, item.Parse())
	}
	var vals []values.Value
	for _, value := range condition.Values {
		vals = append(vals, values.Parse(value))
	}
	return models.Condition{
		Operator:   models.NewOperator(name),
		Conditions: conditions,
		Values:     vals,
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
	res.Task, err = app.Task.Parse()
	for _, item := range app.Options {
		option, _ := item.Parse()
		res.DefinedOptions = append(res.DefinedOptions, option)
	}
	return
}

func (option AppOptionItem) Parse() (res models.AppOptionItem, err error) {
	return models.AppOptionItem{
		Name:     option.Name,
		Type:     option.Type,
		Required: option.Required,
	}, nil
}
