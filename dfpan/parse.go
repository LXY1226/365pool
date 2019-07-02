package dfpan

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"strconv"
)

//license_code = "7c9d535eb56d2690dc09b760574c9a4d" @43
//mac          = "03000200-0400-0500-0006-000700080009" @47

var dowords = [9][]byte{}
var inited = false

/*{
"0987c0e589ed",
"0eebbff94611",
"a2fadc73e4ea",
"6f6c21a4065c",
"48c028081dd8",
"ae90f5f60751",
"85c18a624204",
}*/

func Parse(id []byte) *http.Request {
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
			word.WriteString(hex.EncodeToString(md5w.Sum(nil)))
			var doword bytes.Buffer
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
	var url bytes.Buffer
	url.Write([]byte("http://page2.dfpan.com/view?module=downLoader&vr=2.9.4&fixufid="))
	url.Write(id)
	url.Write([]byte("&action=download&licence=7c9d535eb56d2690dc09b760574c9a4d&dowords="))
	url.Write(dowords[int(rand.Float32()*float32(len(dowords)))])
	url.Write([]byte("&ucode="))
	url.WriteString(ucode)
	url.Write([]byte("&mac=03000200-0400-0500-0006-000700080009"))
	client := &http.Client{}
	request, err := http.NewRequest("GET", url.String(), nil)
	request.Header.Add("User-Agent", "Mozilla/4.0 (compatible; MSIE 8.0; Windows NT 6.0) membership/2 YunDown/2.9.4")
	request.Header.Add("Referer", "http://page2.dfpan.com/fs/"+string(id))
	request.Header.Add("Continue", "Continue")
	request.Header.Add("Host", "page2.dfpan.com")
	request.Header.Add("Accept", "text/html, image/gif, image/jpeg, *; q=.2, */*; q=.2")
	request.Header.Add("Connection", "keep-alive")
	if err != nil {
		log.Panicln("Internal Error: dfpan.Parse", err)
	}
	response, _ := client.Do(request)
	defer response.Body.Close()
	body, err := ioutil.ReadAll(response.Body)
	downUrl := bytes.Split(body, []byte("downUrl:"))[1]
	request, err = http.NewRequest("GET", string(downUrl), nil)
	request.Header.Add("User-Agent", "Mozilla/4.0 (compatible; MSIE 8.0; Windows NT 6.0) membership/2 YunDown/2.9.4")
	request.Header.Add("Referer", "http://page2.dfpan.com/fs/"+string(id)+"/")
	return request
}
