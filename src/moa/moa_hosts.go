package moa

import (
	"container/list"
	"encoding/json"
	"fmt"
	"github.com/blackbeans/goquery"
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
	Info       string `json:"info"`       //启动信息
}

type MoaInStanceManager struct {
	//用于存放服务名称到moa实例的映射
	Instances     *map[string]list.List
	InstanceNames *list.List
}

/**
 *初始化manager
 *
 */
func NewManager() *MoaInStanceManager {
	manager := &MoaInStanceManager{}
	manager.ScheduleInitHosts()

	return manager
}

/**
 *
 * 定时拉取Moa的机器
 */
func (self *MoaInStanceManager) ScheduleInitHosts() {

	self.syncMoaHosts()
	recover()

	//定时任务开启，来获取moa的hosts
	timer := time.NewTicker(SCHEDULE_TIME)
	go func() {
		for {
			select {
			case <-timer.C:
				self.syncMoaHosts()
			}

		}
	}()
}

func (self *MoaInStanceManager) syncMoaHosts() {
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
	defer resp.Body.Close()

	instances := make(map[string]list.List, 50)
	//解析出机器列表
	var hosts []string
	json.Unmarshal(data, &hosts)

	if nil != hosts {
		//如果得到了hosts
		for _, v := range hosts {
			baseUrl := "http://" + v + ":9001"
			doc, err := goquery.NewDocument(baseUrl)
			if nil == err {
				doc.Find("table tbody tr").Each(func(i int, s *goquery.Selection) {

					instance := &MoaInstance{Host: v}
					s.Find("td").Each(func(j int, ss *goquery.Selection) {
						if j > 2 {
							return
						}

						switch j {
						case 0:
							instance.Status = ss.Children().Text()
						case 1:
							instance.Info = ss.Children().Text()
						case 2:
							instance.Name = ss.Children().Text()

						}

					})

					instance.RestartUrl = baseUrl + "/index.html?processname=" + instance.Name + "&amp;action=restart"
					instance.StopUrl = baseUrl + "/index.html?processname=" + instance.Name + "&amp;action=stop"
					v, ok := instances[instance.Name]
					if !ok {
						v = *list.New()
						instances[instance.Name] = v

					}
					// //将该节点推送
					v.PushFront(instance)

				})
			}
		}

		names := list.New()

		for k, v := range instances {
			names.PushBack(k)

			jsonStr, _ := json.Marshal(v)
			fmt.Println(string(jsonStr))
		}

		self.Instances = &instances
		self.InstanceNames = names
	}

}
