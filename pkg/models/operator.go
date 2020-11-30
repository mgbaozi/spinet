package models

import (
	"errors"
)

type Operator interface {
	Do(values []interface{}) (bool, error)
}

type EQ struct{}

func (EQ) Do(values []interface{}) (bool, error) {
	if len(values) < 2 {
		return true, nil
	}
	for i := 0; i < len(values)-1; i++ {
		if values[i] != values[i+1] {
			return false, nil
		}
	}
	return true, nil
}

type And struct{}

func (And) Do(values []interface{}) (bool, error) {
	for _, value := range values {
		if res, ok := value.(bool); ok {
			if !res {
				return false, nil
			}
		} else {
			return false, errors.New("operator 'And' execute failed: can't convert value to boolean")
		}
	}
	return true, nil
}

type Or struct{}

func (Or) Do(values []interface{}) (bool, error) {
	for _, value := range values {
		if res, ok := value.(bool); ok {
			if res {
				return true, nil
			}
		} else {
			return false, errors.New("operator 'Or' execute failed: can't convert value to boolean")
		}
	}
	return false, nil
}
