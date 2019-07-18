package main

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"regexp"
	"strconv"
	"strings"
	"sync"
)

func Parser(r io.Reader) Courses {
	var (
		ret  Courses
		wg   sync.WaitGroup
		lock sync.Mutex
	)
	s, _ := ioutil.ReadAll(r)
	var data map[string]interface{}
	_ = json.Unmarshal(s, &data)
	d := data["datas"].(map[string]interface{})["xskcb"].(map[string]interface{})["rows"].([]interface{})
	for _, v := range d {
		t := v.(map[string]interface{})
		wg.Add(1)
		go func(s map[string]interface{}) {
			name := t["KCM"].(string)
			teacher := t["SKJS"].(string)
			var target *Course
			lock.Lock()
			b, i := ret.CheckRepeat(name, teacher)
			if b {
				target = ret[i]
			} else {
				target = NewCourse(name, teacher)
				ret = append(ret, target)
			}
			lock.Unlock()
			target.itemParser(s)
			wg.Done()
		}(t)
	}
	wg.Wait()
	return ret
}
func parseDetail(s map[string]interface{}) *Detail {
	var place string
	if s["JASMC"] != nil {
		place = s["JASMC"].(string)
	} else {
		place = s["JASDM"].(string)
	}
	return NewDetail(
		parseWeek(s["ZCMC"].(string)),
		parseIndexDay(s["SKXQ"].(string)),
		parseBegin(s["KSJC"].(string)),
		parseEnd(s["JSJC"].(string)),
		place,
	)
}
func parseWeek(s string) []int {
	var ret []int
	str := strings.Split(s, ",")
	p := regexp.MustCompile(`(.*?)周`)
	flag := strings.Contains(s, "单")
	for _, v := range str {
		l := p.FindStringSubmatch(v)[1]
		t := strings.Split(l, "-")
		if len(t) > 1 {
			begin, _ := strconv.Atoi(t[0])
			end, _ := strconv.Atoi(t[1])
			if flag {
				for i := begin; i <= end; i += 2 {
					ret = append(ret, i)
				}
			} else {
				for i := begin; i <= end; i++ {
					ret = append(ret, i)
				}
			}
		} else {
			r, _ := strconv.Atoi(l)
			ret = append(ret, r)
		}
	}
	return ret
}
func parseIndexDay(s string) int {
	r, _ := strconv.Atoi(s)
	return r
}
func parseBegin(s string) int {
	r, _ := strconv.Atoi(s)
	return r
}
func parseEnd(s string) int {
	r, _ := strconv.Atoi(s)
	return r
}
