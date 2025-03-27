package main

import (
	"fmt"
	"html/template"
	"image"
	"image/draw"
	"image/png"
	"log"
	"log/slog"
	"os"
	"os/exec"
	"time"

	calendar "google.golang.org/api/calendar/v3"
)

const IN_HTML_TEMPLATE = "src/cal.template.html"
const OUT_HTML = "out/cal.html"

const FULL_COLOR_PATH = "out/cal.png"
const DITHERED_PATH = "out/dither.bmp"

const NUM_WEEKS = 8

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
	StartTime *time.Time
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
	Weeks     []Week
	Start     time.Time
	CreatedAt string
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
				IsToday:        currDay.Format("2006-01-02") == now.Format("2006-01-02"),
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

	return Calendar{Weeks: weeks, Start: start, CreatedAt: time.Now().In(TZ).Format("2006-01-02 15:04")}
}

func CreateCalendarHTML(start, now time.Time, dayEventsMap map[string]*DayEvents) error {
	c := generateCalendar(start, now, dayEventsMap)
	b, err := os.ReadFile(IN_HTML_TEMPLATE)
	if err != nil {
		log.Panicf("Error creating HTML: %v", err)
	}

	t := template.Must(template.New("cal").Funcs(template.FuncMap{
		"isSameDate":       isSameDate,
		"isBeforeFirstDay": isBeforeFirstDay,
	}).Parse(string(b)))

	f, err := os.Create(OUT_HTML)
	if err != nil {
		return fmt.Errorf("error creating html file: %w", err)
	}
	defer f.Close()
	slog.Info("Created html")

	err = t.Execute(f, c)
	if err != nil {
		return fmt.Errorf("error rendering html template: %w", err)
	}
	return nil
}

func trimScreenshot() error {
	// Open the original PNG
	inFile, err := os.Open(FULL_COLOR_PATH)
	if err != nil {
		panic(err)
	}
	defer inFile.Close()

	// Decode into an Image
	src, err := png.Decode(inFile)
	if err != nil {
		return fmt.Errorf("failed to decode png: %w", err)
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
	outFile, err := os.Create(FULL_COLOR_PATH)
	if err != nil {
		return fmt.Errorf("failed to create outFile: %w", err)
	}
	defer outFile.Close()

	if err = png.Encode(outFile, dst); err != nil {
		return fmt.Errorf("failed to encode png: %w", err)
	}
	slog.Info("Trimmed png")
	return nil
}

func getScreenshot() error {
	const paddingBottom = 200

	wd, err := os.Getwd()
	if err != nil {
		slog.Error("Error getting current working directory", "error", err)
		return err
	}

	cmd := exec.Command(
		"chromium",
		"--headless",
		"--no-sandbox",
		"--disable-gpu",
		fmt.Sprintf("--window-size=%v,%v", EXPORT_WIDTH, EXPORT_HEIGHT+paddingBottom),
		"--force-device-scale-factor=1",
		"--virtual-time-budget=50",
		fmt.Sprintf("--screenshot=%s/%s", wd, FULL_COLOR_PATH),
		fmt.Sprintf("file://%s/%s", wd, OUT_HTML),
	)
	out, err := cmd.Output()
	slog.Info("Took screenshot", "output", out)

	if err != nil {
		slog.Error("Error creating image from HTML", "error", err, "output", out)
		return err
	}
	slog.Info("Created png")

	err = trimScreenshot()
	if err != nil {
		return fmt.Errorf("failed to trim screenshot: %w", err)
	}
	return nil
}

func main() {
	slog.Info("Started")

	if true {
		slog.Info("Taking calendar screenshot")
		now := time.Now()
		offset := (int(now.Weekday()) + 6) % 7
		start := now.AddDate(0, 0, -offset)
		end := start.AddDate(0, 0, NUM_WEEKS*7)

		useCache := false
		createStubEvents := true

		dayEventsMap, err := getData(start, end, createStubEvents, useCache)
		if err != nil {
			slog.Error("Error getting data", "error", err)
			panic(err)
		}
		err = CreateCalendarHTML(start, now, dayEventsMap)
		if err != nil {
			slog.Error("Error creating html", "error", err)
			panic(err)
		}
		err = getScreenshot()
		if err != nil {
			slog.Error("Error getting screenshot", "error", err)
			panic(err)
		}
	}

	err := Dither(FULL_COLOR_PATH, DITHERED_PATH)
	if err != nil {
		slog.Error("Error dithering", "error", err)
		panic(err)
	}
}
