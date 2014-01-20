package control

import (
	"encoding/json"
	"fmt"
	"github.com/blackbeans/goquery"
	"io/ioutil"
	"net/http"
	"sort"
	"time"
)

const (
	BAISC_URL     = "http://kratos.wemomo.com/api/zabbix/hostgroups/"
	SCHEDULE_TIME = 60 * time.Second
)

/**
 * supervisor实例
 */
type SupervisorInstance struct {
	Host        string `json:"host"` //当前机器名
	Name        string `json:"name"` //服务名称
	clusterName string
	RestartUrl  string `json:"restarturl"` //重启url
	StopUrl     string `json:"stopurl"`    // 关闭url
	Status      string `json:"status"`     //当前状态
	Info        string `json:"info"`       //启动信息
}

// /**
//  * 实例名称过滤
//  */
// type IInstanceFilter interface {
// 	Filter(instance SupervisorInstance) bool
// }

type InstanceManager struct {
	filter     func(instance SupervisorInstance) bool
	namefilter func(name string) string
	hostType   string //host类型
	//用于存放服务名称到moa实例的映射
	Instances     map[string][]SupervisorInstance
	InstanceNames []string
}

/**
 *初始化manager
 *
 */
func NewManager(hostType string, namefilter func(name string) string,
	filter func(instance SupervisorInstance) bool) *InstanceManager {
	manager := &InstanceManager{}
	manager.hostType = hostType
	manager.filter = filter
	manager.namefilter = namefilter
	manager.ScheduleInitHosts()
	return manager
}

/**
 *
 * 定时拉取Moa的机器
 */
func (self *InstanceManager) ScheduleInitHosts() {

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

func (self *InstanceManager) syncMoaHosts() {
	resp, err := http.Get(BAISC_URL + self.hostType)
	if nil != err {
		fmt.Printf("获取[%s]机器失败....%s\n", self.hostType, err.Error())
		return
	}

	data, err := ioutil.ReadAll(resp.Body)
	if nil != err {
		fmt.Printf("获取[%s]机器失败....%s\n", self.hostType, err.Error())
		return
	}
	defer resp.Body.Close()

	instances := make(map[string][]SupervisorInstance)
	//解析出机器列表
	var hosts []string
	json.Unmarshal(data, &hosts)
	fmt.Printf("[%s] hosts:%s, hosts arr len:%n\n", self.hostType, string(data), len(hosts))
	if nil != hosts {
		//如果得到了hosts
		for _, v := range hosts {
			baseUrl := "http://" + v + ":9001"
			doc, err := goquery.NewDocument(baseUrl)
			if nil == err {
				doc.Find("table tbody tr").Each(func(i int, s *goquery.Selection) {

					instance := SupervisorInstance{Host: v}
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
							instance.clusterName = self.namefilter(instance.Name)
						}

					})

					//拼接info信息
					instance.Info = instance.Name + "|" + instance.Info

					if len(instance.Name) <= 0 || self.filter(instance) {
						return
					}

					instance.RestartUrl = baseUrl + "/index.html?processname=" + instance.Name + "&amp;action=restart"
					instance.StopUrl = baseUrl + "/index.html?processname=" + instance.Name + "&amp;action=stop"
					v, ok := instances[instance.clusterName]
					if !ok {
						v = make([]SupervisorInstance, 0)

					}

					// //将该节点推送
					instances[instance.clusterName] = append(v, instance)
					// jsonStr, _ := json.Marshal(instance)
					//fmt.Println("______________" + string(jsonStr) + "-----------------" + strconv.Itoa(v.Len()))

				})
			} else {
				fmt.Println(err.Error())
			}
		}

		names := make([]string, 0)

		for k, _ := range instances {
			names = append(names, k)
			// jsonStr, _ := json.Marshal(v.Front().Value)
			// fmt.Println(k + "+++++++++++++++++" + strconv.Itoa(v.Len()) + "--------------------" + string(jsonStr))
		}

		sort.Strings(names)
		self.Instances = instances
		self.InstanceNames = names
	}

}
