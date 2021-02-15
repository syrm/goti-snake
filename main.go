package main

import (
	gamePkg "github.com/blackprism/goti-snake/game"
	"github.com/gen2brain/raylib-go/raylib"
)

const (
	gameSize  = 600
	border    = 20
	gridSize  = 20
	snakeSize = (gameSize - 2*border) / gridSize
)

func main() {
	rl.SetTraceLog(rl.LogError)

	if gridSize < 3 {
		panic("GridSize should be at less 10")
	}

	position := gamePkg.NewCoordinateConverter(
		20,
		60,
		snakeSize,
	)
	snake := gamePkg.NewSnake(gridSize, position)
	snake.Init()

	game := gamePkg.NewBoard(snake, gameSize, gridSize, 20, 0.2, position)
	game.Init()
	game.SpawnApple()

	gameStatus := gamePkg.NewGame

	windowShouldBeClosed := false

	for !windowShouldBeClosed {
		rl.BeginDrawing()

		if gameStatus == gamePkg.NewGame {
			game.AskNewGame()
		}

		for !windowShouldBeClosed {
			if rl.WindowShouldClose() {
				windowShouldBeClosed = true
				break
			}

			gameStatus = game.Loop()

			if gameStatus == gamePkg.NewGame {
				rl.EndDrawing()
				continue
			}

			if gameStatus != gamePkg.Continue {
				rl.EndDrawing()
				break
			}

			rl.EndDrawing()
		}

		switch gameStatus {
		case gamePkg.Victory:
			game.DisplayVictory()
		case gamePkg.GameOver:
			game.DisplayGameOver()
		}

		rl.EndDrawing()
	}

	rl.CloseWindow()
}
