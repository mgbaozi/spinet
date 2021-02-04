package values

type Constant struct {
	value interface{}
}

func (*Constant) New(value map[string]interface{}) Value {
	return &Constant{
		value: value["value"],
	}
}

func (*Constant) Parse(str string) Value {
	return &Constant{
		value: str,
	}
}

func (*Constant) Type() ValueType {
	return ValueTypeConstant
}

func (constant *Constant) Format() interface{} {
	return ""
}

func (constant *Constant) Extract(data map[string]interface{}) (interface{}, error) {
	return constant.value, nil
}
