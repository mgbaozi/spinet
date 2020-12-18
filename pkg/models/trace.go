package models

import (
	"fmt"
	"strings"
)

// TODO: task trace
type Result struct {
	Level   int
	Success bool
	Message string
	Data    interface{}
}

func (result Result) String() string {
	prefix := strings.Repeat("  ", result.Level)
	res := "Failed"
	if result.Success {
		res = "Success"
	}
	return fmt.Sprintf("%s[%s]%s(%v)", prefix, res, result.Message, result.Data)
}

type Trace struct {
	Level   int
	Results []Result
}

//FIXME: different level for each goroutine
func (trace *Trace) Indent() {
	trace.Level += 1
}

func (trace *Trace) UnIndent() {
	trace.Level -= 1
}

func (trace *Trace) Push(success bool, message string, data interface{}) {
	result := Result{
		Level:   trace.Level,
		Success: success,
		Message: message,
		Data:    data,
	}
	trace.Results = append(trace.Results, result)
}

func (trace *Trace) String() string {
	var results []string
	for _, item := range trace.Results {
		results = append(results, item.String())
	}
	return strings.Join(results, "\n")
}
