package main

import (
	"fmt"
	"log"
	"time"

	calendar "google.golang.org/api/calendar/v3"
)

const DEAN_EMAIL = "dean.pilioussis@gmail.com"
const STRUGS_EMAIL = "mcpherson.sarah.a@gmail.com"
const DEAN_TOKEN = "dean_token.json"
const STRUGS_TOKEN = "strugs_token.json"
const CALENDAR_BIRTHDAYS = "403994ecc2585854c8e932c00d1ca82c7cb9b423fdab94e0b5b6be2c56335b9d@group.calendar.google.com"
const CALENDAR_PERSONAL = "primary"

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

func getData(start, end time.Time) map[string]*DayEvents {
	dayEventsMap := make(map[string]*DayEvents)
	dean_pre := getDataForCal(start, end, DEAN_TOKEN, CALENDAR_PERSONAL)
	strugs_pre := getDataForCal(start, end, STRUGS_TOKEN, CALENDAR_PERSONAL)
	birthdays := getDataForCal(start, end, DEAN_TOKEN, CALENDAR_BIRTHDAYS)

	dean, shared := filterShared(dean_pre.Items, STRUGS_EMAIL)
	strugs, _ := filterShared(strugs_pre.Items, DEAN_EMAIL)

	getEventsForDays(dayEventsMap, shared, "e-shared")
	getEventsForDays(dayEventsMap, dean, "e-dean")
	getEventsForDays(dayEventsMap, strugs, "e-strugs")
	getEventsForDays(dayEventsMap, birthdays.Items, "e-birthday")
	return dayEventsMap
}
