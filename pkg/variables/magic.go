package variables

import (
	"github.com/mgbaozi/spinet/pkg/models"
	"time"
)

/*
Magic variables:
Global:
$DATE: current date
$TIME: current timestamp
$RUNTIME: cluster / standalone
*/

func init() {
	models.RegisterMagicVariable(MagicVariableNone{})
	models.RegisterMagicVariable(MagicVariableDate{})
	models.RegisterMagicVariable(MagicVariableTime{})
}

type MagicVariableNone struct{}

func (MagicVariableNone) Name() string {
	return "none"
}

func (MagicVariableNone) New(interface{}) models.MagicVariable {
	return MagicVariableNone{}
}

func (MagicVariableNone) Data() interface{} {
	return nil
}

type MagicVariableDate struct {
	date interface{}
}

func (MagicVariableDate) New(value interface{}) models.MagicVariable {
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

func (MagicVariableTime) New(value interface{}) models.MagicVariable {
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
