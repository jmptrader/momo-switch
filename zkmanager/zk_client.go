package zkmanager

import (
	"fmt"
	"github.com/blackbeans/zk"
	"strconv"
	"strings"
	"time"
)

type ZKManager struct {
	session *zk.Session
}

func NewZKManager(zkhosts string) *ZKManager {
	if len(zkhosts) <= 0 {
		fmt.Println("使用默认zkhosts！|localhost:2181\n")
		zkhosts = "localhost:2181"
	} else {
		fmt.Printf("使用zkhosts:[%s]！\n", zkhosts)
	}

	conf := &zk.Config{Addrs: strings.Split(zkhosts, ",")[1:], Timeout: 5 * time.Second}

	ss, err := zk.Dial(conf)
	if nil != err {
		panic("连接zk失败..." + err.Error())
		return nil
	}

	return &ZKManager{session: ss}

}

func (self *ZKManager) SetGoRedisSwitch(path string, data string) bool {

	_, err := self.session.Set(path, []byte(data), -1)
	if nil != err {
		return false
	}

	fmt.Println("设置[" + path + "] 成功! [" + data + "]")
	return true
}

func (self *ZKManager) TraverseCreatePath(path string, defaultValue string) {

	blocks := strings.Split(path, "/")
	root := ""
	for i := 1; i < len(blocks); i++ {
		root = root + "/" + blocks[i]
		fmt.Println(root)
		exist, _, err := self.session.Exists(root, nil)
		//创建节点
		if !exist && nil == err {
			id, _ := self.session.Create(root, []byte(defaultValue), zk.CreatePersistent, zk.AclOpen)
			fmt.Printf("创建节点:[%s],id:[%s]\n", root, id)
		} else if nil != err {
			fmt.Println(err.Error())
		}

	}
}

func (self *ZKManager) Get(path string) bool {

	val, _, err := self.session.Get(path, nil)
	if nil != err {
		fmt.Println("获取[" + path + "] 失败!")
		return false
	}
	bv, _ := strconv.ParseBool(string(val))
	return bv
}
