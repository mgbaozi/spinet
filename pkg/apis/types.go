package apis

import "github.com/mgbaozi/spinet/pkg/models"

type Meta struct {
	Name      string `json:"name" yaml:"name"`
	Namespace string `json:"namespace" yaml:"namespace"`
}

type Task struct {
	Meta       `json:",inline" yaml:",inline"`
	Dictionary map[string]interface{} `json:"dictionary" yaml:"dictionary"`
	Triggers   []Trigger              `json:"triggers" yaml:"triggers"`
	Conditions []Condition            `json:"conditions" yaml:"conditions"`
	Inputs     []Step                 `json:"inputs" yaml:"inputs"`
	Outputs    []Step                 `json:"outputs" yaml:"outputs"`
}

type Trigger struct {
	Type    string                 `json:"type" yaml:"type"`
	Options map[string]interface{} `json:"options" yaml:"options"`
}

type Condition struct {
	Operator   string        `json:"operator" yaml:"operator"`
	Conditions []Condition   `json:"conditions" yaml:"conditions"`
	Values     []interface{} `json:"values" yaml:"values"`
}

type Step struct {
	App          string                 `json:"app" yaml:"app"`
	Options      map[string]interface{} `json:"options" yaml:"options"`
	Mapper       map[string]interface{} `json:"mapper" yaml:"mapper"`
	Conditions   []Condition            `json:"conditions" yaml:"conditions"`
	Dependencies []Step                 `json:"dependencies" yaml:"dependencies"`
}

type Output struct {
	App        string                 `json:"app" yaml:"app"`
	Options    map[string]interface{} `json:"options" yaml:"options"`
	Mapper     map[string]interface{} `json:"mapper" yaml:"mapper"`
	Conditions []Condition            `json:"conditions" yaml:"conditions"`
	Outputs    []Output               `json:"outputs" yaml:"outputs"`
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
