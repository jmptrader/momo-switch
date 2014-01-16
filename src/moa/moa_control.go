package moa

import (
	"encoding/json"
	"net/http"
)

type MoaControl struct {
	moamanager *MoaInStanceManager
}

func InitMoaControl() *MoaControl {
	manager := &MoaInStanceManager{}
	manager.ScheduleInitHosts()

	return &MoaControl{moamanager: manager}
}

func (self *MoaControl) HandleQueryMoaNameQ(resp http.ResponseWriter, req *http.Request) {
	instanceNames := self.moamanager.InstanceNames

	names, err := json.Marshal(instanceNames)
	if nil != err {

	}
	resp.Write(names)
}
