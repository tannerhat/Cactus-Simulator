package game

import (
	"image/color"

	"github.com/hajimehoshi/ebiten"
)

type Shape struct {
	X         int
	Y         int
	Cells     [][]bool
	Gameboard GameBoard
	color     color.Color
}

func (s *Shape) Name() string {
	return "shape"
}

func NewShape(x int, y int, width int, height int, color color.Color) *Shape {
	s := &Shape{
		X:     x,
		Y:     y,
		color: color,
	}

	s.Cells = make([][]bool, width)
	for i := range s.Cells {
		s.Cells[i] = make([]bool, height)
	}

	return s
}

func (s *Shape) Draw(screen *ebiten.Image, scale int) {
	cellImage, _ := ebiten.NewImage(scale, scale, ebiten.FilterDefault)
	cellImage.Fill(s.color)

	for x := range s.Cells {
		for y := range s.Cells[x] {
			if s.Cells[x][y] {
				op := &ebiten.DrawImageOptions{}
				op.GeoM.Translate(float64((s.X+x)*scale), float64((s.Y+y)*scale))
				screen.DrawImage(cellImage, op)
			}
		}
	}
}

func (s *Shape) Update() {
	return
}

func (s *Shape) AddToBoard(gameboard GameBoard) {
	s.Gameboard = gameboard
	return
}

func (s *Shape) Width() int {
	return len(s.Cells)
}

func (s *Shape) Height() int {
	// TODO: enforce a width > 0 and height > 0 requirement on solid
	return len(s.Cells[0])
}
