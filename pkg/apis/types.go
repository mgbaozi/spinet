package apis

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
	Conditions []interface{}          `json:"conditions,omitempty" yaml:"conditions"`
	Inputs     []Step                 `json:"inputs,omitempty" yaml:"inputs"`
	Outputs    []Step                 `json:"outputs,omitempty" yaml:"outputs"`
	Aggregator map[string]interface{} `json:"aggregator,omitempty" yaml:"aggregator"`
}

type Trigger struct {
	Type    string                 `json:"type" yaml:"type"`
	Options map[string]interface{} `json:"options" yaml:"options"`
}

type Step struct {
	App          string                 `json:"app" yaml:"app"`
	Options      map[string]interface{} `json:"options" yaml:"options"`
	Mapper       map[string]interface{} `json:"mapper,omitempty" yaml:"mapper"`
	Conditions   []interface{}          `json:"conditions,omitempty" yaml:"conditions"`
	Dependencies []Step                 `json:"dependencies,omitempty" yaml:"dependencies"`
}

type AppOptionItem struct {
	Name     string `json:"name" yaml:"name"`
	Type     string `json:"type" yaml:"type"`
	Required bool   `json:"required" yaml:"required"`
}

type CustomApp struct {
	Task    `json:",inline" yaml:",inline"`
	Options []AppOptionItem `json:"options" yaml:"options"`
}

type App struct {
	Name    string          `json:"name" yaml:"name"`
	Options []AppOptionItem `json:"options" yaml:"options"`
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
