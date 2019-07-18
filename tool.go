package main

import (
	"math"
	"strconv"
)

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
