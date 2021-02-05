package game

type CoordinateConverter struct {
	borderX int32
	borderY int32
	grid    int32
}

func NewCoordinateConverter(borderX int32, borderY int32, grid int32) CoordinateConverter {
	return CoordinateConverter{
		borderX: borderX,
		borderY: borderY,
		grid:    grid,
	}
}

func (position *CoordinateConverter) XToPixel(x int32) int32 {
	return position.borderX + x*position.grid
}

func (position *CoordinateConverter) YToPixel(y int32) int32 {
	return position.borderY + y*position.grid
}
