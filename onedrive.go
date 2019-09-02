package main

import (
	"fmt"
	"github.com/valyala/fastjson"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
)

type passport string

type oPassport struct {
}

const defpath = "https://guanghou-my.sharepoint.com/_api/v2.0"

/*
const (
	api_url = "https://graph.microsoft.com/v1.0"
	oauth_url = "https://login.microsoftonline.com/common/oauth2/v2.0" //
	client_secret = "fpvbADS606?$?)vdrHRKZ90"
	client_id = "c604ab67-5104-4aa0-accc-28c2e64b003d"
	OAUTH_REDIRECT_URI =
	"client_id":     "c604ab67-5104-4aa0-accc-28c2e64b003d",
	"client_secret": "fpvbADS606?$?)vdrHRKZ90",
	"redirect_uri":  "http://localhost:8000/callback",
	"resource":      "https%3A%2F%2Fgitaccuacnz2-my.sharepoint.com%2F",
)
https://login.microsoftonline.com/common/oauth2/v2.0/authorize?
client_id=c604ab67-5104-4aa0-accc-28c2e64b003d
&scope=offline_access+files.readwrite.all
&response_type=code
&redirect_uri=http://localhost:8000/callback

func login()
*/

var req http.Request
var access_token string

func initOnedrive() {
	refreshToken()
}

func refreshToken() {
	refreshtoken, err := ioutil.ReadFile("refreshtoken")
	if err != nil {
		refreshtoken = mkLogin()
	}
	refresh_token := string(refreshtoken)
	go func() {
		for {
			resp, err := http.PostForm("https://login.microsoftonline.com/common/oauth2/v2.0/token", map[string][]string{
				"client_id":     {"c604ab67-5104-4aa0-accc-28c2e64b003d"},
				"client_secret": {"fpvbADS606?$?)vdrHRKZ90"},
				"refresh_token": {refresh_token},
				"grant_type":    {"refresh_token"},
				"redirect_uri":  {"http://localhost:8000/callback"},
			})
			Logln("Onedrive: Token Refresh")
			res, _ := ioutil.ReadAll(resp.Body)
			if err != nil || resp.StatusCode != 200 {
				fmt.Println("Onedrive: get Token Error")
				os.Exit(1)
			}
			defer resp.Body.Close()
			refresh_token = fastjson.GetString(res, "refresh_token")
			access_token = "Bearer " + fastjson.GetString(res, "access_token")
			time.Sleep(30 * time.Minute)
		}
	}()
}

func mkLogin() []byte {
	fmt.Println("URL: https://login.microsoftonline.com/common/oauth2/v2.0/authorize?client_id=c604ab67-5104-4aa0-accc-28c2e64b003d&scope=offline_access%20files.readwrite.all&response_type=code&redirect_uri=http://localhost:8000/callback")
	var auth_code string
	_, _ = fmt.Scanln(&auth_code)
	code := strings.Split(strings.Split(auth_code, "&")[0], "code=")[1]
	resp, _ := http.PostForm("https://login.microsoftonline.com/common/oauth2/v2.0/token", map[string][]string{
		"client_id":     {"c604ab67-5104-4aa0-accc-28c2e64b003d"},
		"client_secret": {"fpvbADS606?$?)vdrHRKZ90"},
		"code":          {code},
		"grant_type":    {"authorization_code"},
		"redirect_uri":  {"http://localhost:8000/callback"},
	})
	res, err := ioutil.ReadAll(resp.Body)
	if err != nil || resp.StatusCode != 200 {
		fmt.Println("Onedrive: get Token Error")

		os.Exit(1)
	}
	defer resp.Body.Close()
	refresh_token := fastjson.GetBytes(res, "refresh_token")
	_ = ioutil.WriteFile("refreshtoken", refresh_token, os.ModeType)
	return refresh_token
}

func mkUploadRequest(filename string) *url.URL {
	resp := odPOST(fmt.Sprintf("https://graph.microsoft.com/v1.0/me/drive/root:%s:/createUploadSession", filename),
		"{\"item\":{\"@name.conflictBehavior\":\"replace\"}}")
	res, _ := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	if resp.StatusCode == 200 {
		Logln(fastjson.GetString(res, "uploadUrl"))
		u, _ := url.Parse(fastjson.GetString(res, "uploadUrl"))
		return u
	} else {
		Logln("Onedrive: Create Upload Session Error")
		Logln("Onedrive: %d %s", resp.StatusCode, res)
		return nil
	}
}

func odGET(subURL string) *http.Response {
	uri, _ := url.Parse(fmt.Sprintf("https://gitaccuacnz2-my.sharepoint.com%s", subURL))
	req := http.Request{
		Method: http.MethodGet,
		URL:    uri,
		Header: map[string][]string{
			"Authorization": {access_token},
			"Content-Type":  {"application/json"},
		},
	}
	resp, err := client.Do(&req)
	if err != nil {
		print("http Error")
		print(err)
		os.Exit(1)
	}
	return resp
}

func odPOST(URL string, data string) *http.Response {
	uri, _ := url.Parse(URL)
	req := http.Request{
		Method: http.MethodPost,
		URL:    uri,
		Header: map[string][]string{
			"Authorization": {access_token},
			"Content-Type":  {"application/json"},
		},
		Body: ioutil.NopCloser(strings.NewReader(data)),
	}
	resp, err := client.Do(&req)
	if err != nil {
		print("http Error")
		print(err)
		os.Exit(1)
	}
	return resp
}
