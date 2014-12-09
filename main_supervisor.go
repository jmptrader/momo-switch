package main

import (
	"fmt"
	"github.com/blackbeans/goquery"
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

func main() {
	hosts := []string{"task001.m6",
		"task002.m6",
		"task003.m6",
		"task004.m6",
		"task005.m6",
		"task006.m6",
		"task007.m6",
		"task008.m6",
		"task009.m6",
		"task010.m6",
	}
	store := make(map[string][]string, 20)
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
					}

				})

				arr, ok := store[instance.Name]
				if !ok {
					arr = make([]string, 0, 10)

				}
				arr = append(arr, v)
				store[instance.Name] = arr

			})
		}
	}

	for k, v := range store {
		label := k + "\t\t\t"
		for i, host := range v {
			label += host
			if i < len(v)-1 {
				label += ","
			}
		}

		fmt.Printf("* \t%s\n", label)
	}

}
