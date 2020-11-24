package models

type Request struct {

}

type Result struct {
	Variables map[string]string
}

type ProcessorFunc func(data string, req *Request) (res *Result)

type ExecutorFunc func(data interface{}) error