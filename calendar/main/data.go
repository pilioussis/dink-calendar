package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"sort"
	"strings"
	"time"

	calendar "google.golang.org/api/calendar/v3"
)

type CalendarSource struct {
	Name  string
	Token string
	ID    string
}

const DEAN_EMAIL = "dean.pilioussis@gmail.com"
const STRUGS_EMAIL = "mcpherson.sarah.a@gmail.com"

var CALENDAR_DEAN = CalendarSource{
	Name:  "dean-personal",
	Token: DEAN_TOKEN,
	ID:    "primary",
}

var CALENDAR_STRUGS = CalendarSource{
	Name:  "strugs-personal",
	Token: DEAN_TOKEN,
	ID:    "mcpherson.sarah.a@gmail.com",
}

var CALENDAR_BIRTHDAYS = CalendarSource{
	Name:  "birthdays",
	Token: DEAN_TOKEN,
	ID:    "403994ecc2585854c8e932c00d1ca82c7cb9b423fdab94e0b5b6be2c56335b9d@group.calendar.google.com",
}

var CALENDAR_HOLIDAYS = CalendarSource{
	Name:  "holidays",
	Token: DEAN_TOKEN,
	ID:    "ee92ea54f1e2e5fab5aee1a88873031d57d2ea0b164a6968a4a943c9121bf292@group.calendar.google.com",
}

const CACHE_FOLDER = "out/cache"

var TZ = mustLoadLocation("Australia/Melbourne")

func getCacheFile(cal CalendarSource) string {
	return fmt.Sprintf("%s/%s.json", CACHE_FOLDER, cal.Name)
}

