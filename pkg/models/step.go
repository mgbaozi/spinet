package models

import (
	"fmt"
	"github.com/mgbaozi/spinet/pkg/common/utils"
	"github.com/mgbaozi/spinet/pkg/operators"
	"k8s.io/klog/v2"
	"sync"
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
	if res, err := ProcessSteps(ctx, step.Dependencies, nil); err != nil || !res {
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

func ProcessSteps(ctx Context, steps []Step, vars MagicVariables) (res bool, err error) {
	ctx.Trace(true, "processing "+ctx.shader.Key(), nil)
	defer func() {
		ctx.Trace(err == nil, "process "+ctx.shader.Key()+" finished", res)
	}()
	var dependencyResults []interface{}
	for index := range steps {
		key := fmt.Sprintf("step-%d", index)
		magic := map[string]interface{}{
			"__step__": index,
		}
		if vars != nil {
			for m, v := range vars {
				magic[m] = v
			}
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

func ConcurrencyProcessSteps(ctx Context, steps []Step, vars MagicVariables) (res bool, err error) {
	ctx.Trace(true, "processing "+ctx.shader.Key(), nil)
	defer func() {
		ctx.Trace(err == nil, "process "+ctx.shader.Key()+" finished", res)
	}()
	var dependencyResults = make([]interface{}, len(steps))
	var wg sync.WaitGroup
	for index := range steps {
		go func(index int) {
			defer wg.Done()
			key := fmt.Sprintf("step-%d", index)
			magic := map[string]interface{}{
				"__step__": index,
			}
			if vars != nil {
				for m, v := range vars {
					magic[m] = v
				}
			}
			wg.Add(1)
			if res, err := steps[index].Process(ctx.Sub(key, magic)); err != nil {
				dependencyResults[index] = err
			} else {
				dependencyResults[index] = res
			}
		}(index)
	}
	wg.Wait()
	//FIXME: handle error
	var val interface{}
	val, err = operators.New("and").Do(dependencyResults)
	return utils.ToBoolean(val), err
}
