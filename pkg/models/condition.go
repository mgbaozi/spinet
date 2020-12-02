package models

type Condition struct {
	Operator   Operator
	Conditions []Condition
	Values     []Value
}

func (condition Condition) Exec(dictionary, appdata interface{}) (bool, error) {
	if len(condition.Conditions) > 0 {
		return ProcessConditions(condition.Operator, condition.Conditions, dictionary, appdata)
	}
	var values []interface{}
	for _, value := range condition.Values {
		extracted, err := value.Extract(dictionary, appdata)
		if err != nil {
			return false, err
		}
		values = append(values, extracted)
	}
	return condition.Operator.Do(values)
}

func ProcessConditions(operator Operator, conditions []Condition, dictionary, appdata interface{}) (bool, error) {
	var values []interface{}
	for _, condition := range conditions {
		res, err := condition.Exec(dictionary, appdata)
		if err != nil {
			return false, err
		}
		values = append(values, res)
	}
	return operator.Do(values)
}
