package operators

type EQ struct{}

func (EQ) Name() string {
	return "eq"
}

func (op EQ) Do(values []interface{}) (res bool, err error) {
	defer func() {
		logOperatorResult(op.Name(), res, err)
	}()
	if len(values) < 2 {
		return true, nil
	}
	for i := 0; i < len(values)-1; i++ {
		if !isEqual(values[i], values[i+1]) {
			return false, nil
		}
	}
	return true, nil
}

type Greater struct{}

func (Greater) Name() string {
	return "gt"
}

func (op Greater) Do(values []interface{}) (res bool, err error) {
	defer func() {
		logOperatorResult(op.Name(), res, err)
	}()
	if len(values) < 2 {
		return true, nil
	}
	for i := 0; i < len(values)-1; i++ {
		res, err := compare(values[i], values[i+1])
		if err != nil {
			return false, err
		}
		if res != CompareResultGreater {
			return false, nil
		}
	}
	return true, nil
}

type Less struct{}

func (Less) Name() string {
	return "lt"
}

func (op Less) Do(values []interface{}) (res bool, err error) {
	defer func() {
		logOperatorResult(op.Name(), res, err)
	}()
	if len(values) < 2 {
		return true, nil
	}
	for i := 0; i < len(values)-1; i++ {
		res, err := compare(values[i], values[i+1])
		if err != nil {
			return false, err
		}
		if res != CompareResultLess {
			return false, err
		}
	}
	return true, nil
}

type LE struct{}

func (LE) Name() string {
	return "le"
}

func (op LE) Do(values []interface{}) (res bool, err error) {
	defer func() {
		logOperatorResult(op.Name(), res, err)
	}()
	if len(values) < 2 {
		// TODO: need return error when not have enough values
		return true, nil
	}
	for i := 0; i < len(values)-1; i++ {
		if isEqual(values[i], values[i+1]) {
			continue
		}
		res, err := compare(values[i], values[i+1])
		if err != nil {
			return false, err
		}
		if res != CompareResultLess && res != CompareResultEqual {
			return false, err
		}
	}
	return true, nil
}

type GE struct{}

func (GE) Name() string {
	return "ge"
}

func (op GE) Do(values []interface{}) (res bool, err error) {
	defer func() {
		logOperatorResult(op.Name(), res, err)
	}()
	if len(values) < 2 {
		// TODO: need return error when not have enough values
		return true, nil
	}
	for i := 0; i < len(values)-1; i++ {
		if isEqual(values[i], values[i+1]) {
			continue
		}
		res, err := compare(values[i], values[i+1])
		if err != nil {
			return false, err
		}
		if res != CompareResultGreater && res != CompareResultEqual {
			return false, err
		}
	}
	return true, nil
}
