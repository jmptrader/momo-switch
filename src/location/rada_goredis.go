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

	//关闭好友雷达
	SWITCH_OFF_GOREDIS = "/switch/location_notify/switchOff_GoRedis"
	//是否打开好友雷达的goredis读
	SWITCH_ON_GOREDIS_READ = "/switch/location_notify/switchOn_GoRedis_Read"
	//是否打开好友雷达
	SWITCH_ON_RADAR = "/switch/location_notify/switchOn_friend_radar"
	//好友雷达的LOG开关
	SWITCH_ON_RADAR_LOG = "/switch/location_notify/switchOn_friend_radar_record"
	//带关注和群组的留言板
	SWITCH_ON_FEED_V2 = "/switch/feed/v2_switch"
	//附近列表用户信用等级开关
	SWITCH_ON_GEO_UPDATE_CREDIT = "/switch/geo_update/credit_switch"
)

type OpTag struct {
	Label  string `json:"label"`
	Status bool   `json:"status"`
}

type RadaGoRedis struct {
	zkmanager *zkmanager.ZKManager
}

func NewRadaGoRedis(zkmanager *zkmanager.ZKManager) *RadaGoRedis {

	switches := [][]string{
		{SWITCH_OFF_GOREDIS, "true"},
		{SWITCH_ON_GOREDIS_READ, "false"},
		{SWITCH_ON_RADAR, "true"},
		{SWITCH_ON_RADAR_LOG, "true"},
		{SWITCH_ON_FEED_V2, "true"},
		{SWITCH_ON_GEO_UPDATE_CREDIT, "true"}}

	for _, v := range switches {
		//创建节点
		zkmanager.TraverseCreatePath(v[0], v[1])
		data := zkmanager.Get(v[0])
		fmt.Printf("获取[%s]数据成功!data:[%b]\n", v, data)
	}

	return &RadaGoRedis{zkmanager: zkmanager}
}

func (self *RadaGoRedis) HandleLocationNotifySwitchQ(resp http.ResponseWriter, req *http.Request) {

	switchOff_GoReDdis_bool := self.zkmanager.Get(SWITCH_OFF_GOREDIS)

	switchOn_GoRedis_Read_bool := self.zkmanager.Get(SWITCH_ON_GOREDIS_READ)

	SWITCH_ON_RADAR_bool := self.zkmanager.Get(SWITCH_ON_RADAR)

	//好友雷达日志开关
	switch_on_radar_log_bool := self.zkmanager.Get(SWITCH_ON_RADAR_LOG)

	switch_on_feed_v2 := self.zkmanager.Get(SWITCH_ON_FEED_V2)

	switch_on_geo_update_credit := self.zkmanager.Get(SWITCH_ON_GEO_UPDATE_CREDIT)

	tags := []OpTag{
		OpTag{Label: "switchOn_friend_radar", Status: SWITCH_ON_RADAR_bool},
		OpTag{Label: "switchOn_GoRedis", Status: !switchOff_GoRedis_bool},
		OpTag{Label: "switchOn_GoRedis_Read", Status: switchOn_GoRedis_Read_bool},
		OpTag{Label: "switchOn_radar_log", Status: switch_on_radar_log_bool},
		OpTag{Label:"switch_on_feed_v2",Status:switch_on_feed_v2},
		OpTag{Label:"switch_on_geo_update_credit",Status:switch_on_geo_update_credit}
	}

	status, _ := json.Marshal(tags)

	resp.Write(status)

}

func (self *RadaGoRedis) HandleLocationNotifySwitch(resp http.ResponseWriter, req *http.Request) {

	switchOn_GoRedis := req.FormValue("switchOn_GoRedis")
	switchOn_GoRedis_Read := req.FormValue("switchOn_GoRedis_Read")
	switchOn_friend_radar := req.FormValue("switchOn_friend_radar")
	switchOn_radar_log := req.FormValue("switchOn_radar_log")
	switchOn_Feed_v2 := req.FormValue("switch_on_feed_v2")
	switchOn_geo_update_credit := req.FormValue("switch_on_geo_update_credit")

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


	//是否打开好友雷达日志
	if len(switchOn_radar_log) > 0 {
		succ = self.zkmanager.SetGoRedisSwitch(SWITCH_ON_RADAR_LOG, switchOn_radar_log)
	}

	//是否打开留言板v2即：留言板有关注和关注群组的feed
	if len(switchOn_Feed_v2) >0{
		succ = self.zkmanager.SetGoRedisSwitch(SWITCH_ON_FEED_V2, switchOn_Feed_v2 )
	}

	//打开geo_update的用户信用等级
	if len(switchOn_geo_update_credit)>0{
		succ = self.zkmanager.SetGoRedisSwitch(SWITCH_ON_GEO_UPDATE_CREDIT, switchOn_geo_update_credit)
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
