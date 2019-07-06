package main

import (
	"bufio"
	"fmt"
	"github.com/lxy1226/365pool/dfpan"
	"io"
	"math"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync/atomic"
	"time"
)

// uint up to 256
const (
	downloadItems   = 3
	downloadThreads = 8
	retryTimes      = 10
	pieceSize       = 320 * 1024
)

func main() {
	refreshChan := make(chan uint)
	reqs, err := dfpan.Parse([]byte("5adcczb414b67"))
	if err != nil {
		//Error
		fmt.Println("Parse Error: ", err)
	} else {
		go download(refreshChan, reqs, "D:/Download/", "")
	}
	time.Sleep(2 ^ 10*time.Hour)

}

func read_conf() {
	f, err := os.Open("conf/global.conf")
	if err != nil {
		conf := ask_conf()
		conf += 1
	}
	defer func() { _ = f.Close() }()
}
func bytesToSize(length uint64) string {
	// https://blog.csdn.net/a99361481/article/details/81751231
	var k = 1024 // or 1024
	var sizes = []string{" B", "KB", "MB", "GB"}
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

func download(taskrefreshChan chan uint, r []http.Request, dir string, filename string) {
	var poss [downloadThreads]uint64
	refreshChan := make(chan uint)
	var size uint64
	var reqs []http.Request
	client := http.Client{}
	client.Timeout = 10 * time.Second
	for i := 0; i < len(r); i++ {
		req := r[i]
		req.Method = "HEAD"
		resp, err := client.Do(&req)
		if err == nil {
			asize, _ := strconv.ParseInt(resp.Header.Get("Content-Length"), 10, 64)
			nsize := uint64(asize)
			if size == 0 {
				size = nsize
			} else if size != nsize {
				fmt.Println("Incorrect Requests")
			}
			if filename == "" {
				filename = resp.Header.Get("Content-Disposition")
				filename = strings.Split(filename, "filename=")[1]
				println(filename)
			}
			reqs = append(reqs, r[i])
		}
	}
	if reqs == nil {
		fmt.Println("No Usable URL")
		taskrefreshChan <- 1
		return
	}
	var ranges [downloadThreads][2]uint64
	if size != 0 {
		var blockCount, threadCount, blocksextra, pos uint64
		blockCount = uint64(size / pieceSize)
		threadCount = uint64(blockCount / downloadThreads)
		blocksextra = blockCount % downloadThreads
		for i := 0; i < downloadThreads-1; i++ {
			ranges[i][0] = pos
			poss[i] = pos
			pos += threadCount * pieceSize
			if blocksextra > 0 {
				pos += pieceSize
				blocksextra--
			}
			ranges[i][1] = pos
		}
		ranges[downloadThreads-1][0] = pos
		ranges[downloadThreads-1][1] = size
	} else {
		fmt.Println("Download for 0B?")
		taskrefreshChan <- 1
		return
	}
	f, err := os.Create(dir + filename)
	if err != nil {
		panic(err)
	}
	err = f.Truncate(int64(size))
	if err != nil {
		panic(err)
	}
	for id := uint(0); id < downloadThreads; id++ {
		req := reqs[int(rand.Float32()*float32(len(reqs)))]
		go goDownload(id, refreshChan, req, &ranges[id][0], &ranges[id][1], f)
	}
	ticker := time.NewTicker(1 * time.Second)
	for {
		select {
		case /*id :=*/ <-refreshChan:
			//TODO Add Dynamic Part Alloc
			println("go exited")
		case <-ticker.C:
			var speed uint64
			for id := 0; id < downloadThreads; id++ {
				speed += ranges[id][0] - poss[id]
				poss[id] = ranges[id][0]
			}
			print("\r" + bytesToSize(speed))
		}
	}
}

func CopyReq(req http.Request) *http.Request {
	requestNew := new(http.Request)
	*requestNew = req
	fmt.Printf("%p %p", req, *requestNew)
	return requestNew
}

/*
func CopyReq(r *http.Request)  http.Request{
	requestNew := http.Request{}
	requestNew.Header = r.Header
	requestNew.URL = r.URL
	requestNew.Method = r.Method
	return requestNew
}
*/

func goDownload(id uint, refreshChan chan uint, basereq http.Request, pos *uint64, end *uint64, f *os.File) {
	fmt.Printf("[%d]Started %s-%s\n", id, bytesToSize(*pos), bytesToSize(*end))
	client := http.Client{}
	for {
		for i := 0; i < retryTimes+1; i++ {
			//req := CopyReq(&basereq)
			req := basereq
			newreq := CopyReq(req)
			newreq.Method = "GET"
			newreq.Header.Set("Range", fmt.Sprintf("bytes=%d-%d", *pos, *end))
			//resp, err := client.Do(&req)
			resp, err := client.Do(newreq)
			raw := resp.Body
			defer raw.Close()
			reader := bufio.NewReaderSize(raw, pieceSize)
			buff := make([]byte, pieceSize)
			for {
				nr, er := reader.Read(buff)
				if nr > 0 {
					nw, ew := f.WriteAt(buff[0:nr], int64(*pos))
					fmt.Printf("[%d]WriteAt %d %X\n", id, *pos, buff[0:7])
					if nw > 0 {
						atomic.StoreUint64(pos, *pos+uint64(nr))
						if *pos >= *end {
							refreshChan <- id
							return
						}
					}
					if ew != nil {
						err = ew
						break
					}
					if nr != nw {
						err = io.ErrShortWrite
						break
					}
				}
				if er != nil {
					if er != io.EOF {
						err = er
					}
					break
				}
			}
			if err != nil {
				fmt.Println("Download Error: ", err)
			}

		}
	}
}
