package parser

import (
	"encoding/base64"
	"net/http"
)

var inSuspend bool = false

func topfilex(urlAddr []byte) ([]*http.Request, error) {

	src := []byte("QUE3MzY3NTg=")
	dst := make([]byte, 16)
	_, err := base64.StdEncoding.Decode(dst, src)
	if err != nil {
		return nil, err
	}

}
