package values

type Expression struct {
	value interface{}
}

func (*Expression) New(value map[string]interface{}) Value {
	return &Expression{
		value: value["value"],
	}
}

func (*Expression) Parse(str string) Value {
	return &Expression{
		value: str,
	}
}

func (*Expression) Type() ValueType {
	return ValueTypeExpression
}

func (variable *Expression) Format() string {
	return ""
}

func (variable *Expression) Extract(data map[string]interface{}) (interface{}, error) {
	return nil, nil
}
