package main

import (
	"bufio"
	"fmt"
	"github.com/lxy1226/365pool/dfpan"
	"io/ioutil"
	"math/rand"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)

// uint up to 256
const (
	downloadItems = 3
	pararrel      = 8
	retryTimes    = 10
	splitSize     = 100 * 320 * 1024 //32MiB
	buffSize      = 32 * 1024
)

type downTask struct {
	id      *string
	downReq *http.Request
	start   uint64
	size    *uint64
	upURL   *url.URL
}

var client = http.Client{}

func main() {
	taskChan := make(chan *downTask)
	id := "0fb1_12029127133196364"
	go initDown(taskChan)
	initOnedrive()
	reqs, err := dfpan.Parse([]byte(id))
	if err != nil {
		//Error
		fmt.Println("Parse Error: ", err)
	} else {
		download(&id, taskChan, reqs, "/Guomoo/")
	}
	time.Sleep(2 ^ 10*time.Hour)
}

func initDown(taskChan chan *downTask) {

	var works [pararrel]*downTask
	var refreshChan chan uint8
	mux := new(sync.Mutex)
	for id := uint8(0); id < pararrel; id++ {
		works[id] = <-taskChan
		go goTask(id, refreshChan, works[id], mux)
	}
	ok := true
	for {
		id := <-refreshChan
		works[id], ok = <-taskChan
		if !ok {
			fmt.Println("Successfully Finished")
			os.Exit(0)
		}
		go goTask(id, refreshChan, works[id], mux)
	}
}

func download(id *string, taskChan chan *downTask, r []*http.Request, dir string) {
	var reqs []*http.Request
	client := http.Client{}
	client.Timeout = 10 * time.Second
	size := uint64(0)
	var filename string
	for i := 0; i < len(r); i++ {
		req := r[i]
		req.Method = http.MethodHead
		resp, err := client.Do(req)
		if err == nil {
			asize, _ := strconv.ParseInt(resp.Header.Get("Content-Length"), 10, 64)
			nsize := uint64(asize)
			if size == 0 {
				size = nsize
			} else if size != nsize {
				fmt.Println("Incorrect Requests")
			}
			filename = resp.Header.Get("Content-Disposition")
			filename = strings.Split(filename, "filename=")[1]
			println(filename)
			reqs = append(reqs, r[i])
		}
		req.Method = http.MethodGet
	}
	if reqs == nil {
		fmt.Println("No Usable URL")
		return
	}
	upURL := mkUploadRequest(dir + filename)
	if size != 0 {
		pos := uint64(0)
		for pos = uint64(0); pos < size; pos += splitSize {
			task := downTask{
				id:      id,
				downReq: reqs[int(rand.Float32()*float32(len(reqs)))],
				start:   pos,
				size:    &size,
				upURL:   upURL,
			}
			taskChan <- &task
		}
	} else {
		fmt.Println("Download for 0B?")
		return
	}
}

func goTask(id uint8, refreshChan chan uint8, task *downTask, mux *sync.Mutex) {
	for {
		for i := 0; i < retryTimes+1; i++ {
			end := task.start + splitSize - 1
			if end > *task.size {
				end = *task.size - 1
			}
			Logln(fmt.Sprintf("%d %s %d-%d/%d %d", id, *task.id, task.start, end, *task.size, end-task.start+1))
			mux.Lock()
			task.downReq.Header.Set("Range", fmt.Sprintf("bytes=%d-%d", task.start, end))
			downResp, derr := client.Do(task.downReq)
			mux.Unlock()
			if derr != nil {
				fmt.Println("Download Error: ", derr)
			}
			raw := downResp.Body
			reader := bufio.NewReaderSize(raw, 32*1024*1024)
			upReq := http.Request{
				Method:        http.MethodPut,
				URL:           task.upURL,
				Header:        map[string][]string{"Content-Range": {fmt.Sprintf("bytes %d-%d/%d", task.start, end, task.size)}},
				Body:          ioutil.NopCloser(reader),
				ContentLength: int64(end - task.start + 1),
				Host:          task.upURL.Host,
			}
			_, uerr := client.Do(&upReq)
			if uerr != nil {
				fmt.Println("Upload Error: ", uerr)
			}
			Logln(fmt.Sprintf("%d %s %d-%d/%d %d", id, *task.id, task.start, end, *task.size, end-task.start+1))
			refreshChan <- id
		}
	}
}
