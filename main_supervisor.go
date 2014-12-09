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

	hosts := []string{
		"xxx001.m6",
		"xxx002.m6",
		"xxx003.m6",
		"xxx004.m6",
		"xxx005.m6",
		"xxx006.m6",
		"xxx007.m6",
		"xxx008.m6",
		"xxx009.m6",
		"xxx010.m6",
		"xxx011.m6",
		"xxx012.m6",
		"xxx013.m6",
		"xxx014.m6",
		"xxx015.m6",
		"xxx016.m6",
		"xxx017.m6",
		"xxx018.m6",
		"xxx019.m6",
		"xxx020.m6",
		"xxx021.m6",
		"xxx022.m6",
		"xxx023.m6",
		"xxx024.m6",
		"xxx025.m6",
		"xxx026.m6",
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

				if instance.Status != "running" || instance.Name == "flume-agent" {
					return
				}

				arr, ok := store[instance.Name]
				if !ok {
					arr = make([]string, 0, 10)

				}
				arr = append(arr, v)
				store[instance.Name] = arr

			})
		}
	}

	table := "||||\n" +
		"|:--:|:--:|----|\n" +
		"| 服务名称|数量| 部署机器 |\n"

	for k, v := range store {
		label := ""
		for i, host := range v {
			label += host
			if i < len(v)-1 {
				label += ","
			}
		}

		table += fmt.Sprintf("|%s|%d|%s|\n", k, len(v), label)
	}

	fmt.Println(table)

}
