package models

type Source struct {
}

func (source *Source) TriggerType() string {
	return TriggerTypeActive
}