package main

import (
	"bufio"
	"fmt"
	"github.com/lxy1226/365pool/dfpan"
	"github.com/lxy1226/365pool/types"
	"log"
	"math"
	"net/http"
	"os"
	"strconv"
	"sync/atomic"
	"time"
)
// uint8 up to 256
const (
	downloadItems   = 3
	downloadThreads = 8
	retryTimes      = 10
)

func main() {
	/*fin, err := os.OpenFile("test.bin", os.O_CREATE|os.O_WRONLY, os.ModeType)
	if err != nil {
		panic(err)
	}
	err = fin.Truncate(1048576 * 1024)
	if err != nil {
		panic(err)
	}*/
	refreshChan := make(chan uint8)
	reqs, err := dfpan.Parse([]byte("5adcczb414b67"))
	if err != nil {
		//Error
		fmt.Println("Parse Error: ", err)
	} else {
		err := download(reqs, "D:/Download/", "")
		if err != nil {
			fmt.Println("Download Error: ", err)
		}
	}
}

func read_conf() {
	f, err := os.Open("conf/global.conf")
	if err != nil {
		conf := ask_conf()
		conf += 1
	}
	defer func() { _ = f.Close() }()
}
func bytesToSize(length uint32) string {
	// https://blog.csdn.net/a99361481/article/details/81751231
	var k = 1024 // or 1024
	var sizes = []string{" B", "KB", "MB", "GB", "TB"}
	if length == 0 {
		return "0 B"
	}
	i := math.Floor(math.Log(float64(length)) / math.Log(float64(k)))
	r := float64(length) / math.Pow(float64(k), i)
	return strconv.FormatFloat(r, 'f', 3, 64) + " " + sizes[int(i)]
}

func ask_conf() int {
	return 0
}

func download(r []*http.Request, dir string, filename string) error {
	speeds := [downloadThreads]uint32{}
	refreshChan := make(chan uint8)
	var size int64
	var reqs []*http.Request
	client := &http.Client{}
	client.Timeout = 10
	for i := 0; i < len(r); i++ {
		req := *r[i]
		req.Method = "HEAD"
		resp, err := client.Do(*req)
		if err == nil {
			nsize, _ = strconv.ParseInt(resp.Header.Get("Content-Length"), 10, 64)
			if size == 0 {
				size = nsize
			} else if size != nsize {
				return &types.MyError{"Requests Size not match"}
			}
			reqs = append(reqs, r[i])
		}
	}

	for i := 0; i < downloadThreads; i++ {
		start = 0
		go goDownload()
	}
	select {
	case /*i := */ <-refresh_chan:

	}
	go add_speed(&speeds[0])
	time.Sleep(1 * time.Second)
	fmt.Println(speeds)
	return nil
}

func goDownload(id int, refreshChan chan uint8, basereq *http.Request, pos uint64, size uint64, fIo *bufio.Writer) {
	req := *basereq
	//https://gocn.vip/question/666
	req.Header.Add("Range", fmt.Sprintf("bytes=%d-%d", start, start+size))
	//https://blog.csdn.net/a99361481/article/details/81751231
	client := &http.Client{}
	client.Timeout = 10 * time.Second
	resp, _ := client.Do(req)
	completed := false
	for !completed {
		if err != nil {
			for i := 0, i < retryTimes, i++ {
				fmt.println("Download Error: ")
			}
			
		}
		atsomic.StoreUint32(speed, 3)
	}
	i -> refresh_chan
}

func add_speed(speed *uint32) {
	
}
