package main

import (
	"bytes"
	"fmt"
	uuid "github.com/satori/go.uuid"
	"strconv"
	"time"
)

type Event struct {
	start       string
	end         string
	created     string
	modified    string
	stamp       string
	uid         string
	class       string
	description string
	location    string
	sequence    int
	status      string
	summary     string
}

type Events []*Event

func NewEvent(start string, end string, description string, location string, sequence int,
	summary string, status string, class string) *Event {
	stamp := DateToStamp(time.Now())
	uid, _ := uuid.NewV4()
	return &Event{start: start, end: end, created: stamp, modified: stamp, stamp: stamp,
		status: status, class: class, uid: uid.String(),
		description: description, location: location, sequence: sequence, summary: summary}
}

func DateToStamp(t time.Time) string {
	p := t.UTC()
	return fmt.Sprintf("%04d%02d%02dT%02d%02d%02dZ", p.Year(), p.Month(), p.Day(), p.Hour(), p.Minute(), p.Second())
}

func (e *Event) String() string {
	var ret bytes.Buffer
	ret.WriteString("BEGIN:VEVENT\n")
	ret.WriteString("DTSTART:" + e.start + "\n")
	ret.WriteString("DTEND:" + e.end + "\n")
	ret.WriteString("CREATED:" + e.created + "\n")
	ret.WriteString("LAST-MODIFIED:" + e.modified + "\n")
	ret.WriteString("DTSTAMP:" + e.stamp + "\n")
	ret.WriteString("UID:" + e.uid + "\n")
	ret.WriteString("CLASS:" + e.class + "\n")
	ret.WriteString("DESCRIPTION:" + e.description + "\n")
	ret.WriteString("LOCATION:" + e.location + "\n")
	ret.WriteString("SEQUENCE:" + strconv.Itoa(e.sequence) + "\n")
	ret.WriteString("STATUS:" + e.status + "\n")
	ret.WriteString("SUMMARY:" + e.summary + "\n")
	ret.WriteString("END:VEVENT\n")
	return ret.String()
}

type Calendar struct {
	version  string
	prodid   string
	calscale string
	method   string
	calname  string
	timezone string
	event    Events
}

//
//func NewCalendar(version string, prodid string, calscale string, method string, calname string, timezone string, event Events) *Calendar {
//	return &Calendar{version: version, prodid: prodid, calscale: calscale, method: method, calname: calname, timezone: timezone, event: event}
//}

func NewCalendar(calname string, event Events) *Calendar {
	return &Calendar{version: "2.0", prodid: "-//Google Inc//Google Calendar 70.9054//EN",
		calscale: "GREGORIAN", method: "PUBLISH", timezone: "Asia/Shanghai",
		calname: calname, event: event}
}

func (c *Calendar) String() string {
	var ret bytes.Buffer
	ret.WriteString("BEGIN:VCALENDAR\n")
	ret.WriteString("PRODID:" + c.prodid + "\n")
	ret.WriteString("CALSCALE:" + c.calscale + "\n")
	ret.WriteString("METHOD:" + c.method + "\n")
	ret.WriteString("X-WR-CALNAME:" + c.calname + "\n")
	ret.WriteString("X-WR-TIMEZONE:" + c.timezone + "\n")
	for _, v := range c.event {
		ret.WriteString(v.String())
	}
	ret.WriteString("END:VCALENDAR\n")
	return ret.String()
}
