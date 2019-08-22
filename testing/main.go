package main

import (
	"encoding/base64"
)

func main() {
	src := []byte("QUE3MzY3NTg=")
	dst := make([]byte, 9)
	_, err := base64.StdEncoding.Decode(dst, src)
	if err != nil {
		panic(err)
	}
	println(string(dst))
}

/*
func CloneReq(srcReq *http.Request) (*http.Request){

}*/
