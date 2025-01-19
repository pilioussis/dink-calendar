package main

import (
	"fmt"
	"html/template"
	"image"
	"image/draw"
	"image/png"
	"log"
	"os"
	"os/exec"
	"time"

	calendar "google.golang.org/api/calendar/v3"
)

const PROJ_PATH = "/code"

const IN_HTML_TEMPLATE = "src/cal.template.html"
const OUT_HTML = "out/cal.html"

const FULL_PNG = "out/cal.png"
const DITHER_PNG = "out/dither.png"

const NUM_WEEKS = 30
const TZ = "Australia/Melbourne"

const EXPORT_WIDTH, EXPORT_HEIGHT = 1600, 1200

type SameDayEvent struct {
	Event     *calendar.Event
	StartTime *time.Time
	Class     string
}

type MultiDayEvent struct {
	Event     *calendar.Event
	StartDate time.Time
	EndDate   time.Time
	Position  int
	Class     string
}

type Day struct {
	Date               time.Time
	MonthBoundaryRight bool
	MonthBoundaryTop   bool
	IsToday            bool
	SameDayEvents      []*SameDayEvent
	MultiDayEvents     []*MultiDayEvent
	MultiDayMax        int
	MonthLabel         string
	Holiday            string // Is "" if not holiday
}

type Week struct {
	Days    []Day
	IsFirst bool
	IsLast  bool
}

type Calendar struct {
	Weeks    []Week
	Timezone string
}

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
	SameDay     []*SameDayEvent
	MultiDay    []*MultiDayEvent
	MultiDayMax int
	Holiday     string
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

func generateCalendar(start, now time.Time, dayEventsMap map[string]*DayEvents) Calendar {
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
				MultiDayMax:    dayEvents.MultiDayMax,
				MonthLabel:     monthLabel,
				Holiday:        dayEvents.Holiday,
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

	return Calendar{Weeks: weeks, Timezone: TZ}
}

func isSameDate(t1, t2 time.Time) bool {
	y1, m1, d1 := t1.Date()
	y2, m2, d2 := t2.Date()
	return y1 == y2 && m1 == m2 && d1 == d2
}

func CreateCalendarHTML(start, now time.Time, dayEventsMap map[string]*DayEvents) {
	c := generateCalendar(start, now, dayEventsMap)
	b, err := os.ReadFile(IN_HTML_TEMPLATE)
	if err != nil {
		log.Panicf("Error creating HTML: %v", err)
	}

	tmpl := string(b)

	t := template.Must(template.New("cal").Funcs(template.FuncMap{
		"isSameDate": isSameDate,
	}).Parse(tmpl))

	f, err := os.Create(OUT_HTML)
	if err != nil {
		log.Panicf("Error creating HTML: %v", err)
	}

	defer f.Close()

	err = t.Execute(f, c)
	fmt.Println("Created html")

	if err != nil {
		log.Panicf("Error creating HTML: %v", err)
	}
}

func trimScreenshot() {
	// Open the original PNG
	inFile, err := os.Open(FULL_PNG)
	if err != nil {
		panic(err)
	}
	defer inFile.Close()

	// Decode into an Image
	src, err := png.Decode(inFile)
	if err != nil {
		panic(err)
	}

	// Desired width/height
	// (Modify these to whatever you need)
	targetWidth := EXPORT_WIDTH
	targetHeight := EXPORT_HEIGHT

	// Create a new image with the desired cropped size
	rect := image.Rect(0, 0, targetWidth, targetHeight)
	dst := image.NewRGBA(rect)

	// Draw just the top-left portion onto dst
	draw.Draw(dst, rect, src, image.Point{0, 0}, draw.Src)

	// Write the new image to file
	outFile, err := os.Create(FULL_PNG)
	if err != nil {
		panic(err)
	}
	defer outFile.Close()

	if err = png.Encode(outFile, dst); err != nil {
		panic(err)
	}
	fmt.Println("Trimmed png")
}

func getScreenshot() {
	const paddingBottom = 200

	cmd := exec.Command(
		"chromium",
		"--headless",
		"--no-sandbox",
		fmt.Sprintf("--window-size=%v,%v", EXPORT_WIDTH, EXPORT_HEIGHT+paddingBottom),
		"--force-device-scale-factor=1",
		"--virtual-time-budget=5000",
		fmt.Sprintf("--screenshot=%s/%s", PROJ_PATH, FULL_PNG),
		fmt.Sprintf("file://%s/%s", PROJ_PATH, OUT_HTML),
	)
	out, err := cmd.Output()

	if err != nil {
		fmt.Println("Error creating image from HTML:", err, out)
		return
	}
	fmt.Println("Created png")

	trimScreenshot()
}

func main() {
	now := time.Now().AddDate(0, 0, 0)
	offset := (int(now.Weekday()) + 6) % 7
	start := now.AddDate(0, 0, -offset)

	var dayEventsMap map[string]*DayEvents
	if skip := true; skip {
		dayEventsMap = getCachedData()
	} else {
		dayEventsMap = getData(start, start.AddDate(0, 0, NUM_WEEKS*7))
	}

	CreateCalendarHTML(start, now, dayEventsMap)
	getScreenshot()
	Dither(FULL_PNG, DITHER_PNG)
}
