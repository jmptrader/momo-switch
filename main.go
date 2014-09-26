package main

// import "github.com/jmhodges/levigo"
import (
	"fmt"
	// "os"
	"encoding/json"
	"net/http"
	_ "net/http/pprof"
	"strings"
)

// import "code.google.com/p/goprotobuf/proto"

const (
	DIC    = "0123456789bcdefghjkmnpqrstuvwxyz"
	BITNUM = 30
)

func main() {

	// fmt.Println(strings.Split("solr-shard20-8985", "-shard")[0])

	// opt := levigo.NewOptions()
	// opt.SetCreateIfMissing(true)
	// opt.SetBlockSize(32 * 1024)
	// opt.SetCompression(levigo.SnappyCompression)

	// fileinfo, err := os.Stat("/data/demo")
	// if os.IsNotExist(err) || !fileinfo.IsDir() {
	// 	os.MkdirAll("/data/demo", os.ModePerm)
	// 	fmt.Println("创建目录成功")
	// } else if nil != err {
	// 	fmt.Println(err)
	// 	return
	// }

	// db, err := levigo.Open("/data/demo/db", opt)
	// if nil != err {
	// 	panic(err)
	// }

	// wopt := levigo.NewWriteOptions()
	// batchWr := levigo.NewWriteBatch()

	// batchWr.Put([]byte("hello"), []byte("value"))
	// batchWr.Put([]byte("hello1"), []byte("value1"))
	// batchWr.Put([]byte("hello2"), []byte("value2"))

	// db.Write(wopt, batchWr)
	// defer wopt.Close()

	// ropt := levigo.NewReadOptions()
	// ropt.SetFillCache(false)
	// data, _ := db.Get(ropt, []byte("hello"))

	// fmt.Printf("%s", string(data))

	// it := db.NewIterator(ropt)
	// it.Seek([]byte("hello1"))
	// fmt.Printf("iterator:%s\n", it.Valid())
	// for it.Valid() {

	// 	fmt.Printf("%s|%s\n", string(it.Key()), string(it.Value()))
	// 	it.Next()
	// }

	// defer ropt.Close()
	//
	// fmt.Println(decodeLatLng("wqh5tf"))

	http.HandleFunc("/q", decodeGeoCode)
	http.ListenAndServe(":19870", nil)
}

type Location struct {
	Lat  float64 `json:"lat"`
	Lng  float64 `json:"lng"`
	Succ int32   `json:"succ"`
}

func decodeGeoCode(rw http.ResponseWriter, req *http.Request) {
	code := req.FormValue("code")
	loc := &Location{Lat: 0.0, Lng: 0.0}
	if len(code) <= 0 {
		loc.Succ = 500
		fmt.Println("不能为空！")
		return
	}

	lat, lng := decodeLatLng(code)
	// fmt.Printf("%f,%f", lat, lng)
	loc.Lat = lat
	loc.Lng = lng
	data, _ := json.Marshal(loc)

	rw.Write(data)

}

func decodeLatLng(geocode string) (lat, lng float64) {

	var codebit int64 = 0
	var i int = 0
	reader := strings.NewReader(geocode)
	for ; i < len(geocode); i++ {
		b, _, _ := reader.ReadRune()

		c := fmt.Sprintf("%c", b)
		idx := strings.Index(DIC, c)
		// fmt.Printf("%s,%d\n", c, idx)
		codebit = (codebit << 5) | int64((idx & 0x1F))

	}
	//fmt.Println(codebit)

	var latbit int64 = 0
	var lngbit int64 = 0

	//解析经纬度bit
	//
	//
	for i = 0; i < BITNUM; i++ {
		if i%2 == 0 {
			lngbit = (lngbit << 1) | ((codebit >> uint32((BITNUM - i - 1))) & 0x01)
		} else {
			latbit = (latbit << 1) | ((codebit >> uint32((BITNUM - i - 1))) & 0x01)
		}
	}

	// fmt.Printf("%d,%d\n", latbit, lngbit)
	// 开始计算一下中心点的经纬度

	return splitNum(15, latbit, -90.0, 90.0), splitNum(15, lngbit, -180.0, 180.0)
}

func splitNum(bitIdx int, bit int64, baseNeg float64, basePos float64) float64 {

	var tmpL float64 = baseNeg
	var tmpR float64 = basePos

	var tmpMid float64 = (basePos + baseNeg) / 2
	// fmt.Println(tmpMid)
	for i := 0; i < bitIdx; i++ {
		val := (bit >> uint32((bitIdx - i - 1))) & 0x01
		if val == 1 {
			tmpL = tmpMid
		} else {
			tmpR = tmpMid

		}
		tmpMid = (tmpL + tmpR) / 2
		// fmt.Printf("[%f,%f]\t", tmpL, tmpR)
	}

	return tmpMid

}
