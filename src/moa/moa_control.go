package moa

import (
	"encoding/json"
	"fmt"
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

	instance := req.FormValue("instance")

	//没有参数，那么就查询所有的服务
	if len(instance) <= 0 {
		instanceNames := self.moamanager.InstanceNames

		namesJson, _ := json.Marshal(instanceNames)
		fmt.Println("query:" + string(namesJson))
		resp.Write(namesJson)
	} else {

		v, ok := self.moamanager.Instances[instance]
		if ok {
			names, _ := json.Marshal(v)
			fmt.Println("query:" + string(names))
			resp.Write(names)
		}
	}

}
