package main

import (
	"sync"
	"time"
)

//8:30-9:15
//9:20-10:05
//10:25-11:10
//11:15-12:00
//summer:
//14:30-15:15
//15:20-16:05
//16:25-17.10
//17:15-18:00

//19:00-19:45
//19:50-20:35
var timeTable = map[int]time.Duration{
	0: 0,

	1: 8*time.Hour + 30*time.Minute,
	2: 9*time.Hour + 20*time.Minute,
	3: 10*time.Hour + 25*time.Minute,
	4: 11*time.Hour + 15*time.Minute,

	5: 14*time.Hour + 30*time.Minute,
	6: 15*time.Hour + 20*time.Minute,
	7: 16*time.Hour + 25*time.Minute,
	8: 17*time.Hour + 15*time.Minute,

	9:  19 * time.Hour,
	10: 19*time.Hour + 50*time.Minute,
}

func Generator(name string, cs *Courses, start time.Time) string {
	c := NewCalendar(name, cs.generateEvents(start))
	return c.String()
}
func fixTime(t time.Time) time.Time {
	//before 5.1 and after 10.1
	if t.After(time.Date(t.Year(), 10, 1, 0, 0, 0, 0, time.Local)) ||
		t.Before(time.Date(t.Year(), 5, 1, 0, 0, 0, 0, time.Local)) {
		return t.Add(-45 * time.Minute)
	}
	return t
}
func (cs Courses) generateEvents(start time.Time) Events {
	var (
		ret  Events
		wg   sync.WaitGroup
		lock sync.Mutex
	)
	for _, course := range cs {
		for _, item := range course.Data {
			wg.Add(1)
			go func(d *Detail, c *Course) {
				lock.Lock()
				e := generateFromDetail(d, start, c.Name, c.Teacher)
				ret = append(ret, e...)
				lock.Unlock()
				wg.Done()
			}(item, course)
		}
	}
	wg.Wait()
	return ret
}
func generateFromDetail(d *Detail, start time.Time, name string, teacher string) Events {
	var ret Events
	for i, v := range d.Week {
		var (
			b time.Time
			e time.Time
		)
		startDay := start.AddDate(0, 0, 7*(v-1)+d.IndexDay-1)
		_b := startDay.Add(timeTable[d.Begin])
		_e := startDay.Add(timeTable[d.End] + 45*time.Minute)
		if d.Begin >= 5 {
			b = fixTime(_b)
		} else {
			b = _b
		}
		if d.End >= 5 {
			e = fixTime(_e)
		} else {
			e = _e
		}
		ret = append(ret, NewEvent(
			DateToStamp(b),
			DateToStamp(e),
			teacher,
			d.Place,
			i,
			name,
			"CONFIRMED",
			"PUBLIC",
		))
	}
	return ret
}
