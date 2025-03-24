package main

import (
	"log"
	"time"
)

func mustLoadLocation(name string) *time.Location {
	loc, err := time.LoadLocation(name)
	if err != nil {
		log.Panicf("Failed to load timezone: %s", name)
	}
	return loc
}

func getTimeFromString(dateTimeString string, dateString string) time.Time {
	var t time.Time
	if dateTimeString != "" {
		parsed, err := time.ParseInLocation(time.RFC3339, dateTimeString, TZ)
		if err != nil {
			log.Panicf("Error loading dateTimeString: %+v", dateTimeString)
		}
		t = parsed
	} else if dateString != "" {
		parsed, err := time.ParseInLocation(time.DateOnly, dateString, TZ)
		if err != nil {
			log.Panicf("Error loading dateString: %+v", dateString)
		}
		t = parsed
	} else {
		panic("Both date strings empty")
	}

	return t.In(TZ)
}

func isSameDate(t1, t2 time.Time) bool {
	y1, m1, d1 := t1.Date()
	y2, m2, d2 := t2.Date()
	return y1 == y2 && m1 == m2 && d1 == d2
}

func isLessThanDate(t1, t2 time.Time) bool {
	y1, m1, d1 := t1.Date()
	y2, m2, d2 := t2.Date()
	return y1 < y2 || (y1 == y2 && (m1 < m2 || (m1 == m2 && d1 < d2)))
}

func isBeforeFirstDay(currDay, calStart, eventStart time.Time) bool {
	startedBeforeCalendar := isLessThanDate(eventStart, calStart)
	isFirstDayOfCalendar := isSameDate(currDay, calStart)
	return startedBeforeCalendar && isFirstDayOfCalendar
}
