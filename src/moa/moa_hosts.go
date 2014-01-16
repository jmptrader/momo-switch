package moa

import (
	"container/list"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

const (
	MOA_HOST_URL  = "http://kratos.wemomo.com/api/zabbix/hostgroups/moa"
	SCHEDULE_TIME = 60 * time.Second
)

/**
 * moa实例
 */
type MoaInstance struct {
	Host       string `json:"host"`       //当前机器名
	Name       string `json:"name"`       //服务名称
	RestartUrl string `json:"restarturl"` //重启url
	StopUrl    string `json:"stopurl"`    // 关闭url
	Status     string `json:"status"`     //当前状态
}

type MoaInStanceManager struct {
	//用于存放服务名称到moa实例的映射
	Instances map[string]*list.List
}

/**
 *
 * 定时拉取Moa的机器
 */
func (manager *MoaInStanceManager) ScheduleInitHosts() {

	syncMoaHosts()
	//定时任务开启，来获取moa的hosts
	timer := time.NewTicker(SCHEDULE_TIME)
	go func() {
		for {
			select {
			case <-timer.C:
				syncMoaHosts()
			}

		}
	}()
}

func syncMoaHosts() {
	resp, err := http.Get(MOA_HOST_URL)
	if nil != err {
		fmt.Println("获取MOA机器失败...." + err.Error())
		return
	}

	data, err := ioutil.ReadAll(resp.Body)
	if nil != err {
		fmt.Println("获取MOA机器失败...." + err.Error())
		return
	}

	//解析出机器列表
	var hosts []string
	json.Unmarshal(data, &hosts)

	if nil != hosts {
		//如果得到了hosts
		for _, v := range hosts {
			fmt.Println("机器名:" + v)
		}

	}

}
