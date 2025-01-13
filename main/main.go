package main

import (
	"fmt"
	"html/template"
	"log"
	"os"
	"os/exec"
	"time"

	calendar "google.golang.org/api/calendar/v3"
)

type MultiDayEvent struct {
	Event     *calendar.Event
	StartDate time.Time
	EndDate   time.Time
	Position  int
}

type Day struct {
	Date               time.Time
	MonthBoundaryRight bool
	MonthBoundaryTop   bool
	IsToday            bool
	SameDayEvents      []*calendar.Event
	MultiDayEvents     []*MultiDayEvent
	MonthLabel         string
}

type Week struct {
	Days    []Day
	IsFirst bool
	IsLast  bool
}

type Calendar struct {
	Weeks []Week
}

const NUM_WEEKS = 52

const TZ = "Australia/Melbourne"

func getTimeFromString(dateTimeString string, dateString string) time.Time {
	loc, err := time.LoadLocation(TZ)
	if err != nil {
		log.Panicf("Error loading location: %s", TZ)
	}
	if dateTimeString != "" {
		t, err := time.ParseInLocation(time.RFC3339, dateTimeString, loc)
		if err != nil {
			log.Panicf("Error loading dateTimeString: %+v", dateTimeString)
		}
		return t
	}

	if dateString != "" {
		t, err := time.ParseInLocation(time.DateOnly, dateString, loc)
		if err != nil {
			log.Panicf("Error loading dateString: %+v", dateString)
		}
		return t
	}
	panic("Both date strings empty")
}

func getMapKey(t time.Time) string {
	return t.Format("2006-01-02")
}

type DayEvents struct {
	SameDay  []*calendar.Event
	MultiDay []*MultiDayEvent
}

func getEventsForDays(events *calendar.Events) map[string]*DayEvents {
	days := make(map[string]*DayEvents)

	for _, e := range events.Items {
		start := getTimeFromString(e.Start.DateTime, e.Start.Date)
		end := getTimeFromString(e.End.DateTime, e.End.Date)

		fmt.Printf("%#v -- %#v %s %s %#v\n", start, end, e.End.Date, e.End.DateTime, e.Summary)

		y1, m1, d1 := start.Date()
		y2, m2, d2 := end.Date()

		sameDay := (y1 == y2 && m1 == m2 && d1 == d2)

		if sameDay {
			v, ok := days[getMapKey(start)]
			if !ok {
				v = &DayEvents{}
				days[getMapKey(start)] = v
			}
			v.SameDay = append(v.SameDay, e)
		} else {
			end = end.Add(-time.Duration(24) * time.Hour) // google does exclusive days
			for current := start; !current.After(end); current = current.AddDate(0, 0, 1) {
				v, ok := days[getMapKey(current)]
				if !ok {
					v = &DayEvents{}
					days[getMapKey(current)] = v
				}
				v.MultiDay = append(v.MultiDay, &MultiDayEvent{
					Event:     e,
					StartDate: start,
					EndDate:   end,
					Position:  0,
				})
			}
		}
	}

	return days
}

func generateCalendar() Calendar {
	now := time.Now().AddDate(0, 0, 0)
	offset := (int(now.Weekday()) + 6) % 7
	start := now.AddDate(0, 0, -offset)

	fmt.Println("Getting data")
	events := getData(start, start.AddDate(0, 0, NUM_WEEKS*7))

	dayEventsMap := getEventsForDays(events)
	fmt.Printf("Got %d events\n", len(events.Items))

	var weeks []Week
	currDay := start

	for w := 0; w < NUM_WEEKS; w++ {
		var days []Day
		for d := 0; d < 7; d++ {

			monthLabel := ""
			if currDay.Day() == 1 {
				monthLabel = currDay.Format("Jan")
			}

			dayEvents, ok := dayEventsMap[getMapKey(currDay)]

			if !ok {
				dayEvents = &DayEvents{}
			}

			days = append(days, Day{
				Date:           currDay,
				IsToday:        currDay == now,
				SameDayEvents:  dayEvents.SameDay,
				MultiDayEvents: dayEvents.MultiDay,
				MonthLabel:     monthLabel,
			})
			currDay = currDay.AddDate(0, 0, 1)
		}
		weeks = append(weeks, Week{
			Days:    days,
			IsFirst: w == 0,
			IsLast:  w == NUM_WEEKS-1,
		})
	}

	// Mark month boundaries
	for w := 0; w < len(weeks); w++ {
		for d := 0; d < len(weeks[w].Days); d++ {
			// Check horizontal boundary
			if d < 6 {
				if weeks[w].Days[d].Date.Month() != weeks[w].Days[d+1].Date.Month() {
					weeks[w].Days[d].MonthBoundaryRight = true
				}
			}
			// Check vertical boundary
			if w > 0 {
				if weeks[w].Days[d].Date.Month() != weeks[w-1].Days[d].Date.Month() {
					weeks[w].Days[d].MonthBoundaryTop = true
				}
			}
		}
	}

	return Calendar{Weeks: weeks}
}

func isSameDate(t1, t2 time.Time) bool {
	fmt.Printf("%#v %#v\n", t1, t2)
	y1, m1, d1 := t1.Date()
	y2, m2, d2 := t2.Date()
	return y1 == y2 && m1 == m2 && d1 == d2
}

// CreateCalendarHTML returns the rendered HTML for the 8-week calendar.
func CreateCalendarHTML() error {
	c := generateCalendar()
	b, err := os.ReadFile("src/cal.template.html")
	if err != nil {
		return err
	}

	tmpl := string(b)

	t := template.Must(template.New("cal").Funcs(template.FuncMap{
		"isSameDate": isSameDate,
	}).Parse(tmpl))

	f, err := os.Create("/code/out/cal.html")
	fmt.Println("Created HTML")
	if err != nil {
		return err
	}

	defer f.Close()

	err = t.Execute(f, c)
	return err
}

func main() {
	err := CreateCalendarHTML()
	if err != nil {
		fmt.Print(err)
		return
	}

	cmd := exec.Command(
		"wkhtmltoimage",
		"--enable-local-file-access",
		"--zoom", "10",
		"--width", "800",
		"file:///code/out/cal.html",
		"out/cal.png",
	)
	out, err := cmd.Output()

	if err != nil {
		fmt.Println("Error creating image from HTML:", err, out)
		return
	}
	fmt.Println("Created png")
}
