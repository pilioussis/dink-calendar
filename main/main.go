package main

import (
	"fmt"
	"html/template"
	"os"
	"os/exec"
	"slices"
	"time"

	calendar "google.golang.org/api/calendar/v3"
)

type Day struct {
	Date               time.Time
	MonthBoundaryRight bool
	MonthBoundaryTop   bool
	IsToday            bool
	Events             []*calendar.Event
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

func sameDayRFC3339(a time.Time, b, tz string) (bool, error) {
	if b == "" {
		return false, nil
	}
	loc, err := time.LoadLocation(tz)
	if err != nil {
		return false, err
	}

	t1 := a
	t2, err := time.Parse(time.RFC3339, b)

	if err != nil {
		return false, err
	}

	y1, m1, d1 := t1.In(loc).Date()
	y2, m2, d2 := t2.In(loc).Date()

	return (y1 == y2 && m1 == m2 && d1 == d2), nil
}

func generateCalendar() Calendar {
	now := time.Now().AddDate(0, 0, 0)
	offset := (int(now.Weekday()) + 6) % 7
	start := now.AddDate(0, 0, -offset)

	fmt.Println("Getting data")
	events := getData(start, start.AddDate(0, 0, NUM_WEEKS*7))
	fmt.Printf("Got %d events\n", len(events.Items))

	// for _, i := range events.Items {
	// 	start := i.Start.DateTime
	// 	if start == "" {
	// 		start = i.Start.Date
	// 	}
	// 	fmt.Printf("%s - %s\n", start, i.Summary)
	// }

	var weeks []Week
	day := start

	for w := 0; w < NUM_WEEKS; w++ {
		var days []Day
		for d := 0; d < 7; d++ {
			dayEvents := slices.Collect(
				func(yield func(*calendar.Event) bool) {
					for _, v := range events.Items {
						d := v.Start.DateTime
						if d == "" {
							if v.Start.Date != "" {
								d = v.Start.Date + "T00:00:00-08:00"
							} else {
								fmt.Printf("%+v ||| %+v", v, v.Start)
								os.Exit(1)
							}
						}
						sameDay, err := sameDayRFC3339(day, d, "Australia/Melbourne")
						if err != nil {
							fmt.Println("Error equating dates", err)
						}
						if sameDay {
							if !yield(v) {
								return // triggered in "break"
							}
						}
					}
				},
			)
			monthLabel := ""
			if day.Day() == 1 {
				monthLabel = day.Format("Jan")
			}
			days = append(days, Day{
				Date:       day,
				IsToday:    day == now,
				Events:     dayEvents,
				MonthLabel: monthLabel,
			})
			day = day.AddDate(0, 0, 1)
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

// CreateCalendarHTML returns the rendered HTML for the 8-week calendar.
func CreateCalendarHTML() error {
	c := generateCalendar()
	b, err := os.ReadFile("src/cal.template.html")
	if err != nil {
		return err
	}

	tmpl := string(b)

	t := template.Must(template.New("cal").Parse(tmpl))

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
