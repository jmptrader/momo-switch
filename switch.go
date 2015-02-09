package main

import (
	"flag"
	"fmt"
	"momo-switch/control"
	"momo-switch/location"
	"momo-switch/zkmanager"
	"net/http"
)

func main() {

	fmt.Println("-zkhosts=localhost:2181 来定义zookeeper的地址!")
	zkhosts := flag.String("zkhosts", "vm-search-001:2181,vm-search-002:2181,vm-search-003:2181", "输入zookeeper地址...请用逗号分隔")
	flag.Parse()

	control := control.InitControl()

	zkmanager := zkmanager.NewZKManager(*zkhosts)

	radaGoRedis := location.NewRadaGoRedis(zkmanager)

	http.HandleFunc("/switch/location_noftiy/conf", radaGoRedis.HandleLocationNotifySwitch)
	http.HandleFunc("/switch/location_noftiy/q", radaGoRedis.HandleLocationNotifySwitchQ)
	http.HandleFunc("/switch/moa/q_instances", control.HandleQueryMoaNameQ)
	http.HandleFunc("/switch/solr/q_instances", control.HandleQueryMoaNameQ)
	http.HandleFunc("/switch/flume/q_instances", control.HandleQueryMoaNameQ)
	http.HandleFunc("/switch/task/q_instances", control.HandleQueryMoaNameQ)
	http.HandleFunc("/switch/trade/q_instances", control.HandleQueryMoaNameQ)
	http.ListenAndServe(":7979", nil)
}