func getDataForCal(start, end time.Time, cal CalendarSource, useCache bool) *calendar.Events {
	if useCache {
		return getCachedData(cal)
	}
	client := getToken(cal.Token)
	srv, err := calendar.New(client)
	if err != nil {
		log.Fatalf("Unable to retrieve Calendar client: %v", err)
	}

	fmt.Printf("Get events: %s", cal.Name)
	events, err := srv.Events.List(cal.ID).
		ShowDeleted(false).
		SingleEvents(true).
		TimeMin(start.UTC().Format(time.RFC3339)).
		TimeMax(end.UTC().Format(time.RFC3339)).
		OrderBy("startTime").
		Do()
	if err != nil {
		log.Fatalf("\nUnable to retrieve events: %v", err)
	}

	os.MkdirAll(CACHE_FOLDER, os.ModePerm)
	file, err := os.Create(getCacheFile(cal))
	if err != nil {
		log.Fatal("\nSaving cache error:", err)
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(events); err != nil {
		log.Fatal(err)
	}

	fmt.Println("\rGot events:", len(events.Items), cal.Name)
	return events
}

func getEventsForDays(dayEventsMap map[string]*DayEvents, events []*calendar.Event, calname string) {

	lowestFreeSpot := make(map[int]time.Time)

	for _, e := range events {
		start := getTimeFromString(e.Start.DateTime, e.Start.Date)
		end := getTimeFromString(e.End.DateTime, e.End.Date)

		if e.End.Date != "" {
			// google does exclusive days
			end = end.Add(-time.Duration(24) * time.Hour)
		}

		y1, m1, d1 := start.Date()
		y2, m2, d2 := end.Date()
		sameDay := (y1 == y2 && m1 == m2 && d1 == d2)

		if sameDay {
			v, ok := dayEventsMap[getMapKey(start)]
			if !ok {
				v = &DayEvents{}
				dayEventsMap[getMapKey(start)] = v
			}
			startTime := &start
			if e.Start.DateTime == "" {
				startTime = nil
			}
			v.SameDay = append(v.SameDay, &SameDayEvent{Event: e, StartTime: startTime, Class: calname})
		} else {
			for pos, t := range lowestFreeSpot {
				if start.After(t) {
					delete(lowestFreeSpot, pos)
				}
			}
			lowest := -1
			for i := range 10 {
				if _, ok := lowestFreeSpot[i]; !ok {
					lowest = i
					break
				}
			}

			lowestFreeSpot[lowest] = end

			init := time.Date(start.Year(), start.Month(), start.Day(), 0, 0, 0, 0, TZ)
			for current := init; !current.After(end); current = current.AddDate(0, 0, 1) {
				v, ok := dayEventsMap[getMapKey(current)]
				if !ok {
					v = &DayEvents{}
					dayEventsMap[getMapKey(current)] = v
				}
				if len(e.Summary) < 25 {
					e.Summary = e.Summary + strings.Repeat("-", max(16-len(e.Summary), 0))
				}

				var startTime *time.Time
				if e.Start.DateTime != "" {
					// Only set time if event has a time
					startTime = &start
				}
				v.MultiDay = append(v.MultiDay, &MultiDayEvent{
					Event:     e,
					StartDate: start,
					EndDate:   end,
					StartTime: startTime,
					Position:  lowest,
					Class:     calname,
				})
				v.MultiDayMax = max(v.MultiDayMax, lowest)

				sort.Slice(v.MultiDay, func(i, j int) bool {
					return v.MultiDay[i].Position < v.MultiDay[j].Position
				})
			}
		}
	}
}

func filterShared(events []*calendar.Event, email string) ([]*calendar.Event, []*calendar.Event) {
	individual := []*calendar.Event{}
	shared := []*calendar.Event{}
	for _, e := range events {
		skip := false
		for _, a := range e.Attendees {
			if a.Email == email {
				skip = true
				shared = append(shared, e)
			}
		}
		if !skip {
			individual = append(individual, e)
		}
	}
	return individual, shared
}

func getCachedData(cal CalendarSource) *calendar.Events {
	fmt.Println("Got cached events:", cal.Name)
	cache_str, err := os.ReadFile(getCacheFile(cal))
	var events *calendar.Events

	if err != nil {
		log.Panicln("Cache file not found", cal.ID, err)
	}

	err = json.Unmarshal([]byte(cache_str), &events)
	if err != nil {
		log.Panicln("Error loading cache", err)
	}

	return events
}

func getData(start, end time.Time, createStubEvents, useCache bool) map[string]*DayEvents {
	dayEventsMap := make(map[string]*DayEvents)
	dean_pre := getDataForCal(start, end, CALENDAR_DEAN, useCache)
	strugs_pre := getDataForCal(start, end, CALENDAR_STRUGS, useCache)
	birthdays := getDataForCal(start, end, CALENDAR_BIRTHDAYS, useCache)
	holidays := getDataForCal(start, end, CALENDAR_HOLIDAYS, useCache)

	dean, shared := filterShared(dean_pre.Items, STRUGS_EMAIL)
	strugs, _ := filterShared(strugs_pre.Items, DEAN_EMAIL)

	if createStubEvents {
		times := [][2]string{
			{"2025-02-12T13:15:00+11:00", "2025-02-16"},
		}
		stub_str, err := os.ReadFile("src/stub.json")
		if err != nil {
			log.Panicln("Event stub file not found:", err)
		}

		for _, t := range times {
			var fakeEvent *calendar.Event
			err := json.Unmarshal([]byte(stub_str), &fakeEvent)
			if err != nil {
				log.Panicln("Unable to create fake events", err)
			}
			if strings.Contains(t[0], ":") {
				fakeEvent.Start.DateTime = t[0]
			} else {
				fakeEvent.Start.Date = t[0]
			}
			if strings.Contains(t[1], ":") {
				fakeEvent.End.DateTime = t[1]
			} else {
				fakeEvent.End.Date = t[1]
			}
			dean = append(dean, fakeEvent)
		}

	}

	getEventsForDays(dayEventsMap, shared, "e-shared")
	getEventsForDays(dayEventsMap, dean, "e-dean")
	getEventsForDays(dayEventsMap, strugs, "e-strugs")
	getEventsForDays(dayEventsMap, birthdays.Items, "e-birthday")

	for _, e := range holidays.Items {
		kTime := getTimeFromString(e.Start.DateTime, e.Start.Date)
		v, ok := dayEventsMap[getMapKey(kTime)]
		if !ok {
			v = &DayEvents{}
			dayEventsMap[getMapKey(kTime)] = v

		}
		v.Holiday = e.Summary
	}

	return dayEventsMap
}
