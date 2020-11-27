package models

type Timer struct {
	Period int
	ch <-chan struct{}
}

func NewTimer(options map[string]interface{}) Timer {
	return Timer{
		ch: make(<-chan struct{}),
	}
}

func (timer Timer) Triggered() <-chan struct{} {
	ch := make(<-chan struct{})
	return ch
}