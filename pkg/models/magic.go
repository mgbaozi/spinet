package models

import "time"

/*
Magic variables:
Global:
$DATE: current date
$TIME: current timestamp
$RUNTIME: cluster / standalone
*/

func init() {
	RegisterMagicVariable(MagicVariableNone{})
	RegisterMagicVariable(MagicVariableDate{})
	RegisterMagicVariable(MagicVariableTime{})
}

type MagicVariable interface {
	New(value interface{}) MagicVariable
	Name() string
	Data() interface{}
}

type MagicVariableNone struct{}

func (MagicVariableNone) Name() string {
	return "none"
}

func (MagicVariableNone) New(interface{}) MagicVariable {
	return MagicVariableNone{}
}

func (MagicVariableNone) Data() interface{} {
	return nil
}

type MagicVariableDate struct {
	date interface{}
}

func (MagicVariableDate) New(value interface{}) MagicVariable {
	date := time.Now().Format("2006-01-02")
	return MagicVariableDate{
		date: date,
	}
}

func (MagicVariableDate) Name() string {
	return "date"
}

func (m MagicVariableDate) Data() interface{} {
	return m.date
}

type MagicVariableTime struct {
	time interface{}
}

func (MagicVariableTime) New(value interface{}) MagicVariable {
	time := time.Now().Unix()
	return MagicVariableTime{
		time: time,
	}
}

func (MagicVariableTime) Name() string {
	return "time"
}

func (m MagicVariableTime) Data() interface{} {
	return m.time
}
