package main

import (
	"fmt"
	"image"
	"image/color"
	"image/png"
	"log"
	"os"
	"strconv"

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

func Dither(inPath, outPath string) {
	src, err := getImageFromFilePath(inPath)

	// Set the expected size that you want:
	dd := image.NewRGBA(image.Rect(0, 0, 1600, 960))

	// Resize:
	draw.NearestNeighbor.Scale(dd, dd.Rect, src, src.Bounds(), draw.Over, nil)

	img := image.Image(dd)
	if err != nil {
		log.Panicln("Image does not exist", img)
	}

	d := dither.NewDitherer(COLOR_6_EXTENDED)
	d.Matrix = dither.FloydSteinberg

	img = d.Dither(img)

	outfile, err := os.Create(outPath)
	if err != nil {
		panic(err)
	}
	if err = png.Encode(outfile, img); err != nil {
		log.Printf("failed to encode: %v", err)
	}

	fmt.Println("Created dither")
}
