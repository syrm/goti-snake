// +build pixelbench

package main

import (
	"fmt"
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
	"github.com/faiface/pixel/pixelgl"
	"golang.org/x/image/colornames"
	_ "image/png"
	"math/rand"
	"time"
)

const snakeBody = "./assets/snake-body.png"
const test = "./assets/hiking.png"

const size = 600
const border = 10
const gridSize = 20
const snakeSize = (size - 2 * border) / gridSize

type Direction int

const (
	Left Direction = 0
	Right Direction = 1
	Up Direction = 2
	Down Direction = 3
)

func main() {
	_ = snakeBody
	_ = test

	pixelgl.Run(run)
}

func run() {
	cfg := pixelgl.WindowConfig{
		Title:  "Goti Snake ♿",
		Bounds: pixel.R(0, 0, size, size),
		VSync:  false,
	}

	win, err := pixelgl.NewWindow(cfg)

	if err != nil {
		panic(fmt.Sprintf("NewWindow: %s", err))

	}

	win.Clear(colornames.Green)
	win.SetSmooth(false)

	frames := 0
	second := time.Tick(time.Second)
	squaresCount := 20000

	x := make([]int32, squaresCount)
	y := make([]int32, squaresCount)

	rand.Seed(time.Now().UnixNano())

	for i := 0; i < squaresCount; i++ {
		x[i] = rand.Int31n(600)
		y[i] = rand.Int31n(600)
	}

	last := time.Now()
	imd := imdraw.New(nil)

	for !win.Closed() {
		dt := time.Since(last).Seconds()
		last = time.Now()
		_ = dt

		imd.Clear()
		for i := 0; i < squaresCount; i++ {
			xfactor := int32(1)
			if rand.Intn(2) == 0 {
				xfactor *= -1
			}

			yfactor := int32(1)
			if rand.Intn(2) == 0 {
				yfactor *= -1
			}

			x[i] += xfactor
			y[i] += yfactor

			imd.Color = pixel.RGB(1, 0, 0)
			imd.Push(pixel.V(0 + float64(x[i]), 0 + float64(y[i])))
			imd.Push(pixel.V(0 + float64(x[i]), 8 + float64(y[i])))
			imd.Push(pixel.V(8 + float64(x[i]), 8 + float64(y[i])))
			imd.Push(pixel.V(8 + float64(x[i]), 0 + float64(y[i])))
			imd.Polygon(0)
		}

		win.Clear(colornames.White)
		imd.Draw(win)
		win.Update()

		frames++
		select {
		case <-second:
			win.SetTitle(fmt.Sprintf("%s | FPS: %d", cfg.Title, frames))
			frames = 0
		default:
		}
	}

	println("Tu as perdu")
}
