package triggers

import (
	"github.com/mgbaozi/spinet/pkg/logging"
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
	TimerOptions
	ch      chan struct{}
	running bool
}

func NewTimer(options map[string]interface{}) *Timer {
	return &Timer{
		TimerOptions: NewTimerOptions(options),
		ch:           make(chan struct{}),
	}
}

func (*Timer) New(options map[string]interface{}) models.Trigger {
	return NewTimer(options)
}

func (*Timer) Name() string {
	return "timer"
}

func (timer *Timer) run() {
	ch := time.Tick(time.Duration(timer.Period) * time.Second)
	logging.Info("Start timer with period: %d", timer.Period)
	for {
		_ = <-ch
		logging.Debug("Tick!")
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
