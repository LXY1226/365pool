package dfpan

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"github.com/lxy1226/365pool/types"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"strconv"
)

//license_code = "7c9d535eb56d2690dc09b760574c9a4d" @43
//mac          = "03000200-0400-0500-0006-000700080009" @47

var dowords = [9][]byte{}
var inited = false

func Parse(id []byte) ([]http.Request, error) {
	if !inited {
		resp, err := http.Get("http://page2.dfpan.com/downloader/webip.jsp")
		if err != nil {
			log.Fatalln("Get IP Error", err)
		}
		ip, _ := ioutil.ReadAll(resp.Body)
		ip = bytes.Split(ip, []byte("\r"))[0]
		_ = resp.Body.Close()
		var word bytes.Buffer
		for i := 1; i < 10; i++ {
			md5w := md5.New()
			md5w.Write(ip)
			md5w.Write([]byte("7c9d535eb56d2690dc09b760574c9a4dkieliOAwii*&^543uy(t<bvfe?PQZW"))
			md5w.Write(word.Bytes())
			word.Reset()
			word.Grow(32)
			word.WriteString(hex.EncodeToString(md5w.Sum(nil)))
			var doword bytes.Buffer
			doword.Grow(12)
			doword.Write(word.Bytes()[0:7])
			doword.WriteString(strconv.FormatInt(int64(i), 10))
			doword.Write(word.Bytes()[8:12])
			dowords[i-1] = doword.Bytes()
		}
		inited = true
	}
	md5w := md5.New()
	md5w.Write(id)
	md5w.Write([]byte("kieliOAwii*&^543uy(t<bvfe?PQZW"))
	sum := md5w.Sum(nil)
	ucode := hex.EncodeToString(sum[4:8])
	var tmpurl bytes.Buffer
	tmpurl.Write([]byte("http://page2.dfpan.com/view?module=downLoader&vr=2.9.4&fixufid="))
	tmpurl.Write(id)
	tmpurl.Write([]byte("&action=download&licence=7c9d535eb56d2690dc09b760574c9a4d&dowords="))
	tmpurl.Write(dowords[int(rand.Float32()*float32(len(dowords)))])
	tmpurl.Write([]byte("&ucode="))
	tmpurl.WriteString(ucode)
	tmpurl.Write([]byte("&mac=03000200-0400-0500-0006-000700080009"))
	client := &http.Client{}
	request := http.Request{}
	request.URL, _ = url.Parse(tmpurl.String())
	request.Header = http.Header{}
	request.Header.Add("User-Agent", "Mozilla/4.0 (compatible; MSIE 8.0; Windows NT 6.0) membership/2 YunDown/2.9.4")
	request.Header.Add("Referer", "http://page2.dfpan.com/fs/"+string(id))
	response, _ := client.Do(&request)
	defer response.Body.Close()
	body, err := ioutil.ReadAll(response.Body)
	if len(bytes.Split(body, []byte("downUrl:"))) != 2 {
		return nil, &types.MyError{string(body)}
	}
	if err != nil {
		log.Panicln("Internal Error: dfpan.Parse", err)
	}
	downUrl := bytes.Split(body, []byte("downUrl:"))[1]
	request.URL, _ = url.Parse(string(downUrl))
	requests := []http.Request{request}
	return requests, nil
}
