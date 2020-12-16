package models

import "k8s.io/klog/v2"

func (step *Step) processConditions(ctx *Context, appdata interface{}) (bool, error) {
	return ProcessConditions(NewOperator("and"), step.Conditions, ctx.Dictionary, appdata)
}

func (step *Step) Process(ctx *Context) (bool, error) {
	if res, err := processInputs(ctx, step.Dependencies); err != nil || !res {
		return res, err
	}
	app := step.App
	klog.V(3).Infof("Running app: %s", app.AppName())
	var data interface{}
	err := app.Execute(AppModeInput, ctx, &data)
	if err != nil {
		klog.V(3).Infof("Execute app failed: %v", err)
		return false, err
	}
	ctx.AppData = append(ctx.AppData, data)
	ProcessMapper(ctx, step.Mapper, data)
	res, err := step.processConditions(ctx, data)
	if err != nil {
		klog.V(3).Infof("Process conditions of app %s failed: %v", app.AppName(), err)
	}
	return res, err
}

func processInputs(ctx *Context, inputs []Step) (bool, error) {
	var dependencyResults []interface{}
	for index := range inputs {
		if res, err := inputs[index].Process(ctx); err != nil {
			return res, err
		} else {
			dependencyResults = append(dependencyResults, res)
		}
	}
	return NewOperator("and").Do(dependencyResults)
}