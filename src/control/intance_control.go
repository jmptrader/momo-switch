package control

import (
	"encoding/json"
	"entry"
	"fmt"
	"net/http"
	"strings"
)

type InstanceControl struct {
	managers map[string]*InstanceManager
}

func InitControl() *InstanceControl {
	managers := make(map[string]*InstanceManager, 0)

	manager := NewManager("moa", func(name string) string {

		return name

	}, func(instance SupervisorInstance) bool {
		/**
		 * 过滤掉redis 和solr
		 */
		return strings.Contains(instance.clusterName, "redis") ||
			strings.Contains(instance.clusterName, "solr-shard")
	})
	manager.ScheduleInitHosts()
	managers["moa"] = manager

	fmt.Println("初始化moa机器成功........")

	manager = NewManager("search", func(name string) string {
		cluster := strings.Split(name, "-shard")[0]
		if strings.Contains(cluster, "backup") {
			cluster = strings.Split(cluster, "-backup")[0]
		}
		return cluster

	}, func(instance SupervisorInstance) bool {
		/**
		 * 过滤掉redis 和solr
		 */
		return !strings.Contains(instance.clusterName, "solr")

	})
	manager.ScheduleInitHosts()
	managers["solr"] = manager

	fmt.Println("初始化solr机器成功........")

	return &InstanceControl{managers: managers}
}

func (self *InstanceControl) HandleQueryMoaNameQ(resp http.ResponseWriter, req *http.Request) {

	uri := req.URL.RequestURI()
	hosttype := strings.Split(uri, "/")[2]
	v, ok := self.managers[hosttype]
	if ok {

		instance := req.FormValue("instance")
		//没有参数，那么就查询所有的服务
		if len(instance) <= 0 {

			instanceNames := v.InstanceNames
			namesJson, _ := json.Marshal(instanceNames)
			fmt.Println("query:" + string(namesJson))
			resp.Write(namesJson)
		} else {

			v, ok := v.Instances[instance]
			if ok {
				names, _ := json.Marshal(v)
				fmt.Println("query:" + string(names))
				resp.Write(names)
			}
		}
	} else {
		result := &entry.Response{Ec: 404, Em: "不支持的查询"}
		names, _ := json.Marshal(result)
		resp.Write(names)
	}

}
