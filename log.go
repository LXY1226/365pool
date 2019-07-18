package main

import (
	"fmt"
	"time"
)

func Logln(a ...interface{}) {
	fmt.Println(GetTimeStr(), a)
}

func Log(a ...interface{}) {
	fmt.Print(GetTimeStr(), a)
}

func GetTimeStr() string {
	t := time.Now()
	return fmt.Sprintf("[%d:%d:%d] ", t.Hour(), t.Minute(), t.Second())
}
