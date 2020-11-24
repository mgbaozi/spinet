package models

type Hook struct {
	Name string
}

func (hook *Hook) TriggerType() string {
	return TriggerTypePassive
}
