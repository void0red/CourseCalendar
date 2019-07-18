package main

import (
	"encoding/json"
	"flag"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

var (
	courseAppid = "4770397878132218"
	termUrl     = "http://ehall.xidian.edu.cn/jwapp/sys/wdkb/modules/jshkcb/dqxnxq.do"
	serviceUrl  = "http://ehall.xidian.edu.cn/jwapp/sys/wdkb/modules/xskcb/xskcb.do"
	username    string
	password    string
	startDate   string
	icsName     string
)

func init() {
	flag.StringVar(&username, "u", "", "username")
	flag.StringVar(&password, "p", "", "password")
	flag.StringVar(&startDate, "d", "", "start date(default: today)\n for example:2019.7.18")
}
func main() {
	flag.Parse()
	start := handleDate(startDate)
	if username == "" || password == "" {
		flag.Usage()
		return
	}
	s := NewSession().Login(username, password).Service(courseAppid)
	term := getTerm(s)
	icsName = username + " - " + term
	courses := getCourse(s, term)
	c := Parser(courses)
	data := Generator(icsName, &c, start)
	_ = ioutil.WriteFile(icsName+".ics", []byte(data), 0775)
}
func getTerm(s *Session) string {
	c := (*http.Client)(s)

	req, _ := http.NewRequest("POST", termUrl, nil)
	resp, _ := c.Do(req)
	body, _ := ioutil.ReadAll(resp.Body)

	var t map[string]interface{}
	_ = json.Unmarshal(body, &t)
	term := t["datas"].(map[string]interface{})["dqxnxq"].(map[string]interface{})["rows"].([]interface{})[0].(map[string]interface{})["DM"].(string)

	return term
}
func getCourse(s *Session, term string) io.Reader {
	c := (*http.Client)(s)

	req, _ := http.NewRequest("POST", serviceUrl, strings.NewReader(url.Values{
		"XNXQDM": []string{term},
	}.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	resp, _ := c.Do(req)

	return resp.Body
}
func handleDate(t string) time.Time {
	if t == "" {
		return time.Now()
	} else {
		s := strings.Split(t, ".")
		if len(s) < 3 {
			return time.Now()
		}
		year, _ := strconv.Atoi(s[0])
		month, _ := strconv.Atoi(s[1])
		day, _ := strconv.Atoi(s[2])
		return time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.Local)
	}
}
