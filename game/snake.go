package game

import (
	rl "github.com/gen2brain/raylib-go/raylib"
	"math"
	"math/rand"
	"time"
)

type Direction int

const (
	Left  Direction = 0
	Right Direction = 1
	Up    Direction = 2
	Down  Direction = 3
)

func newPosition(x int32, y int32) Position {
	return Position{
		x,
		y,
	}
}

type Position struct {
	x int32
	y int32
}

func NewSnake(grid int32, coordinateConverter CoordinateConverter) *Snake {
	return &Snake{
		grid:                grid,
		head:                0,
		length:              0,
		body:                make([]Position, grid*grid+1),
		direction:           Up,
		coordinateConverter: coordinateConverter,
		needToGrow:          false,
	}
}

type Snake struct {
	grid                int32
	head                int
	length              int
	body                []Position
	direction           Direction
	coordinateConverter CoordinateConverter
	needToGrow          bool
}

func (snake *Snake) Init() {
	rand.Seed(time.Now().UnixNano())

	startX := rand.Int31n(snake.grid)
	startY := rand.Int31n(snake.grid)

	snake.setBody(snake.head, newPosition(startX, startY))
	snake.head = 0
	snake.length = 1

	if startX <= int32(math.Floor(float64(snake.grid)*0.33)) && startY <= int32(math.Floor(float64(snake.grid)*0.33)) { // Left Top
		switch rand.Intn(2) {
		case 0:
			snake.direction = Down
		case 1:
			snake.direction = Right
		}
	} else if startX <= int32(math.Floor(float64(snake.grid)*0.33)) && startY >= int32(math.Floor(float64(snake.grid)*0.66)) { // Left Bottom
		switch rand.Intn(2) {
		case 0:
			snake.direction = Up
		case 1:
			snake.direction = Right
		}
	} else if startX >= int32(math.Floor(float64(snake.grid)*0.66)) && startY <= int32(math.Floor(float64(snake.grid)*0.33)) { // Right Top
		switch rand.Intn(2) {
		case 0:
			snake.direction = Down
		case 1:
			snake.direction = Left
		}
	} else if startX >= int32(math.Floor(float64(snake.grid)*0.66)) && startY >= int32(math.Floor(float64(snake.grid)*0.66)) { // Right Bottom
		switch rand.Intn(2) {
		case 0:
			snake.direction = Up
		case 1:
			snake.direction = Left
		}
	} else if startX <= int32(math.Floor(float64(snake.grid)*0.33)) { // Left
		switch rand.Intn(3) {
		case 0:
			snake.direction = Down
		case 1:
			snake.direction = Right
		case 2:
			snake.direction = Up
		}
	} else if startY <= int32(math.Floor(float64(snake.grid)*0.33)) { // Top
		switch rand.Intn(3) { // Left
		case 0:
			snake.direction = Down
		case 1:
			snake.direction = Right
		case 2:
			snake.direction = Left
		}
	} else if startX >= int32(math.Floor(float64(snake.grid)*0.66)) { // Right
		switch rand.Intn(3) {
		case 0:
			snake.direction = Up
		case 1:
			snake.direction = Right
		case 2:
			snake.direction = Down
		}
	} else if startY >= int32(math.Floor(float64(snake.grid)*0.66)) { // Bottom
		switch rand.Intn(3) {
		case 0:
			snake.direction = Up
		case 1:
			snake.direction = Right
		case 2:
			snake.direction = Left
		}
	}
}

func (snake *Snake) Size() int {
	return snake.length
}

func (snake *Snake) GetFreeCells() [][]int32 {
	var freeCells [][]int32
	for gridX := int32(0); gridX < snake.grid; gridX++ {
		for gridY := int32(0); gridY < snake.grid; gridY++ {
			found := false
			for index := snake.head; index >= snake.head-(snake.length-1); index-- {
				if snake.getBody(index).x == gridX && snake.getBody(index).y == gridY {
					found = true
					break
				}
			}
			if found == false {
				freeCells = append(freeCells, []int32{gridX, gridY})
			}
		}
	}

	return freeCells
}

