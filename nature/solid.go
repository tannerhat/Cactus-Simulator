package nature

import (
	"image/color"
)

type solid struct {
	*shape
}

func NewSolid(x int, y int, width int, height int, color color.Color) *solid {
	s := &solid{
		shape: NewShape(x, y, width, height, color),
	}

	return s
}

func (s *solid) AddToBoard(gameboard GameBoard) {
	for x := range s.cells {
		for y := range s.cells[x] {
			if s.cells[x][y] {
				gameboard.SetEntity(s, s.x+x, s.y+y)
			}
		}
	}
}
