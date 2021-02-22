package values

import (
	"fmt"
	"github.com/mgbaozi/spinet/pkg/operators"
)

type Expression struct {
	operator operators.Operator
	values   []Value
}

func (*Expression) New(value map[string]interface{}) Value {
	//TODO: type checking
	name := value["operator"].(string)
	operator := operators.New(name)
	expression := &Expression{
		operator: operator,
	}
	values := value["values"].([]interface{})
	for _, item := range values {
		expression.values = append(expression.values, Parse(item))
	}
	return expression
}

func (*Expression) Parse(str string) Value {
	panic("unimplemented error")
	return &Expression{}
}

func (*Expression) Type() ValueType {
	return ValueTypeExpression
}

func (variable *Expression) Format() interface{} {
	return ""
}

func (variable *Expression) Extract(data map[string]interface{}) (interface{}, error) {
	var values []interface{}
	for _, item := range variable.values {
		if value, err := item.Extract(data); err != nil {
			return nil, err
		} else {
			values = append(values, value)
		}
	}
	return variable.operator.Do(values)
}

func (variable *Expression) String() string {
	return fmt.Sprintf("Expression [op=%s,values=%v]",
		variable.operator.Name(), variable.values)
}
