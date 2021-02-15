package game

import (
	rl "github.com/gen2brain/raylib-go/raylib"
	"math/rand"
	"time"
)

func NewBoard(snake *Snake, size int32, grid int32, border int32, speed float32, position CoordinateConverter) *Board {
	return &Board{
		snake:      snake,
		size:       size,
		grid:       grid,
		speed:      speed,
		menuSize:   40,
		border:     border,
		frames:     0,
		position:   position,
		status:     NewGame,
		lastStatus: NewGame,
	}
}

type Board struct {
	snake      *Snake
	size       int32
	grid       int32
	speed      float32
	menuSize   int32
	border     int32
	frames     int32
	position   CoordinateConverter
	apple      Apple
	cellSize   int32
	status     Status
	lastStatus Status
}

type Apple struct {
	x int32
	y int32
}

type Status int

const (
	Victory  Status = 1
	GameOver Status = 2
	Continue Status = 3
	NewGame  Status = 4
)

func (board *Board) Init() {
	rand.Seed(time.Now().UnixNano())

	board.cellSize = (board.size - 2*board.border) / board.grid
	rl.SetConfigFlags(rl.FlagMsaa4xHint | rl.FlagVsyncHint)
	rl.InitWindow(board.size, board.size+board.menuSize, "Goti Board")
	rl.SetTargetFPS(60)
}

func (board *Board) Reset() {
	board.frames = 0
	board.snake.Init()
}

func (board *Board) Draw() {
	rl.ClearBackground(rl.White)
	board.drawMenu()
	board.drawBackground()
	board.drawApple()
	board.snake.Draw(board.cellSize)

	board.frames++
}

func (board *Board) CheckApple() {
	if board.snake.AppleEatable(board.apple) {
		board.snake.AppleEated(board.apple)
		board.SpawnApple()
	}
}

func (board *Board) AutoMove() {
	if board.frames%int32(60*board.speed) == 0 {
		board.snake.Move()
	}
}

func (board *Board) KeyListener() {
	if rl.IsKeyPressed(rl.KeyN) || rl.IsKeyPressed(rl.KeyEnter) || rl.IsKeyPressed(rl.KeyKpEnter) {
		board.NewGame()
		board.status = Continue
		return
	}

	if board.status == NewGame {
		return
	}

	if rl.IsKeyPressed(rl.KeyRight) {
		if !board.snake.GoingToDirection(Right) {
			board.status = GameOver
			return
		}
	}

	if rl.IsKeyPressed(rl.KeyLeft) {
		if !board.snake.GoingToDirection(Left) {
			board.status = GameOver
			return
		}
	}

	if rl.IsKeyPressed(rl.KeyUp) {
		if !board.snake.GoingToDirection(Up) {
			board.status = GameOver
			return
		}
	}

	if rl.IsKeyPressed(rl.KeyDown) {
		if !board.snake.GoingToDirection(Down) {
			board.status = GameOver
			return
		}
	}
}

func (board *Board) drawMenu() {
	rl.DrawFPS(10, 10)
}

func (board *Board) drawBackground() {
	rl.DrawRectangle(0, board.menuSize, board.size, board.size, rl.NewColor(66, 66, 66, 255))
	rl.DrawRectangle(board.border, board.menuSize+board.border, board.size-2*board.border, board.size-2*board.border, rl.White)
}

func (board *Board) drawApple() {
	rl.DrawRectangle(
		board.position.XToPixel(board.apple.x),
		board.position.YToPixel(board.apple.y),
		board.cellSize,
		board.cellSize,
		rl.NewColor(192, 70, 67, 255),
	)
}

func (board *Board) SpawnApple() {
	freeCells := board.snake.GetFreeCells()

	if len(freeCells) == 0 {
		return
	}

	freeCell := freeCells[rand.Intn(len(freeCells))]
	board.apple = Apple{
		x: freeCell[0],
		y: freeCell[1],
	}
}

func (board *Board) GetGameStatus() Status {
	if board.snake.Size() == int(board.grid*board.grid)+1 {
		board.status = Victory
		return board.status
	}

	if board.snake.IsOutside(0, 0, board.grid, board.grid) || board.snake.IsEatingItSelf() {
		board.status = GameOver
		return board.status
	}

	board.status = Continue
	return board.status
}

func (board *Board) Loop() Status {
	board.KeyListener()

	if board.status == NewGame {
		board.Draw()
		switch board.lastStatus {
		case Victory:
			board.DisplayVictory()
		case GameOver:
			board.DisplayGameOver()
		}

		board.DisplayAskNewGame()

		return board.status
	}

	if board.status != Continue {
		return board.status
	}

	if board.GetGameStatus() != Continue {
		return board.status
	}

	board.Draw()
	board.CheckApple()
	board.AutoMove()

	return board.GetGameStatus()
}

func (board *Board) NewGame() {
	board.Reset()
	board.SpawnApple()
}

func (board *Board) AskNewGame() {
	board.status = NewGame
	board.DisplayAskNewGame()
}

func (board *Board) DisplayAskNewGame() {
	rl.DrawText("Press [ENTER] for a New game", board.size/2-270, 80, 20, rl.Black)
}

func (board *Board) DisplayVictory() {
	board.lastStatus = Victory
	board.status = NewGame
	rl.DrawText("Victory !", board.size/2-150, 2, 40, rl.Black)
}

func (board *Board) DisplayGameOver() {
	board.lastStatus = GameOver
	board.status = NewGame
	rl.DrawText("Tu as perdu !", board.size/2-150, 2, 40, rl.Black)
}
