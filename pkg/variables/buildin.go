package variables

import (
	"github.com/mgbaozi/spinet/pkg/models"
	"time"
)

/*
Build-in variables:
Global:
$DATE: current date
$TIME: current timestamp
$RUNTIME: cluster / standalone
*/

func init() {
	models.RegisterBuildInVariable(BuildInVariableNone{})
	models.RegisterBuildInVariable(BuildInVariableDate{})
	models.RegisterBuildInVariable(BuildInVariableTime{})
}

type BuildInVariableNone struct{}

func (BuildInVariableNone) Name() string {
	return "none"
}

func (BuildInVariableNone) New(interface{}) models.BuildInVariable {
	return BuildInVariableNone{}
}

func (BuildInVariableNone) Data() interface{} {
	return nil
}

type BuildInVariableDate struct {
	date interface{}
}

func (BuildInVariableDate) New(value interface{}) models.BuildInVariable {
	date := time.Now().Format("2006-01-02")
	return BuildInVariableDate{
		date: date,
	}
}

func (BuildInVariableDate) Name() string {
	return "date"
}

func (m BuildInVariableDate) Data() interface{} {
	return m.date
}

type BuildInVariableTime struct {
	time interface{}
}

func (BuildInVariableTime) New(value interface{}) models.BuildInVariable {
	time := time.Now().Unix()
	return BuildInVariableTime{
		time: time,
	}
}

func (BuildInVariableTime) Name() string {
	return "time"
}

func (m BuildInVariableTime) Data() interface{} {
	return m.time
}
