package lighttp

import (
	"errors"
	"net"
	"net/url"
)

var dnsResolver net.Resolver

type LightReq struct {
	u url.URL
	header []byte

}

func GetHttp(u url.URL) (error) {
	dialer := net.Dialer{Resolver: &dnsResolver}
	conn, err := dialer.Dial("tcp", u.Host)
	if err != nil {return err}
	conn.
}