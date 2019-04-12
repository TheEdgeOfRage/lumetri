//
// main.go
// Copyright (C) 2019 pavle <pavle.portic@tilda.center>
//
// Distributed under terms of the BSD-3-Clause license.
//

package main

import (
	"fmt"
	"image"
	"image/color"
	_ "image/jpeg"
	"image/png"
	"log"
	// "math"
	"os"
	"time"
)

func elapsed() func() {
	start := time.Now()
	return func() {
		fmt.Printf("%v\n", time.Since(start))
	}
}

const width uint = 512

func main() {
	defer elapsed()()

	imgfile, err := os.Open("in.jpg")
	if err != nil {
		log.Fatal(err)
	}
	defer imgfile.Close()

	img, _, err := image.Decode(imgfile)
	if err != nil {
		log.Fatal(err)
	}

	bounds := img.Bounds()
	scale := float64(width) / float64(bounds.Max.X - bounds.Min.X)
	coeff := uint8(scale * 4)
	if coeff == 0 {
		coeff = 1
	}
	fmt.Printf("scale: %v, coeff: %v\n", scale, coeff)

	scope := image.NewRGBA(image.Rect(0, 0, int(width), 256))
	scope_bounds := scope.Bounds()

	for y := scope_bounds.Min.Y; y < scope_bounds.Max.Y; y++ {
		for x := scope_bounds.Min.X; x < scope_bounds.Max.X; x++ {
			scope.SetRGBA(x, y, color.RGBA{0, 0, 0, 255})
		}
	}

	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			r, g, b, _:= img.At(x, y).RGBA()
			ir := uint8(r >> 8)
			ig := uint8(g >> 8)
			ib := uint8(b >> 8)

			tr, tg, tb, _ := scope.At(int(float64(x) * scale), 255 - int(ir)).RGBA()
			sr := uint8(tr >> 8)
			sg := uint8(tg >> 8)
			sb := uint8(tb >> 8)
			if sr > 255 - coeff {
				sr = 255
			} else {
				sr += coeff
			}
			scope.SetRGBA(int(float64(x) * scale), 255 - int(ir), color.RGBA{sr, sg, sb, 255})

			tr, tg, tb, _ = scope.At(int(float64(x) * scale), 255 - int(ig)).RGBA()
			sr = uint8(tr >> 8)
			sg = uint8(tg >> 8)
			sb = uint8(tb >> 8)
			if sg > 255 - coeff {
				sg = 255
			} else {
				sg += coeff
			}
			scope.SetRGBA(int(float64(x) * scale), 255 - int(ig), color.RGBA{sr, sg, sb, 255})

			tr, tg, tb, _ = scope.At(int(float64(x) * scale), 255 - int(ib)).RGBA()
			sr = uint8(tr >> 8)
			sg = uint8(tg >> 8)
			sb = uint8(tb >> 8)
			if sb > 255 - coeff {
				sb = 255
			} else {
				sb += coeff
			}
			scope.SetRGBA(int(float64(x) * scale), 255 - int(ib), color.RGBA{sr, sg, sb, 255})
		}
	}

	outFile, _ := os.Create("output.png")
	defer outFile.Close()
	png.Encode(outFile, scope)
}

