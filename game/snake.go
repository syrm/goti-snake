package game

import (
	rl "github.com/gen2brain/raylib-go/raylib"
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

	if startX < 4 && startY < 4 {
		switch rand.Intn(2) {
		case 0:
			snake.direction = Down
		case 1:
			snake.direction = Right
		}
	} else if startX < 4 && startY > snake.grid-5 {
		switch rand.Intn(2) {
		case 0:
			snake.direction = Up
		case 1:
			snake.direction = Right
		}
	} else if startX > snake.grid-5 && startY < 4 {
		switch rand.Intn(2) {
		case 0:
			snake.direction = Down
		case 1:
			snake.direction = Left
		}
	} else if startX > snake.grid-5 && startY > snake.grid-5 {
		switch rand.Intn(2) {
		case 0:
			snake.direction = Up
		case 1:
			snake.direction = Left
		}
	} else if startX < 4 {
		switch rand.Intn(3) {
		case 0:
			snake.direction = Down
		case 1:
			snake.direction = Right
		case 2:
			snake.direction = Up
		}
	} else if startY < 4 {
		switch rand.Intn(3) {
		case 0:
			snake.direction = Down
		case 1:
			snake.direction = Right
		case 2:
			snake.direction = Left
		}
	} else if startX > snake.grid-5 {
		switch rand.Intn(3) {
		case 0:
			snake.direction = Up
		case 1:
			snake.direction = Right
		case 2:
			snake.direction = Down
		}
	} else if startY > snake.grid-5 {
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
	previousCoord := newPosition(0, 0)
	for index := snake.head; index >= snake.head-(snake.length-1); index-- {
		coord := snake.getBody(index)

		if index != snake.head && index == snake.head-(snake.length-1) {
			snake.drawTail(size, previousCoord, coord)
		} else {
			rl.DrawRectangle(
				snake.coordinateConverter.XToPixel(coord.x),
				snake.coordinateConverter.YToPixel(coord.y),
				size,
				size,
				rl.DarkGreen,
			)
		}

		if index == snake.head {
			for _, eye := range snake.xyToLeftEye(size, coord.x, coord.y, snake.direction) {
				rl.DrawCircle(eye[0], eye[1], float32(size/8), rl.DarkPurple)
			}
		}

		previousCoord = newPosition(
			coord.x,
			coord.y,
		)
	}
}

func (snake *Snake) drawTail(size int32, previousBodyCoord Position, coord Position) {
	if previousBodyCoord.x < coord.x {
		rl.DrawCircle(
			snake.coordinateConverter.XToPixel(coord.x),
			snake.coordinateConverter.YToPixel(coord.y)+size/2,
			float32(size/2),
			rl.DarkGreen,
		)
		return
	}

	if previousBodyCoord.x > coord.x {
		rl.DrawCircle(
			snake.coordinateConverter.XToPixel(coord.x)+size,
			snake.coordinateConverter.YToPixel(coord.y)+size/2,
			float32(size/2),
			rl.DarkGreen,
		)
		return
	}

	if previousBodyCoord.y > coord.y {
		rl.DrawCircle(
			snake.coordinateConverter.XToPixel(coord.x)+size/2,
			snake.coordinateConverter.YToPixel(coord.y)+size,
			float32(size/2),
			rl.DarkGreen,
		)
		return
	}

	if previousBodyCoord.y < coord.y {
		rl.DrawCircle(
			snake.coordinateConverter.XToPixel(coord.x)+size/2,
			snake.coordinateConverter.YToPixel(coord.y),
			float32(size/2),
			rl.DarkGreen,
		)
		return
	}
}

func (snake *Snake) xyToLeftEye(size int32, x int32, y int32, direction Direction) [2][2]int32 {
	switch direction {
	case Left:
		return [2][2]int32{
			{
				snake.coordinateConverter.XToPixel(x) + size/8,
				snake.coordinateConverter.YToPixel(y) + size - size/8,
			},
			{
				snake.coordinateConverter.XToPixel(x) + size/8,
				snake.coordinateConverter.YToPixel(y) + size/8,
			},
		}
	case Right:
		return [2][2]int32{
			{
				snake.coordinateConverter.XToPixel(x) + size - size/8,
				snake.coordinateConverter.YToPixel(y) + size - size/8,
			},
			{
				snake.coordinateConverter.XToPixel(x) + size - size/8,
				snake.coordinateConverter.YToPixel(y) - size/8 + 2*size/8,
			},
		}
	case Up:
		return [2][2]int32{
			{
				snake.coordinateConverter.XToPixel(x) + size - size/8,
				snake.coordinateConverter.YToPixel(y) + size/8,
			},
			{
				snake.coordinateConverter.XToPixel(x) + size/8,
				snake.coordinateConverter.YToPixel(y) + size/8,
			},
		}
	case Down:
		return [2][2]int32{
			{
				snake.coordinateConverter.XToPixel(x) + size - size/8,
				snake.coordinateConverter.YToPixel(y) + size - size/8,
			},
			{
				snake.coordinateConverter.XToPixel(x) - size/8 + 2*size/8,
				snake.coordinateConverter.YToPixel(y) + size - size/8,
			},
		}
	}

	return [2][2]int32{}
}

func (snake *Snake) setBody(index int, position Position) {
	snake.body[index%len(snake.body)] = position
}

func (snake *Snake) getBody(index int) Position {
	return snake.body[index%len(snake.body)]
}
