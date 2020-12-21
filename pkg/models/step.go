package models

import (
	"fmt"
	"k8s.io/klog/v2"
)

func (step *Step) processConditions(ctx *Context, appdata interface{}) (bool, error) {
	return ProcessConditions(NewOperator("and"), step.Conditions, ctx.Dictionary, appdata)
}

func (step *Step) Process(ctx *Context, key string) (res bool, err error) {
	ctx.Trace.Indent()
	ctx.Trace.Push(true, "processing "+key, nil)
	defer func() {
		if err != nil {
			klog.V(3).Infof("Execute app %s failed: %v", step.App.AppName(), err)
		}
		ctx.Trace.Push(err == nil, "process "+key+" finished", res)
		ctx.Trace.UnIndent()
	}()
	if res, err := processSteps(ctx, step.Dependencies, key); err != nil || !res {
		return res, err
	}
	app := step.App
	klog.V(3).Infof("Running app: %s", app.AppName())
	var data interface{}
	if err = app.Execute(ctx, &data); err != nil {
		return false, err
	}
	ctx.AppData[key] = data
	ProcessMapper(ctx, step.Mapper, data)
	return step.processConditions(ctx, data)
}

func processSteps(ctx *Context, steps []Step, prefix string) (res bool, err error) {
	ctx.Trace.Indent()
	ctx.Trace.Push(true, "processing "+prefix, nil)
	defer func() {
		ctx.Trace.Push(err == nil, "process "+prefix+" finished", res)
		ctx.Trace.UnIndent()
	}()
	var dependencyResults []interface{}
	for index := range steps {
		key := fmt.Sprintf("%s.%d", prefix, index)
		if res, err := steps[index].Process(ctx, key); err != nil {
			return res, err
		} else {
			dependencyResults = append(dependencyResults, res)
		}
	}
	return NewOperator("and").Do(dependencyResults)
}
