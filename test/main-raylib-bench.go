// +build raylibbench

package main

import (
	"github.com/gen2brain/raylib-go/raylib"
	"math/rand"
	"time"
)

const (
	window_size = 600
)

// rl.DrawLine(18, 42, screenWidth-18, 42, rl.Black)
// rl.DrawRectangleLines(screenWidth/4*2-40, 320, 80, 60, rl.Orange)

func main()  {
	rl.InitWindow(window_size, window_size, "Goti Snake")

	rl.SetTargetFPS(6000)

	squaresCount := 20000
	x := make([]int32, squaresCount)
	y := make([]int32, squaresCount)

	rand.Seed(time.Now().UnixNano())

	for i := 0; i < squaresCount; i++ {
		x[i] = rand.Int31n(window_size)
		y[i] = rand.Int31n(window_size)
	}

	for !rl.WindowShouldClose() {
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
		}

		rl.BeginDrawing()
		rl.ClearBackground(rl.RayWhite)
		for i := 0; i < squaresCount; i++ {
			rl.DrawRectangle(x[i], y[i], 8, 8, rl.Red)
		}
		rl.DrawFPS(2, 2)
		rl.EndDrawing()
	}

	rl.CloseWindow()
}
