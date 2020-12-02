package apis

type Meta struct {
	Name string `json:"name" yaml:"name"`
}

type Task struct {
	Meta
	Dictionary map[string]interface{} `json:"dictionary" yaml:"dictionary"`
	Triggers   []Trigger              `json:"triggers" yaml:"triggers"`
	Conditions []Condition            `json:"conditions" yaml:"conditions"`
	Inputs     []Input                `json:"inputs" yaml:"inputs"`
	Outputs    []Output               `json:"outputs" yaml:"outputs"`
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

type Input struct {
	App        string                 `json:"app" yaml:"app"`
	Options    map[string]interface{} `json:"options" yaml:"options"`
	Mapper     map[string]interface{} `json:"mapper" yaml:"mapper"`
	Conditions []Condition            `json:"conditions" yaml:"conditions"`
	Inputs     []Input                `json:"inputs" yaml:"inputs"`
}

type Output struct {
	App        string                 `json:"app" yaml:"app"`
	Options    map[string]interface{} `json:"options" yaml:"options"`
	Conditions []Condition            `json:"conditions" yaml:"conditions"`
	Outputs    []Output               `json:"outputs" yaml:"outputs"`
}
