package main

import (
	"image"
	"image/color"
	"image/png"
	"math/cmplx"
	"os"
	"sync"

	"github.com/mpragliola/stopwatch"
)

const maxIterations = 40

func main() {
	s := stopwatch.NewStart()

	width := 1200
	height := 800

	upLeft := image.Point{0, 0}
	lowRight := image.Point{width, height}

	img := image.NewRGBA(image.Rectangle{upLeft, lowRight})

	wg := sync.WaitGroup{}

	for y := 0; y < height; y++ {
		wg.Add(1)
		go func(y int) {
			defer wg.Done()

			for x := 0; x < width; x++ {
				// scale coordinates
				c := complex(float64(x), float64(y))/400.0 - 2 - 1i
				i := 0.0
				// A theorem guarantees that z will diverge if it crosses Abs(z) = 2
				for z := c; cmplx.Abs(z) < 2.0 && i < maxIterations; i++ {
					// Mandelbrot's fractal iterative formula (z[n+1] = z[n]² + c)
					z = z*z + c
				}

				// Calculate shate based on the iterations
				h := uint8(i / maxIterations * 255.0)
				h2 := uint8((i / maxIterations) * (i / 90) * 255.0)
				img.Set(x, y, color.RGBA{h, h, h2, 0xff})
			}

		}(y)
	}

	wg.Wait()

	f, _ := os.Create("image.png")
	png.Encode(f, img)

	s.Dump()
}
