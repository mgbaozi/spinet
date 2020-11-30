package models

type Condition struct {
	Operator   Operator
	Conditions []Condition
	Values     []Value
}

func (condition Condition) Exec(data interface{}) (bool, error) {
	if len(condition.Conditions) > 0 {
		return ProcessConditions(condition.Operator, condition.Conditions, data)
	}
	var values []interface{}
	for _, value := range condition.Values {
		extracted, err := value.Extract(data)
		if err != nil {
			return false, err
		}
		values = append(values, extracted)
	}
	return condition.Operator.Do(values)
}

func ProcessConditions(operator Operator, conditions []Condition, data interface{}) (bool, error) {
	var values []interface{}
	for _, condition := range conditions {
		res, err := condition.Exec(data)
		if err != nil {
			return false, err
		}
		values = append(values, res)
	}
	return operator.Do(values)
}

func ProcessCommonConditions(conditions []Condition, ctx *Context) (bool, error) {
	return ProcessConditions(And{}, conditions, ctx.Dictionary)
}

func ProcessAppConditions(app string, conditions []Condition, ctx *Context) (bool, error) {
	data := ctx.Data[app]
	return ProcessConditions(And{}, conditions, data)
}
