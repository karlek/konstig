// Konstig is a trigonometric strange attractor generator.
package main

import (
	"flag"
	"fmt"
	"log"
	"math"
	"math/rand"
	"os"
	"runtime"
	"time"

	"image"
	"image/color"
	"image/draw"
	"image/png"
)

var (
	// Co-efficients for the strange attractor.
	a, b, c, d, e, f, g, h float64
	// Width and height of the output image.
	width, height int
	// Number of iterations to plot.
	iterations int64
	// Output filename.
	filename string
	// Zoom level.
	zoom float64

	// Frequencies to change the color gradient.
	f1, f2, f3 float64
	// Phases to change the color gradient.
	p1, p2, p3 float64
	// The color gradient.
	rainbow []color.RGBA
)

func init() {
	rand.Seed(time.Now().UnixNano())

	// Co-efficients.
	flag.Float64Var(&a, "a", sign()*2*rand.Float64(), "a co-efficient")
	flag.Float64Var(&b, "b", sign()*2*rand.Float64(), "b co-efficient")
	flag.Float64Var(&c, "c", sign()*2*rand.Float64(), "c co-efficient")
	flag.Float64Var(&d, "d", sign()*2*rand.Float64(), "d co-efficient")
	flag.Float64Var(&e, "e", sign()*2*rand.Float64(), "e co-efficient")
	flag.Float64Var(&f, "f", sign()*2*rand.Float64(), "f co-efficient")
	flag.Float64Var(&g, "g", sign()*2*rand.Float64(), "g co-efficient")
	flag.Float64Var(&h, "h1", sign()*2*rand.Float64(), "h co-efficient")

	flag.Int64Var(&iterations, "i", 1000000, "iterations")

	// Color settings.
	flag.Float64Var(&f1, "f1", 0.01, "frequency1")
	flag.Float64Var(&f2, "f2", 0.01, "frequency2")
	flag.Float64Var(&f3, "f3", 0.01, "frequency3")
	flag.Float64Var(&p1, "p1", 0, "phase1")
	flag.Float64Var(&p2, "p2", 2, "phase2")
	flag.Float64Var(&p3, "p3", 4, "phase3")

	flag.StringVar(&filename, "o", fmt.Sprintf("%f_%f_%f_%f_%f_%f_%f_%f.png", a, b, c, d, e, f, g, h), "output filename")
	/// Debug options.
	// flag.StringVar(&filename, "o", "a.png", "output filename")
	// fmt.Println(fmt.Sprintf("%f_%f_%f_%f.png", a, b, c, d))

	// Dimensions.
	flag.IntVar(&width, "w", 5000, "width")
	flag.IntVar(&height, "h", 5000, "height")
	flag.Float64Var(&zoom, "z", 400.0, "zoom")
}

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	flag.Parse()
	err := attractor()
	if err != nil {
		log.Fatalln(err)
	}
}

/// Better documentation.
// gradient creates a list of colors.
func gradient(f1, f2, f3, p1, p2, p3 float64, center, width, length int) []color.RGBA {
	var pal = make([]color.RGBA, length)

	// We make smoother transition from color to color.
	var s float64
	for i := 0; i < length; i++ {
		s += 0.2
		r := int(math.Sin(f1*s+p1)*float64(width)) + center
		g := int(math.Sin(f2*s+p2)*float64(width)) + center
		b := int(math.Sin(f3*s+p3)*float64(width)) + center
		pal[i] = color.RGBA{uint8(r), uint8(g), uint8(b), 255}
	}
	return pal
}

// attractor creates a trigonometric strange attractor.
func attractor() (err error) {
	// Variable gradient size based on image dimensions.
	gsize := int(dist(width/2, height/2, 0, 0))

	// Smooth gradient.
	rainbow = gradient(f1, f2, f3, p1, p2, p3, 127, 128, gsize)

	// Output image with black background.
	img := image.NewRGBA(image.Rect(0, 0, width, height))
	// black := color.RGBA{0, 0, 0, 255}
	black := color.RGBA{255, 255, 255, 255}
	draw.Draw(img, img.Bounds(), &image.Uniform{black}, image.ZP, draw.Src)

	// Starting point. Arbitrary non-zero value. Because the sin function would
	// be constantly zero.
	x, y := 0.1, 0.1

	var p image.Point
	var i int64
	for i = 0; i < 100; i++ {
		xnew := a*math.Cos(y*b) + c*math.Cos(x*d)
		ynew := e*math.Sin(x*f) + g*math.Sin(y*h)

		// Save the new point.
		x = xnew
		y = ynew
	}
	pChan := make(chan image.Point, 10000)
	go func() {
		for i = 0; i < iterations; i++ {
			// Compute a new point using the trigonometric strange attractor
			// equations.
			xnew := a*math.Cos(y*b) + c*math.Cos(x*d)
			ynew := e*math.Sin(x*f) + g*math.Sin(y*h)

			// Save the new point.
			x = xnew
			y = ynew

			// Center the attractor both horizontally and vertically.
			p.X = int(zoom*x) + width/2
			p.Y = int(zoom*y) + height/2

			pChan <- p
		}
		close(pChan)
	}()
	for p := range pChan {
		setPt(&p, img)
	}

	if allBlack(img) {
		fmt.Println(fmt.Sprintf("konstig -a %f -b %f -c %f -d %f -e %f -f %f -g %f -h1 %f -o a.png", a, b, c, d, e, f, g, h))
		return nil
	}
	// Create the output image file.
	return save(img)
}

func allBlack(i *image.RGBA) bool {
	atLeast := width * height / 1000
	num := 0
	for x := 0; x < width; x++ {
		for y := 0; y < height; y++ {
			r, g, b, _ := i.At(x, y).RGBA()
			if !isBlack(r, g, b) {
				num++
			}
		}
	}
	fmt.Println(num, atLeast)
	return num < atLeast
}

// setPt sets the point p to the color in the gradient relative to the distance
// from origo.
func setPt(p *image.Point, img *image.RGBA) {
	// Generate the current pixels color.
	col := gradientIndex(dist(width/2, height/2, p.X, p.Y))

	// Take the old color of the pixel.
	r, g, b, _ := img.At(p.X, p.Y).RGBA()
	if isBlack(r, g, b) {
		img.Set(p.X, p.Y, col)
		return
	}

	// Make the current pixel a little brighter.
	newc := add(col, color.RGBA{uint8(r), uint8(g), uint8(b), 255})
	img.Set(p.X, p.Y, newc)
}

// isBlack returns true if the color is black.
func isBlack(r, g, b uint32) bool {
	return uint8(r) == 255 && uint8(g) == 255 && uint8(b) == 255
}

// save creates an output image file.
func save(img *image.RGBA) (err error) {
	out, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer out.Close()

	fmt.Println("[!]    Done:", filename)
	return png.Encode(out, img)
}

// sign randomly returns either 1 or -1.
func sign() float64 {
	a := rand.Intn(2)
	if a == 0 {
		return 1.0
	}
	return -1.0
}

// dist calculates the distance between two points.
func dist(cx, cy, x, y int) float64 {
	dx, dy := float64(cx-x), float64(cy-y)
	return math.Sqrt(math.Pow(dx, 2) + math.Pow(dy, 2))
}

// gradientIndex returns the color at the index nearest c.
func gradientIndex(c float64) color.RGBA {
	return rainbow[int(c)%cap(rainbow)]
}

// add adds two colors together to make a brighter color.
func add(c1 color.RGBA, c2 color.RGBA) color.RGBA {
	return color.RGBA{c1.R + c2.R, c1.G + c2.G, c1.B + c2.B, 255}
}
