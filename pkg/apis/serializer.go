package apis

import "github.com/mgbaozi/spinet/pkg/models"

func SerializeTask(task models.Task) (res Task) {
	res.Meta = Meta{
		Name:      task.Name,
		Namespace: task.Namespace,
	}
	res.Dictionary = task.OriginDictionary
	for _, trigger := range task.Triggers {
		res.Triggers = append(res.Triggers, SerializeTrigger(trigger))
	}
	for _, input := range task.Inputs {
		res.Inputs = append(res.Inputs, SerializeStep(input))
	}
	for _, condition := range task.Conditions {
		res.Conditions = append(res.Conditions, SerializeCondition(condition))
	}
	for _, output := range task.Outputs {
		res.Outputs = append(res.Outputs, SerializeStep(output))
	}
	return
}

func SerializeTrigger(trigger models.Trigger) (res Trigger) {
	res.Type = trigger.TriggerName()
	//TODO: do not use Options function
	res.Options = trigger.Options()
	return
}

func SerializeStep(step models.Step) (res Step) {
	app := step.App
	res.App = app.AppName()
	//TODO: do not use Options function
	res.Options = app.Options()
	for _, condition := range step.Conditions {
		res.Conditions = append(res.Conditions, SerializeCondition(condition))
	}
	//TODO use origin style value
	res.Mapper = make(map[string]interface{})
	for key, value := range step.Mapper {
		res.Mapper[key] = value.Format()
	}
	for _, dep := range step.Dependencies {
		res.Dependencies = append(res.Dependencies, SerializeStep(dep))
	}
	return
}

func SerializeCondition(condition models.Condition) (res Condition) {
	res.Operator = condition.Operator.Name()
	for _, condition := range condition.Conditions {
		res.Conditions = append(res.Conditions, SerializeCondition(condition))
	}
	for _, value := range condition.Values {
		res.Values = append(res.Values, value.Format())
	}
	return
}
