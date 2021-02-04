package models

import (
	"github.com/mgbaozi/spinet/pkg/common/utils"
	"github.com/mgbaozi/spinet/pkg/values"
)

func ProcessConditions(ctx Context, operator Operator, conditions []values.Value) (bool, error) {
	var values []interface{}
	for _, condition := range conditions {
		res, err := condition.Extract(ctx.MergedData())
		if err != nil {
			return false, err
		}
		values = append(values, res)
	}
	val, err := operator.Do(values)
	return utils.ToBoolean(val), err
}
