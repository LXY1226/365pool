package main

import (
	"365pool/dfpan"
	"math"
	"net/http"
	"os"
	"strconv"
)

const (
	downloadItems   = 8
	downloadThreads = 8
)

func main() {
	fin, err := os.OpenFile("test.bin", os.O_CREATE|os.O_WRONLY, os.ModeType)
	if err != nil {
		panic(err)
	}
	err = fin.Truncate(1048576 * 1024)
	if err != nil {
		panic(err)
	}
	dfpan.Parse([]byte("5adcczb414b67"))
}

func read_conf() {
	f, err := os.Open("conf/global.conf")
	if err != nil {
		conf := ask_conf()
		conf += 1
	}
	defer func() { _ = f.Close() }()
}
func bytesToSize(length int) string {
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

func download(r http.Request, f os.File) error {

	return nil
}
