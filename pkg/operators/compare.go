package operators

type CompareResult int

const (
	CompareResultGreater CompareResult = -1
	CompareResultEqual                 = 0
	CompareResultLess                  = 1
)

func compare(lhs, rhs interface{}) CompareResult {
	// lvalue := reflect.ValueOf(lhs)
	// rvalue := reflect.ValueOf(rhs)
	return CompareResultEqual
}
