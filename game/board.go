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
	apple      [2]int32
	snakeSize  int32
	status     Status
	lastStatus Status
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

	board.snakeSize = (board.size - 2*board.border) / board.grid
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
	board.snake.Draw(board.snakeSize)

	board.frames++
}

func (board *Board) CheckApple() {
	if board.snake.AppleEatable(board.apple) {
		board.SpawnApple()
		board.snake.AppleEated()
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

		board.frames = 0
	}

	if rl.IsKeyPressed(rl.KeyLeft) {
		if !board.snake.GoingToDirection(Left) {
			board.status = GameOver
			return
		}

		board.frames = 0
	}

	if rl.IsKeyPressed(rl.KeyUp) {
		if !board.snake.GoingToDirection(Up) {
			board.status = GameOver
			return
		}

		board.frames = 0
	}

	if rl.IsKeyPressed(rl.KeyDown) {
		if !board.snake.GoingToDirection(Down) {
			board.status = GameOver
			return
		}

		board.frames = 0
	}
}

func (board *Board) drawMenu() {
	rl.DrawFPS(10, 10)
}

func (board *Board) drawBackground() {
	rl.DrawRectangle(0, board.menuSize, board.size, board.size, rl.DarkBrown)
	rl.DrawRectangle(board.border, board.menuSize+board.border, board.size-2*board.border, board.size-2*board.border, rl.White)
}

func (board *Board) drawApple() {
	rl.DrawCircle(
		board.position.XToPixel(board.apple[0])+board.size/board.grid/2,
		board.position.YToPixel(board.apple[1])+board.size/board.grid/2,
		float32(board.snakeSize)*0.4,
		rl.Green,
	)
}

func (board *Board) SpawnApple() {
	freeCells := board.snake.GetFreeCells()

	if len(freeCells) == 0 {
		return
	}

	freeCell := freeCells[rand.Intn(len(freeCells))]
	board.apple = [2]int32{
		freeCell[0],
		freeCell[1],
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
