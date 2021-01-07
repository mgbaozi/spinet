package models

import (
	"fmt"
	"strings"
)

type TraceItem interface {
	String() string
}

// TODO: task trace
type Result struct {
	Success bool
	Message string
	Data    interface{}
}

func (result Result) String() string {
	res := "Failed"
	if result.Success {
		res = "Success"
	}
	return fmt.Sprintf("[%s]%s(%v)", res, result.Message, result.Data)
}

type Trace struct {
	State string
	Items []TraceItem
	super *Trace
}

func NewTrace(state string) *Trace {
	return &Trace{
		State: state,
	}
}

func (trace *Trace) Enter(state string) *Trace {
	res := NewTrace(state)
	res.super = trace
	trace.Items = append(trace.Items, res)
	return res
}

func (trace *Trace) Leave() *Trace {
	return trace.super
}

func (trace *Trace) Push(success bool, message string, data interface{}) {
	result := Result{
		Success: success,
		Message: message,
		Data:    data,
	}
	trace.Items = append(trace.Items, result)
}

func (trace *Trace) String() string {
	var results []string
	for _, item := range trace.Items {
		results = append(results, item.String())
	}
	return strings.Join(results, "\n")
}
