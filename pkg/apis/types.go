package apis

import "github.com/mgbaozi/spinet/pkg/models"

type TypeMeta struct {
	Kind string `json:"kind" yaml:"kind"`
}

type Meta struct {
	TypeMeta  `json:",inline" yaml:",inline"`
	Name      string `json:"name" yaml:"name"`
	Namespace string `json:"namespace" yaml:"namespace"`
}

type Namespace struct {
	TypeMeta `json:",inline" yaml:",inline"`
	Name     string `json:"name" yaml:"name"`
}

type Task struct {
	Meta       `json:",inline" yaml:",inline"`
	Dictionary map[string]interface{} `json:"dictionary" yaml:"dictionary"`
	Triggers   []Trigger              `json:"triggers" yaml:"triggers"`
	Conditions []Condition            `json:"conditions,omitempty" yaml:"conditions"`
	Inputs     []Step                 `json:"inputs,omitempty" yaml:"inputs"`
	Outputs    []Step                 `json:"outputs,omitempty" yaml:"outputs"`
}

type Trigger struct {
	Type    string                 `json:"type" yaml:"type"`
	Options map[string]interface{} `json:"options" yaml:"options"`
}

type Condition struct {
	Operator   string        `json:"operator" yaml:"operator"`
	Conditions []Condition   `json:"conditions,omitempty" yaml:"conditions"`
	Values     []interface{} `json:"values" yaml:"values"`
}

type Step struct {
	App          string                 `json:"app" yaml:"app"`
	Options      map[string]interface{} `json:"options" yaml:"options"`
	Mapper       map[string]interface{} `json:"mapper,omitempty" yaml:"mapper"`
	Conditions   []Condition            `json:"conditions,omitempty" yaml:"conditions"`
	Dependencies []Step                 `json:"dependencies,omitempty" yaml:"dependencies"`
}

type CustomApp struct {
	Task  `json:",inline" yaml:",inline"`
	Modes []models.AppMode `json:"modes" yaml:"modes"`
}

type App struct {
	Name    string                 `json:"name" yaml:"name"`
	Options map[string]interface{} `json:"options" yaml:"options"`
	Modes   []models.AppMode       `json:"modes" yaml:"modes"`
}

type Resource interface {
	UrlFormat() string
	Type() string
}

func (task *Task) Type() string {
	return task.Kind
}

func (*Task) URL() string {
	return "/namespaces/%s/tasks"
}

func (app *CustomApp) Type() string {
	return app.Kind
}

func (*CustomApp) URL() string {
	return "/apps"
}
