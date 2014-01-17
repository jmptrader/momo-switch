package main

import "github.com/jmhodges/levigo"
import (
	"fmt"
	"os"
	"strings"
)

// import "code.google.com/p/goprotobuf/proto"

func main() {
	fmt.Println(strings.Split("feed-solr-shard10-9770", "-shard")[0])

	opt := levigo.NewOptions()
	opt.SetCreateIfMissing(true)
	opt.SetBlockSize(32 * 1024)
	opt.SetCompression(levigo.SnappyCompression)

	fileinfo, err := os.Stat("/data/demo")
	if os.IsNotExist(err) || !fileinfo.IsDir() {
		os.MkdirAll("/data/demo", os.ModePerm)
		fmt.Println("创建目录成功")
	} else if nil != err {
		fmt.Println(err)
		return
	}

	db, err := levigo.Open("/data/demo/db", opt)
	if nil != err {
		panic(err)
	}

	wopt := levigo.NewWriteOptions()
	batchWr := levigo.NewWriteBatch()

	batchWr.Put([]byte("hello"), []byte("value"))
	batchWr.Put([]byte("hello1"), []byte("value1"))
	batchWr.Put([]byte("hello2"), []byte("value2"))

	db.Write(wopt, batchWr)
	defer wopt.Close()

	ropt := levigo.NewReadOptions()
	ropt.SetFillCache(false)
	data, _ := db.Get(ropt, []byte("hello"))

	fmt.Printf("%s", string(data))

	it := db.NewIterator(ropt)
	it.Seek([]byte("hello1"))
	fmt.Printf("iterator:%s\n", it.Valid())
	for it.Valid() {

		fmt.Printf("%s|%s\n", string(it.Key()), string(it.Value()))
		it.Next()
	}

	defer ropt.Close()

}
