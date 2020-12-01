package triggers

import (
	"fmt"
	"github.com/mgbaozi/spinet/pkg/models"
	"time"
)

func init() {
	models.RegisterTrigger(&Timer{})
}

type TimerOptions struct {
	Period int
}

func NewTimerOptions(options map[string]interface{}) TimerOptions {
	return TimerOptions{
		Period: options["period"].(int),
	}
}

type Timer struct {
	Options TimerOptions
	ch      chan struct{}
	running bool
}

func NewTimer(options map[string]interface{}) *Timer {
	return &Timer{
		Options: NewTimerOptions(options),
		ch:      make(chan struct{}),
	}
}

func (*Timer) New(options map[string]interface{}) models.Trigger {
	return NewTimer(options)
}

func (*Timer) Name() string {
	return "timer"
}

func (timer *Timer) run() {
	ch := time.Tick(time.Duration(timer.Options.Period) * time.Second)
	fmt.Println("Start timer...")
	for {
		_ = <-ch
		fmt.Println("Tick!")
		timer.ch <- struct{}{}
	}
}

func (timer *Timer) Triggered() <-chan struct{} {
	if !timer.running {
		timer.running = true
		go timer.run()
	}
	return timer.ch
}
