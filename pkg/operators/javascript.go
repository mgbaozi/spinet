package operators

import (
	"fmt"
	"github.com/robertkrimen/otto"
)

const tmpl = `
(function(values){
%s
})(values)
`

type JavaScript struct{}

func (JavaScript) Name() string {
	return "javascript"
}

func (op JavaScript) Do(values []interface{}) (res interface{}, err error) {
	defer func() {
		logOperatorResult(op.Name(), res, err)
	}()
	vm := otto.New()
	vm.Set("values", values)
	return vm.Run(fmt.Sprintf(tmpl, "return 2"))
}
