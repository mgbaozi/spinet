package models

type CustomApp struct {
	Task
	Modes   []AppMode
	Options map[string]interface{}
}

func (custom *CustomApp) New(options map[string]interface{}) App {
	return &CustomApp{
		Task:    custom.Task,
		Modes:   custom.Modes,
		Options: options,
	}
}

func (custom *CustomApp) AppName() string {
	return custom.Name
}

func (custom *CustomApp) Register() {
	RegisterApp(custom)
}

func (custom *CustomApp) AppModes() []AppMode {
	return custom.Modes
}

func (custom *CustomApp) Execute(mode AppMode, ctx *Context, data interface{}) (err error) {
	return
}
