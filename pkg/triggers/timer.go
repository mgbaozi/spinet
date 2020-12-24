package triggers

import (
	"github.com/mgbaozi/spinet/pkg/models"
	"k8s.io/klog/v2"
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
	ctx     models.Context
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

func (*Timer) TriggerName() string {
	return "timer"
}

func (timer *Timer) run() {
	ch := time.Tick(time.Duration(timer.Period) * time.Second)
	klog.V(2).Infof("Start timer with period: %d", timer.Period)
	go func() {
		for {
			_ = <-ch
			klog.V(4).Info("Tick!")
			timer.ch <- struct{}{}
		}
	}()
}

func (timer *Timer) Triggered(ctx models.Context) <-chan struct{} {
	timer.ctx = ctx
	if !timer.running {
		timer.running = true
		timer.run()
	}
	return timer.ch
}
