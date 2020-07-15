package nature

import (
	"image/color"

	"github.com/hajimehoshi/ebiten"
)

type shape struct {
	x     int
	y     int
	cells [][]bool
	color color.Color
}

func (s *shape) Name() string {
	return "shape"
}

func NewShape(x int, y int, width int, height int, color color.Color) *shape {
	s := &shape{
		x:     x,
		y:     y,
		color: color,
	}

	s.cells = make([][]bool, width)
	for i := range s.cells {
		s.cells[i] = make([]bool, height)
	}

	return s
}

func (s *shape) Draw(screen *ebiten.Image, scale int) {
	cellImage, _ := ebiten.NewImage(5, 5, ebiten.FilterDefault)
	cellImage.Fill(s.color)

	for x := range s.cells {
		for y := range s.cells[x] {
			if s.cells[x][y] {
				op := &ebiten.DrawImageOptions{}
				op.GeoM.Translate(float64((s.x+x)*5), float64((s.y+y)*5))
				screen.DrawImage(cellImage, op)
			}
		}
	}
}

func (s *shape) Update(gameboard GameBoard) {
	return
}

func (s *shape) AddToBoard(gameboard GameBoard) {
	return
}

func (s *shape) Width() int {
	return len(s.cells)
}

func (s *shape) Height() int {
	// TODO: enforce a width > 0 and height > 0 requirement on solid
	return len(s.cells[0])
}