func (snake *Snake) AppleEatable(apple [2]int32) bool {
	if apple[0] == snake.getBody(snake.head).x && apple[1] == snake.getBody(snake.head).y {
		return true
	}

	return false
}

func (snake *Snake) AppleEated() {
	snake.needToGrow = true
}

func (snake *Snake) Move() bool {
	if snake.needToGrow {
		snake.needToGrow = false
		snake.length++
	}

	snake.head++

	snake.setBody(snake.head, newPosition(
		snake.getBody(snake.head-1).x,
		snake.getBody(snake.head-1).y,
	))

	switch snake.direction {
	case Up:
		snake.setBody(snake.head, newPosition(
			snake.getBody(snake.head).x,
			snake.getBody(snake.head).y-1,
		))
	case Down:
		snake.setBody(snake.head, newPosition(
			snake.getBody(snake.head).x,
			snake.getBody(snake.head).y+1,
		))
	case Left:
		snake.setBody(snake.head, newPosition(
			snake.getBody(snake.head).x-1,
			snake.getBody(snake.head).y,
		))
	case Right:
		snake.setBody(snake.head, newPosition(
			snake.getBody(snake.head).x+1,
			snake.getBody(snake.head).y,
		))
	}

	return true
}

func (snake *Snake) IsOutside(x1 int32, y1 int32, x2 int32, y2 int32) bool {
	if snake.getBody(snake.head).x < x1 || snake.getBody(snake.head).y < y1 || snake.getBody(snake.head).x >= x2 || snake.getBody(snake.head).y >= y2 {
		return true
	}

	return false
}

func (snake *Snake) IsEatingItSelf() bool {
	head := snake.getBody(snake.head)

	for index := snake.head; index >= snake.head-(snake.length-1); index-- {
		if index != snake.head && head.x == snake.getBody(index).x && head.y == snake.getBody(index).y {
			return true
		}
	}

	return false
}

func (snake *Snake) GoingToDirection(direction Direction) bool {
	switch snake.direction {
	case Up:
		if direction == Down {
			return false
		}
	case Down:
		if direction == Up {
			return false
		}
	case Left:
		if direction == Right {
			return false
		}
	case Right:
		if direction == Left {
			return false
		}
	}

	snake.direction = direction
	snake.Move()

	return true
}

func (snake *Snake) Draw(size int32) {
	for index := snake.head; index >= snake.head-(snake.length-1); index-- {
		coord := snake.getBody(index)

		rl.DrawRectangle(
			snake.coordinateConverter.XToPixel(coord.x),
			snake.coordinateConverter.YToPixel(coord.y),
			size,
			size,
			rl.NewColor(157, 196, 98, 255),
		)

		if index == snake.head {
			switch snake.direction {
			case Up:
				rl.DrawRectangleGradientV(
					snake.coordinateConverter.XToPixel(coord.x),
					snake.coordinateConverter.YToPixel(coord.y),
					size,
					size,
					rl.NewColor(47, 78, 0, 255),
					rl.NewColor(157, 196, 98, 255),
				)
			case Down:
				rl.DrawRectangleGradientV(
					snake.coordinateConverter.XToPixel(coord.x),
					snake.coordinateConverter.YToPixel(coord.y),
					size,
					size,
					rl.NewColor(157, 196, 98, 255),
					rl.NewColor(47, 78, 0, 255),
				)
			case Left:
				rl.DrawRectangleGradientH(
					snake.coordinateConverter.XToPixel(coord.x),
					snake.coordinateConverter.YToPixel(coord.y),
					size,
					size,
					rl.NewColor(47, 78, 0, 255),
					rl.NewColor(157, 196, 98, 255),
				)
			case Right:
				rl.DrawRectangleGradientH(
					snake.coordinateConverter.XToPixel(coord.x),
					snake.coordinateConverter.YToPixel(coord.y),
					size,
					size,
					rl.NewColor(157, 196, 98, 255),
					rl.NewColor(47, 78, 0, 255),
				)
			}
		}
	}
}

func (snake *Snake) setBody(index int, position Position) {
	snake.body[index%len(snake.body)] = position
}

func (snake *Snake) getBody(index int) Position {
	return snake.body[index%len(snake.body)]
}
