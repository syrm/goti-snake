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
	applesEated         []Apple
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

func (snake *Snake) AppleEatable(apple Apple) bool {
	if apple.x == snake.getBody(snake.head).x && apple.y == snake.getBody(snake.head).y {
		return true
	}

	return false
}

func (snake *Snake) AppleEated(apple Apple) {
	snake.applesEated = append(snake.applesEated, apple)
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

	return true
}

func (snake *Snake) Draw(size int32) {
	applesToDigest := make([]bool, len(snake.applesEated))
	degradedStep := 140 / snake.length

	for index := snake.head; index >= snake.head-(snake.length-1); index-- {
		coord := snake.getBody(index)
		bodyColor := []rl.Color{
			rl.NewColor(157, 196, 98, 255),
			rl.NewColor(157, 196, 98, 255),
			rl.NewColor(157, 196, 98, 255),
			rl.NewColor(157, 196, 98, 255),
		}

		if index < snake.head {
			coordNext := snake.getBody(index + 1)

			// haut gauche
			// bas gauche
			// bas droite
			// haut droite

			if coordNext.x > coord.x {
				bodyColor[0] = rl.NewColor(
					uint8(47+degradedStep*(snake.head-index+1)),
					uint8(78+degradedStep*(snake.head-index+1)),
					uint8(0+degradedStep*(snake.head-index+1)),
					255,
				)

				bodyColor[1] = rl.NewColor(
					uint8(47+degradedStep*(snake.head-index+1)),
					uint8(78+degradedStep*(snake.head-index+1)),
					uint8(0+degradedStep*(snake.head-index+1)),
					255,
				)

				bodyColor[2] = rl.NewColor(
					uint8(47+degradedStep*(snake.head-index)),
					uint8(78+degradedStep*(snake.head-index)),
					uint8(0+degradedStep*(snake.head-index)),
					255,
				)

				bodyColor[3] = rl.NewColor(
					uint8(47+degradedStep*(snake.head-index)),
					uint8(78+degradedStep*(snake.head-index)),
					uint8(0+degradedStep*(snake.head-index)),
					255,
				)
			}

			if coordNext.y > coord.y {
				bodyColor[0] = rl.NewColor(
					uint8(47+degradedStep*(snake.head-index+1)),
					uint8(78+degradedStep*(snake.head-index+1)),
					uint8(0+degradedStep*(snake.head-index+1)),
					255,
				)

				bodyColor[3] = rl.NewColor(
					uint8(47+degradedStep*(snake.head-index+1)),
					uint8(78+degradedStep*(snake.head-index+1)),
					uint8(0+degradedStep*(snake.head-index+1)),
					255,
				)

				bodyColor[1] = rl.NewColor(
					uint8(47+degradedStep*(snake.head-index)),
					uint8(78+degradedStep*(snake.head-index)),
					uint8(0+degradedStep*(snake.head-index)),
					255,
				)

				bodyColor[2] = rl.NewColor(
					uint8(47+degradedStep*(snake.head-index)),
					uint8(78+degradedStep*(snake.head-index)),
					uint8(0+degradedStep*(snake.head-index)),
					255,
				)
			}

			if coordNext.x < coord.x {
				bodyColor[2] = rl.NewColor(
					uint8(47+degradedStep*(snake.head-index+1)),
					uint8(78+degradedStep*(snake.head-index+1)),
					uint8(0+degradedStep*(snake.head-index+1)),
					255,
				)

				bodyColor[3] = rl.NewColor(
					uint8(47+degradedStep*(snake.head-index+1)),
					uint8(78+degradedStep*(snake.head-index+1)),
					uint8(0+degradedStep*(snake.head-index+1)),
					255,
				)

				bodyColor[0] = rl.NewColor(
					uint8(47+degradedStep*(snake.head-index)),
					uint8(78+degradedStep*(snake.head-index)),
					uint8(0+degradedStep*(snake.head-index)),
					255,
				)

				bodyColor[1] = rl.NewColor(
					uint8(47+degradedStep*(snake.head-index)),
					uint8(78+degradedStep*(snake.head-index)),
					uint8(0+degradedStep*(snake.head-index)),
					255,
				)
			}

			if coordNext.y < coord.y {
				bodyColor[1] = rl.NewColor(
					uint8(47+degradedStep*(snake.head-index+1)),
					uint8(78+degradedStep*(snake.head-index+1)),
					uint8(0+degradedStep*(snake.head-index+1)),
					255,
				)

				bodyColor[2] = rl.NewColor(
					uint8(47+degradedStep*(snake.head-index+1)),
					uint8(78+degradedStep*(snake.head-index+1)),
					uint8(0+degradedStep*(snake.head-index+1)),
					255,
				)

				bodyColor[0] = rl.NewColor(
					uint8(47+degradedStep*(snake.head-index)),
					uint8(78+degradedStep*(snake.head-index)),
					uint8(0+degradedStep*(snake.head-index)),
					255,
				)

				bodyColor[3] = rl.NewColor(
					uint8(47+degradedStep*(snake.head-index)),
					uint8(78+degradedStep*(snake.head-index)),
					uint8(0+degradedStep*(snake.head-index)),
					255,
				)
			}
		}

		for appleIndex, appleEated := range snake.applesEated {
			if coord.x == appleEated.x && coord.y == appleEated.y {
				coordNext := snake.getBody(index + 1)

				degradedStep := int(140 / math.Max(1, float64(snake.length/2)))

				// haut gauche
				// bas gauche
				// bas droite
				// haut droite

				if coordNext.x > coord.x {
					bodyColor[0] = rl.NewColor(
						uint8(47+math.Min(140, float64(degradedStep*(snake.head-index+1)))),
						uint8(78+math.Min(140, float64(degradedStep*(snake.head-index+1)))),
						uint8(0+math.Min(140, float64(degradedStep*(snake.head-index+1)))),
						255,
					)

					bodyColor[1] = rl.NewColor(
						uint8(47+math.Min(140, float64(degradedStep*(snake.head-index+1)))),
						uint8(78+math.Min(140, float64(degradedStep*(snake.head-index+1)))),
						uint8(0+math.Min(140, float64(degradedStep*(snake.head-index+1)))),
						255,
					)

					bodyColor[2] = rl.NewColor(
						uint8(47+math.Min(140, float64(degradedStep*(snake.head-index+0)))),
						uint8(78+math.Min(140, float64(degradedStep*(snake.head-index+0)))),
						uint8(0+math.Min(140, float64(degradedStep*(snake.head-index+0)))),
						255,
					)

					bodyColor[3] = rl.NewColor(
						uint8(47+math.Min(140, float64(degradedStep*(snake.head-index+0)))),
						uint8(78+math.Min(140, float64(degradedStep*(snake.head-index+0)))),
						uint8(0+math.Min(140, float64(degradedStep*(snake.head-index+0)))),
						255,
					)
				}

				if coordNext.y > coord.y {
					bodyColor[0] = rl.NewColor(
						uint8(47+math.Min(140, float64(degradedStep*(snake.head-index+1)))),
						uint8(78+math.Min(140, float64(degradedStep*(snake.head-index+1)))),
						uint8(0+math.Min(140, float64(degradedStep*(snake.head-index+1)))),
						255,
					)

					bodyColor[3] = rl.NewColor(
						uint8(47+math.Min(140, float64(degradedStep*(snake.head-index+1)))),
						uint8(78+math.Min(140, float64(degradedStep*(snake.head-index+1)))),
						uint8(0+math.Min(140, float64(degradedStep*(snake.head-index+1)))),
						255,
					)

					bodyColor[1] = rl.NewColor(
						uint8(47+math.Min(140, float64(degradedStep*(snake.head-index+0)))),
						uint8(78+math.Min(140, float64(degradedStep*(snake.head-index+0)))),
						uint8(0+math.Min(140, float64(degradedStep*(snake.head-index+0)))),
						255,
					)

					bodyColor[2] = rl.NewColor(
						uint8(47+math.Min(140, float64(degradedStep*(snake.head-index+0)))),
						uint8(78+math.Min(140, float64(degradedStep*(snake.head-index+0)))),
						uint8(0+math.Min(140, float64(degradedStep*(snake.head-index+0)))),
						255,
					)
				}

				if coordNext.x < coord.x {
					bodyColor[2] = rl.NewColor(
						uint8(47+math.Min(140, float64(degradedStep*(snake.head-index+1)))),
						uint8(78+math.Min(140, float64(degradedStep*(snake.head-index+1)))),
						uint8(0+math.Min(140, float64(degradedStep*(snake.head-index+1)))),
						255,
					)

					bodyColor[3] = rl.NewColor(
						uint8(47+math.Min(140, float64(degradedStep*(snake.head-index+1)))),
						uint8(78+math.Min(140, float64(degradedStep*(snake.head-index+1)))),
						uint8(0+math.Min(140, float64(degradedStep*(snake.head-index+1)))),
						255,
					)

					bodyColor[0] = rl.NewColor(
						uint8(47+math.Min(140, float64(degradedStep*(snake.head-index+0)))),
						uint8(78+math.Min(140, float64(degradedStep*(snake.head-index+0)))),
						uint8(0+math.Min(140, float64(degradedStep*(snake.head-index+0)))),
						255,
					)

					bodyColor[1] = rl.NewColor(
						uint8(47+math.Min(140, float64(degradedStep*(snake.head-index+0)))),
						uint8(78+math.Min(140, float64(degradedStep*(snake.head-index+0)))),
						uint8(0+math.Min(140, float64(degradedStep*(snake.head-index+0)))),
						255,
					)
				}

				if coordNext.y < coord.y {
					bodyColor[1] = rl.NewColor(
						uint8(47+math.Min(140, float64(degradedStep*(snake.head-index+1)))),
						uint8(78+math.Min(140, float64(degradedStep*(snake.head-index+1)))),
						uint8(0+math.Min(140, float64(degradedStep*(snake.head-index+1)))),
						255,
					)

					bodyColor[2] = rl.NewColor(
						uint8(47+math.Min(140, float64(degradedStep*(snake.head-index+1)))),
						uint8(78+math.Min(140, float64(degradedStep*(snake.head-index+1)))),
						uint8(0+math.Min(140, float64(degradedStep*(snake.head-index+1)))),
						255,
					)

					bodyColor[0] = rl.NewColor(
						uint8(47+math.Min(140, float64(degradedStep*(snake.head-index+0)))),
						uint8(78+math.Min(140, float64(degradedStep*(snake.head-index+0)))),
						uint8(0+math.Min(140, float64(degradedStep*(snake.head-index+0)))),
						255,
					)

					bodyColor[3] = rl.NewColor(
						uint8(47+math.Min(140, float64(degradedStep*(snake.head-index+0)))),
						uint8(78+math.Min(140, float64(degradedStep*(snake.head-index+0)))),
						uint8(0+math.Min(140, float64(degradedStep*(snake.head-index+0)))),
						255,
					)
				}

				applesToDigest[appleIndex] = true
				break
			}
		}

		rl.DrawRectangleGradientEx(
			rl.NewRectangle(
				float32(snake.coordinateConverter.XToPixel(coord.x)),
				float32(snake.coordinateConverter.YToPixel(coord.y)),
				float32(size),
				float32(size),
			),
			bodyColor[0],
			bodyColor[1],
			bodyColor[2],
			bodyColor[3],
		)

		if index == snake.head {
			color := rl.NewColor(
				uint8(47+degradedStep),
				uint8(78+degradedStep),
				uint8(0+degradedStep),
				255,
			)

			// 1 : 140 ds 140
			// 2 :  70 ds 71
			// 3 :  46

			colors := []rl.Color{}

			// DrawRectangleGradientV
			switch snake.direction {
			case Up:
				colors = []rl.Color{
					rl.NewColor(47, 78, 0, 255),
					color,
					color,
					rl.NewColor(47, 78, 0, 255),
				}

				/*
					if snake.length > 1 {
						coordNext := snake.getBody(index - 1)

						if coordNext.x > coord.x {
							println("a droite")
							colors[3] = rl.NewColor(157, 196, 98, 255)
						}

						if coordNext.y > coord.y {
							println("en bas")
						}

						if coordNext.x < coord.x {
							println("a gauche")
							colors[0] = rl.NewColor(157, 196, 98, 255)
						}

						if coordNext.y < coord.y {
							println("en haut")
						}
					}
				*/

				// haut gauche
				// bas gauche
				// bas droite
				// haut droite

				/*
						rl.NewColor(47, 78, 0, 255),
						rl.NewColor(157, 196, 98, 255),
						rl.NewColor(157, 196, 98, 255),
						rl.NewColor(255, 196, 98, 255),
					)
					/*
									rl.DrawRectangleGradientV(
							snake.coordinateConverter.XToPixel(coord.x),
							snake.coordinateConverter.YToPixel(coord.y),
							size,
							size,
							rl.NewColor(47, 78, 0, 255),
							rl.NewColor(157, 196, 98, 255),
						)
				*/
			case Down:
				colors = []rl.Color{
					color,
					rl.NewColor(47, 78, 0, 255),
					rl.NewColor(47, 78, 0, 255),
					color,
				}
			case Left:
				colors = []rl.Color{
					rl.NewColor(47, 78, 0, 255),
					rl.NewColor(47, 78, 0, 255),
					color,
					color,
				}
			case Right:
				colors = []rl.Color{
					color,
					color,
					rl.NewColor(47, 78, 0, 255),
					rl.NewColor(47, 78, 0, 255),
				}
			}

			rl.DrawRectangleGradientEx(
				rl.NewRectangle(
					float32(snake.coordinateConverter.XToPixel(coord.x)),
					float32(snake.coordinateConverter.YToPixel(coord.y)),
					float32(size),
					float32(size),
				),
				colors[0],
				colors[1],
				colors[2],
				colors[3],
			)
		}
	}

	for index := range snake.applesEated {
		if applesToDigest[index] == false {
			snake.applesEated = append(snake.applesEated[:index], snake.applesEated[index+1:]...)
		}
	}
}

func (snake *Snake) setBody(index int, position Position) {
	snake.body[index%len(snake.body)] = position
}

func (snake *Snake) getBody(index int) Position {
	return snake.body[index%len(snake.body)]
}
