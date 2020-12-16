package apps

import (
	"github.com/mgbaozi/spinet/pkg/models"
)

type Custom struct {
	models.Task
	Modes   []models.AppMode
	Options map[string]interface{}
}

func (custom *Custom) New(options map[string]interface{}) models.App {
	return &Custom{
		Task:    custom.Task,
		Modes:   custom.Modes,
		Options: options,
	}
}

func (custom *Custom) AppName() string {
	return custom.Name
}

func (custom *Custom) Register() {
	models.RegisterApp(custom)
}

func (custom *Custom) AppModes() []models.AppMode {
	return custom.Modes
}

func (custom *Custom) Execute(mode models.AppMode, ctx *models.Context, data interface{}) (err error) {
	return
}
