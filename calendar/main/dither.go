package main

import (
	"fmt"
	"image"
	"image/color"
	"log/slog"
	"os"
	"strconv"

	"golang.org/x/image/bmp"
	"golang.org/x/image/draw"

	dither "github.com/makeworld-the-better-one/dither/v2"
)

var COLOR_6 = []color.Color{
	hexColor("#000000"),
	hexColor("#ffffff"),
	hexColor("#ffff00"),
	hexColor("#ff0000"),
	hexColor("#0000ff"),
	hexColor("#00ff00"),
}

var BLACK_WHITE = []color.Color{
	color.Black,
	color.White,
}

var COLOR_6_EXTENDED = []color.Color{
	hexColor("#000000"),
	hexColor("#ffffff"),
	hexColor("#ffff00"),
	hexColor("#ff0000"),
	hexColor("#0000ff"),
	hexColor("#00ff00"),

	hexColor("#666600"),
	hexColor("#660000"),
	hexColor("#000066"),
	hexColor("#006600"),
}

func hexColor(hex string) color.RGBA {
	values, _ := strconv.ParseUint(string(hex[1:]), 16, 32)
	return color.RGBA{R: uint8(values >> 16), G: uint8((values >> 8) & 0xFF), B: uint8(values & 0xFF), A: 255}
}

func getImageFromFilePath(filePath string) (image.Image, error) {
	f, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	image, _, err := image.Decode(f)
	return image, err
}

func toBlack(v uint32, multiplier float32) uint8 {
	out := uint8(uint32(float32(v)*multiplier) >> 8)
	return out
}

func toWhite(v uint32, multiplier float32) uint8 {
	const UINT32_MAX = 65535
	return uint8(uint32(float32(v)+float32(UINT32_MAX-v)*multiplier) >> 8)
}

func splitColors(palette []color.Color, steps int, fadeToBlack bool) []color.Color {
	// Create variations of each color to give more depth to the display
	// steps = 0 will do nothing
	// steps = 2 will create two extra colors between white and input colors
	for _, p := range palette {
		for step := range steps {
			r, g, b, a := p.RGBA()
			multiplier := float32(step+1) / float32(steps+1)
			f := toWhite
			if fadeToBlack {
				f = toBlack
			}
			palette = append(palette, color.RGBA{
				R: f(r, multiplier),
				G: f(g, multiplier),
				B: f(b, multiplier),
				A: uint8(a),
			})
		}
	}
	return palette

}

func Dither(inPath, outPath string) error {
	skip := false

	outfile, err := os.Create(outPath)
	if err != nil {
		return fmt.Errorf("failed to create output file: %w", err)
	}

	src, err := getImageFromFilePath(inPath)
	if err != nil {
		return fmt.Errorf("source image to dither does not exist: %w", err)
	}

	// Set the expected size that you want:
	dd := image.NewRGBA(image.Rect(0, 0, EXPORT_WIDTH, EXPORT_HEIGHT))

	// copy image in case we want to skip and pass through (not necessary if dithering)
	draw.Draw(dd, dd.Bounds(), src, src.Bounds().Min, draw.Src)

	draw.NearestNeighbor.Scale(dd, dd.Rect, src, src.Bounds(), draw.Over, nil)

	img := image.Image(dd)

	palette := splitColors(COLOR_6, 0, false)

	d := dither.NewDitherer(palette)
	d.Matrix = dither.FloydSteinberg

	if !skip {
		img = d.Dither(img)
	}

	if err = bmp.Encode(outfile, img); err != nil {
		return fmt.Errorf("failed to encode bmp: %w", err)
	}

	slog.Info("Created dither", "outPath", outPath)
	return nil
}
