package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"regexp"
	"strings"
)

type Session http.Client

const (
	chromeAgent = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/75.0.3770.100 Safari/537.36"
	loginUrl    = "http://ids.xidian.edu.cn/authserver/login"
	appUrl      = "http://ehall.xidian.edu.cn//appShow"
)

func NewSession() *Session {
	jar, _ := cookiejar.New(nil)
	return (*Session)(&http.Client{
		Jar: jar,
	})
}

func (s *Session) Login(username string, password string) *Session {
	client := (*http.Client)(s)

	req, _ := http.NewRequest("GET", loginUrl, nil)
	req.Header.Set("User-Agent", chromeAgent)
	resp, _ := client.Do(req)
	body, _ := ioutil.ReadAll(resp.Body)
	ltPattern := regexp.MustCompile("<input type=\"hidden\" name=\"lt\" value=\"(.*)\" />")
	lt := ltPattern.FindSubmatch(body)[0][38:100]
	executionPattern := regexp.MustCompile("<input type=\"hidden\" name=\"execution\" value=\"(.*)\" />")
	execution := executionPattern.FindSubmatch(body)[0][45:49]

	postForm := url.Values{}
	postForm.Set("username", username)
	postForm.Set("password", password)
	postForm.Set("lt", string(lt))
	postForm.Set("execution", string(execution))
	postForm.Set("_eventId", "submit")
	postForm.Set("rmShown", "1")

	req, _ = http.NewRequest("POST", loginUrl, strings.NewReader(postForm.Encode()))
	req.Header.Set("User-Agent", chromeAgent)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	_, _ = client.Do(req)

	req, _ = http.NewRequest("GET", appUrl, nil)
	req.Header.Set("User-Agent", chromeAgent)
	_, _ = client.Do(req)

	return s
}

func (s *Session) Service(appid string) *Session {
	client := (*http.Client)(s)

	serviceUrl := fmt.Sprintf("%s?appId=%s", appUrl, appid)
	req, _ := http.NewRequest("GET", serviceUrl, nil)
	req.Header.Set("User-Agent", chromeAgent)
	_, _ = client.Do(req)

	return s
}
