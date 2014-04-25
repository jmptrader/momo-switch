package location

import (
	"encoding/json"
	"entry"
	"fmt"
	"net/http"
	"strconv"
	"zkmanager"
)

const (
	SWITCH_OFF_GOREDIS     = "/switch/location_notify/switchOff_GoRedis"
	SWITCH_ON_GOREDIS_READ = "/switch/location_notify/switchOn_GoRedis_Read"
	SWITCH_ON_RADAR        = "/switch/location_notify/switchOn_friend_radar"
	SWITCH_ON_RADAR_LOG    = "/switch/location_notify/switchOn_friend_radar_record"
)

type OpTag struct {
	Label  string `json:"label"`
	Status bool   `json:"status"`
}

type RadaGoRedis struct {
	zkmanager *zkmanager.ZKManager
}

func NewRadaGoRedis(zkmanager *zkmanager.ZKManager) *RadaGoRedis {

	//创建节点
	zkmanager.TraverseCreatePath(SWITCH_OFF_GOREDIS, "true")
	zkmanager.TraverseCreatePath(SWITCH_ON_GOREDIS_READ, "false")
	zkmanager.TraverseCreatePath(SWITCH_ON_RADAR, "true")
	zkmanager.TraverseCreatePath(SWITCH_ON_RADAR_LOG, "true")

	for _, v := range []string{SWITCH_OFF_GOREDIS, SWITCH_ON_GOREDIS_READ, SWITCH_ON_RADAR, SWITCH_ON_RADAR_LOG} {
		data := zkmanager.Get(v)
		fmt.Printf("获取[%s]数据成功!data:[%b]\n", v, data)
	}

	return &RadaGoRedis{zkmanager: zkmanager}
}

func (self *RadaGoRedis) HandleLocationNotifySwitchQ(resp http.ResponseWriter, req *http.Request) {

	switchOff_GoRedis_bool := self.zkmanager.Get(SWITCH_OFF_GOREDIS)

	switchOn_GoRedis_Read_bool := self.zkmanager.Get(SWITCH_ON_GOREDIS_READ)

	SWITCH_ON_RADAR_bool := self.zkmanager.Get(SWITCH_ON_RADAR)

	//好友雷达日志开关
	switch_on_radar_log_bool := self.zkmanager.Get(SWITCH_ON_RADAR_LOG)

	tags := []OpTag{
		OpTag{Label: "switchOn_friend_radar", Status: SWITCH_ON_RADAR_bool},
		OpTag{Label: "switchOn_GoRedis", Status: !switchOff_GoRedis_bool},
		OpTag{Label: "switchOn_GoRedis_Read", Status: switchOn_GoRedis_Read_bool},
		OpTag{Label: "switchOn_radar_log", Status: switch_on_radar_log_bool},
	}

	status, _ := json.Marshal(tags)

	resp.Write(status)

}

func (self *RadaGoRedis) HandleLocationNotifySwitch(resp http.ResponseWriter, req *http.Request) {

	switchOn_GoRedis := req.FormValue("switchOn_GoRedis")
	switchOn_GoRedis_Read := req.FormValue("switchOn_GoRedis_Read")
	switchOn_friend_radar := req.FormValue("switchOn_friend_radar")
	switchOn_radar_log := req.FormValue("switchOn_radar_log")

	succ := false
	reponse := &entry.Response{}
	//关闭goredis
	if len(switchOn_GoRedis) > 0 {
		switchOn_GoRedis_bool, _ := strconv.ParseBool(switchOn_GoRedis)
		succ = self.zkmanager.SetGoRedisSwitch(SWITCH_OFF_GOREDIS, strconv.FormatBool(!switchOn_GoRedis_bool))
	}

	//是否打开goredis的读
	if len(switchOn_GoRedis_Read) > 0 {
		succ = self.zkmanager.SetGoRedisSwitch(SWITCH_ON_GOREDIS_READ, switchOn_GoRedis_Read)
	}

	//是否打开好友雷达
	if len(switchOn_friend_radar) > 0 {
		succ = self.zkmanager.SetGoRedisSwitch(SWITCH_ON_RADAR, switchOn_friend_radar)
	}

	if len(switchOn_radar_log) > 0 {
		succ = self.zkmanager.SetGoRedisSwitch(SWITCH_ON_RADAR_LOG, switchOn_radar_log)
	}

	if succ {
		reponse.Ec = 200
		reponse.Em = "设置成功！"
		majson, _ := json.Marshal(reponse)
		resp.Write(majson)
	} else {
		reponse.Ec = 505
		reponse.Em = "设置失败！"
		majson, _ := json.Marshal(reponse)
		resp.Write(majson)
	}
}
