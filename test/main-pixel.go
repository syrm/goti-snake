// +build pixel

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
		VSync:  true,
	}

	win, err := pixelgl.NewWindow(cfg)

	if err != nil {
		panic(fmt.Sprintf("NewWindow: %s", err))

	}

	win.Clear(colornames.Green)
	win.SetSmooth(true)

	var (
		frames = 0
		second = time.Tick(time.Second)
		crawlSpeed = time.Tick(250 * time.Millisecond)
	)

	imd := imdraw.New(nil)
	imd.Color = pixel.RGB(1, 1, 1)
	imd.Push(pixel.V(border-1, border-1))
	imd.Push(pixel.V(border-1, size - border+1))
	imd.Push(pixel.V(size - border+1, size - border+1))
	imd.Push(pixel.V(size - border+1, border-1))
	imd.Polygon(0)
	imd.Draw(win)

	var snake [][]int32

	rand.Seed(time.Now().UnixNano())

	startX := rand.Int31n(gridSize)
	startY := rand.Int31n(gridSize)
	startX = 2
	startY = gridSize - 3

	snake = append(snake, []int32{startX, startY})

	currentDirection := Up
	if startX < 4 && startY < 4 {
		switch  rand.Intn(2) {
		case 0:
			currentDirection = Up
		case 1:
			currentDirection = Right
		}
	} else if startX < 4 {
		switch  rand.Intn(3) {
		case 0:
			currentDirection = Up
		case 1:
			currentDirection = Right
		case 2:
			currentDirection = Down
		}
	} else if startY < 4 {
		switch  rand.Intn(3) {
		case 0:
			currentDirection = Up
		case 1:
			currentDirection = Right
		case 2:
			currentDirection = Left
		}
	} else if startX > gridSize - 5 && startY > gridSize - 5 {
		switch  rand.Intn(2) {
		case 0:
			currentDirection = Down
		case 1:
			currentDirection = Left
		}
	} else if startX > gridSize - 5 {
		switch  rand.Intn(3) {
		case 0:
			currentDirection = Down
		case 1:
			currentDirection = Left
		case 2:
			currentDirection = Up
		}
	} else if startY > gridSize - 5 {
		switch  rand.Intn(3) {
		case 0:
			currentDirection = Down
		case 1:
			currentDirection = Left
		case 2:
			currentDirection = Right
		}
	} else if startX < 4 && startY > gridSize - 5 {
		switch  rand.Intn(2) {
		case 0:
			currentDirection = Down
		case 1:
			currentDirection = Right
		}
	} else if startX > gridSize - 5 && startY < 4 {
		switch  rand.Intn(2) {
		case 0:
			currentDirection = Up
		case 1:
			currentDirection = Left
		}
	}

	last := time.Now()
	_forWin: for !win.Closed() {
		dt := time.Since(last).Seconds()
		last = time.Now()
		_ = dt

		if win.JustPressed(pixelgl.KeyRight) {
			if currentDirection == Left {
				break
			}
			currentDirection = Right
		}

		if win.JustPressed(pixelgl.KeyLeft) {
			if currentDirection == Right {
				break
			}
			currentDirection = Left
		}

		if win.JustPressed(pixelgl.KeyUp) {
			if currentDirection == Down {
				break
			}
			currentDirection = Up
		}

		if win.JustPressed(pixelgl.KeyDown) {
			if currentDirection == Up {
				break
			}
			currentDirection = Down
		}

		win.Clear(colornames.Darkseagreen)
		drawSnake(imd, snake, win, currentDirection)
		win.Update()

		frames++
		select {
		case <-second:
			win.SetTitle(fmt.Sprintf("%s | FPS: %d", cfg.Title, frames))
			frames = 0
		case <-crawlSpeed:
			clearTailSnake(imd, snake, win)

			switch currentDirection {
			case Up:
				snake[0][1]++
			case Down:
				snake[0][1]--
			case Left:
				snake[0][0]--
			case Right:
				snake[0][0]++
			}

			if snake[0][0] < 0 || snake[0][1] < 0 || snake[0][0] >= gridSize || snake[0][1] >= gridSize {
				break _forWin
			}
		default:
		}
	}

	println("Tu as perdu")
}

