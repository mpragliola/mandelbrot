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
const width = 1200
const height = 800
const offset = 2 + 0.5i
const zoom = 1.5
const aliasFactor = 3.0

func main() {
	s := stopwatch.NewStart()

	upLeft := image.Point{0, 0}
	lowRight := image.Point{width, height}
	img := image.NewRGBA(image.Rectangle{upLeft, lowRight})

	wg := sync.WaitGroup{}
	for ry := 0; ry < height; ry++ {
		wg.Add(1)
		go func(ry int) {
			defer wg.Done()

			var y, h, h2 float64
			for rx := 0; rx < width; rx++ {
				h = 0
				h2 = 0

				// Anti-aliasing: for each screen pixel, calculate the average on a NxN
				// subpixel square, where N is the aliasFactor
				for subPixelY := 0.0; subPixelY < aliasFactor; subPixelY++ {
					y = float64(ry) + subPixelY/aliasFactor

					for subPixelX := 0.0; subPixelX < aliasFactor; subPixelX++ {
						x := float64(rx) + subPixelX/aliasFactor

						// scale by zoom and translate by offset from screen coordinates to complex plane
						fac := complex(float64(height), 0.0) * zoom
						c := complex(x, y)/fac - offset
						i := 0.0

						// Iterate until we reach maxIteration or the value is > 2.0 ; a theorem guarantees
						// that z[n] will diverge in the latter case
						for z := c; cmplx.Abs(z) < 2.0 && i < maxIterations; i++ {

							// Mandelbrot's fractal iterative formula (z[n+1] = z[n]² + c)
							z = z*z + c
						}

						// Calculate shade based on the iterations
						h += .8 * i / maxIterations * 255.0
						// Change curve by using square
						h2 += (i / maxIterations) * (i / maxIterations) * 255.0
					}
				}

				sq := float64(aliasFactor * aliasFactor)
				img.Set(rx, ry, color.RGBA{uint8(h / sq), uint8(h / sq), uint8(h2 / sq), 0xff})
			}

		}(ry)
	}

	wg.Wait()

	f, _ := os.Create("image.png")
	png.Encode(f, img)

	s.Dump()
}
