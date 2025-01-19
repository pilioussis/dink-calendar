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

const DEAN_EMAIL = "dean.pilioussis@gmail.com"
const STRUGS_EMAIL = "mcpherson.sarah.a@gmail.com"
const CALENDAR_BIRTHDAYS = "403994ecc2585854c8e932c00d1ca82c7cb9b423fdab94e0b5b6be2c56335b9d@group.calendar.google.com"
const CALENDAR_PERSONAL = "primary"

const CALENDAR_HOLIDAYS = "ee92ea54f1e2e5fab5aee1a88873031d57d2ea0b164a6968a4a943c9121bf292@group.calendar.google.com"

const FILE_CACHE = "out/cached.json"

func getDataForCal(start, end time.Time, tokenFile, calId string) *calendar.Events {
	client := getttt(tokenFile)
	srv, err := calendar.New(client)
	if err != nil {
		log.Fatalf("Unable to retrieve Calendar client: %v", err)
	}
	fmt.Println("Getting events")
	events, err := srv.Events.List(calId).
		ShowDeleted(false).
		SingleEvents(true).
		TimeMin(start.UTC().Format(time.RFC3339)).
		TimeMax(end.UTC().Format(time.RFC3339)).
		OrderBy("startTime").
		Do()
	if err != nil {
		log.Fatalf("Unable to retrieve events: %v", err)
	}

	fmt.Println("Got events", len(events.Items))
	return events
}

func getEventsForDays(dayEventsMap map[string]*DayEvents, events []*calendar.Event, calname string) {

	lowestFreeSpot := make(map[int]time.Time)

	for _, e := range events {
		start := getTimeFromString(e.Start.DateTime, e.Start.Date)
		end := getTimeFromString(e.End.DateTime, e.End.Date)

		y1, m1, d1 := start.Date()
		y2, m2, d2 := end.Add(time.Microsecond * -1).Date() // google does exclusive days

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
			end = end.Add(-time.Duration(24) * time.Hour) // google does exclusive days

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

			for current := start; !current.After(end); current = current.AddDate(0, 0, 1) {
				v, ok := dayEventsMap[getMapKey(current)]
				if !ok {
					v = &DayEvents{}
					dayEventsMap[getMapKey(current)] = v
				}
				if len(e.Summary) < 25 {
					e.Summary = e.Summary + strings.Repeat("-", 25-len(e.Summary))
				}
				v.MultiDay = append(v.MultiDay, &MultiDayEvent{
					Event:     e,
					StartDate: start,
					EndDate:   end,
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

func getCachedData() map[string]*DayEvents {
	fmt.Println("Using cached events")
	cache_str, err := os.ReadFile(FILE_CACHE)
	dayEventsMap := make(map[string]*DayEvents)

	if err != nil {
		log.Panicln("Cache file not found", FILE_CACHE, err)
	}

	err = json.Unmarshal([]byte(cache_str), &dayEventsMap)
	if err != nil {
		log.Panicln("Error loading cache", err)
	}
	return dayEventsMap
}

func getData(start, end time.Time) map[string]*DayEvents {
	dayEventsMap := make(map[string]*DayEvents)
	dean_pre := getDataForCal(start, end, DEAN_TOKEN, CALENDAR_PERSONAL)
	strugs_pre := getDataForCal(start, end, STRUGS_TOKEN, CALENDAR_PERSONAL)
	birthdays := getDataForCal(start, end, DEAN_TOKEN, CALENDAR_BIRTHDAYS)
	holidays := getDataForCal(start, end, DEAN_TOKEN, CALENDAR_HOLIDAYS)

	dean, shared := filterShared(dean_pre.Items, STRUGS_EMAIL)
	strugs, _ := filterShared(strugs_pre.Items, DEAN_EMAIL)

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

	file, err := os.Create(FILE_CACHE)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(dayEventsMap); err != nil {
		log.Fatal(err)
	}

	return dayEventsMap
}