func clearTailSnake(imd *imdraw.IMDraw, snake [][]int32, win *pixelgl.Window) {
	for index, coord := range snake {
		if index == len(snake)-1 {
			imd.Color = colornames.White
			imd.Push(xyToTopLeftCorner(coord[0], coord[1]))
			imd.Push(xyToTopRightCorner(coord[0], coord[1]))
			imd.Push(xyToBottomRightCorner(coord[0], coord[1]))
			imd.Push(xyToBottomLeftCorner(coord[0], coord[1]))
			imd.Polygon(0)
			imd.Draw(win)
		}
	}
}

func drawSnake(imd *imdraw.IMDraw, snake [][]int32, win *pixelgl.Window, direction Direction) {
	for index, coord := range snake {
		imd.Color = colornames.Darkgreen
		imd.Push(xyToTopLeftCorner(coord[0], coord[1]))
		imd.Push(xyToTopRightCorner(coord[0], coord[1]))
		imd.Push(xyToBottomRightCorner(coord[0], coord[1]))
		imd.Push(xyToBottomLeftCorner(coord[0], coord[1]))
		imd.Polygon(0)

		if index == 0 {
			imd.Color = colornames.Darkred
			for _, eye := range xyToLeftEye(coord[0], coord[1], direction) {
				imd.Push(eye)
			}
			imd.Circle(snakeSize / 8, 0)
		}

		imd.Draw(win)
	}
}

func xyToLeftEye(x int32, y int32, direction Direction) [2]pixel.Vec {
	switch direction {
	case Left:
		return [2]pixel.Vec{
			pixel.V(
				border + float64(x*snakeSize+snakeSize/8),
				border + float64(y*snakeSize+snakeSize-snakeSize/8),
			),
			pixel.V(
				border + float64(x*snakeSize+snakeSize/8),
				border + float64(y*snakeSize+snakeSize/8),
			),
		}
	case Right:
		return [2]pixel.Vec{
			pixel.V(
				border + float64(x*snakeSize+snakeSize-snakeSize/8),
				border + float64(y*snakeSize+snakeSize-snakeSize/8),
			),
			pixel.V(
				border + float64(x*snakeSize+snakeSize-snakeSize/8),
				border + float64(y*snakeSize-snakeSize/8+2*snakeSize/8),
			),
		}
	case Down:
		return [2]pixel.Vec{
			pixel.V(
				border + float64(x*snakeSize+snakeSize-snakeSize/8),
				border + float64(y*snakeSize+snakeSize/8),
			),
			pixel.V(
				border + float64(x*snakeSize+snakeSize/8),
				border + float64(y*snakeSize+snakeSize/8),
			),
		}
	case Up:
		return [2]pixel.Vec{
			pixel.V(
				border + float64(x*snakeSize+snakeSize-snakeSize/8),
				border + float64(y*snakeSize+snakeSize-snakeSize/8),
			),
			pixel.V(
				border + float64(x*snakeSize-snakeSize/8+2*snakeSize/8),
				border + float64(snakeSize+y*snakeSize-snakeSize/8),
			),
		}
	}

	return [2]pixel.Vec{}
}

func xyToTopLeftCorner(x int32, y int32) pixel.Vec {
	return pixel.V(
		border + float64(snakeSize + x * snakeSize),
		border + float64(y * snakeSize),
	)
}

func xyToTopRightCorner(x int32, y int32) pixel.Vec {
	return pixel.V(
		border + float64(snakeSize + x * snakeSize),
		border + float64(snakeSize + y * snakeSize),
	)
}

func xyToBottomLeftCorner(x int32, y int32) pixel.Vec {
	return pixel.V(
		border + float64(x * snakeSize),
		border + float64(y * snakeSize),
	)
}

func xyToBottomRightCorner(x int32, y int32) pixel.Vec {
	return pixel.V(
		border + float64(x * snakeSize),
		border + float64(snakeSize + y * snakeSize),
	)
}
