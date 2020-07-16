package game

import (
	"image/color"
)

type Solid struct {
	*Shape
}

func NewSolid(x int, y int, width int, height int, color color.Color) *Solid {
	s := &Solid{
		Shape: NewShape(x, y, width, height, color),
	}

	return s
}

func (s *Solid) AddToBoard(gameboard GameBoard) {
	s.Shape.AddToBoard(gameboard)
	for x := range s.Cells {
		for y := range s.Cells[x] {
			if s.Cells[x][y] {
				s.Gameboard.SetEntity(s, s.X+x, s.Y+y)
			}
		}
	}
}
