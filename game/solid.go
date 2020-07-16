package game

import (
	"image/color"
)

// Solid is an extension of Shape with a physical presence on the board
type Solid struct {
	*Shape
}

func NewSolid(x int, y int, width int, height int, color color.Color) *Solid {
	s := &Solid{
		Shape: NewShape(x, y, width, height, color),
	}

	return s
}

// AddToBoard calls adds the Solid to any gameboard locations indicated by the Cells matrix
func (s *Solid) AddToBoard(gameboard Gameboard) {
	s.Shape.AddToBoard(gameboard)
	for x := range s.Cells {
		for y := range s.Cells[x] {
			if s.Cells[x][y] {
				s.Gameboard.SetEntity(s, s.X+x, s.Y+y)
			}
		}
	}
}
