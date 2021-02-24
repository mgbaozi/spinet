package models

import (
	"fmt"
	"github.com/mgbaozi/spinet/pkg/common/utils"
	"github.com/mgbaozi/spinet/pkg/operators"
	"k8s.io/klog/v2"
)

func (step *Step) processConditions(ctx Context) (bool, error) {
	return ProcessConditions(ctx, operators.New("and"), step.Conditions)
}

func (step *Step) Process(ctx Context) (res bool, err error) {
	ctx.Trace(true, "processing "+ctx.shader.Key(), nil)
	defer func() {
		if err != nil {
			klog.V(3).Infof("Execute app %s failed: %v", step.App.AppName(), err)
		}
		ctx.Trace(err == nil, "process "+ctx.shader.Key()+" finished", res)
	}()
	if res, err := processSteps(ctx, step.Dependencies); err != nil || !res {
		return res, err
	}
	app := step.App
	klog.V(3).Infof("Running app: %s", app.AppName())
	var data interface{}
	if err = app.Execute(ctx, &data); err != nil {
		return false, err
	}
	ctx.SetAppData(data)
	ctx.Mapper(step.Mapper)
	return step.processConditions(ctx)
}

func processSteps(ctx Context, steps []Step) (res bool, err error) {
	ctx.Trace(true, "processing "+ctx.shader.Key(), nil)
	defer func() {
		ctx.Trace(err == nil, "process "+ctx.shader.Key()+" finished", res)
	}()
	var dependencyResults []interface{}
	for index := range steps {
		key := fmt.Sprintf("step-%d", index)
		magic := map[string]interface{}{
			"__index__": index,
		}
		if res, err := steps[index].Process(ctx.Sub(key, magic)); err != nil {
			return res, err
		} else {
			dependencyResults = append(dependencyResults, res)
		}
	}
	var val interface{}
	val, err = operators.New("and").Do(dependencyResults)
	return utils.ToBoolean(val), err
}
