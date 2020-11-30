package models

type Operator interface {
	Do(ctx *Context, values ...Value) bool
}

type EQ struct{}

func (EQ) Do(ctx *Context, values ...Value) bool {
	if len(values) < 2 {
		return true
	}
	for i := 0; i < len(values)-1; i++ {
		if !values[i].Equals(values[i+1]) {
			return false
		}
	}
	return true
}
