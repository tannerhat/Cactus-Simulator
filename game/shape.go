package game

import (
	"image/color"

	"github.com/hajimehoshi/ebiten"
)

// Shape is a simple non physical entity. Non physical means that it does not take up any space on the gameboard but will still be drawn.
type Shape struct {
	// The gameboard X coordinate of Shape
	X int
	// The gameboard Y coordinate of Shape
	Y int
	// Cells tracks which of the cells in the box bounded by (X,Y) and (X+width,Y+height) are actually part of the Shape.
	Cells     [][]bool
	Gameboard Gameboard
	color     color.Color
	layer     int
}

// New shape returns a Shape that is located at gameboard coordinates (x,y) with an empty Cells matrix of size width x height.
func NewShape(x int, y int, width int, height int, layer int, color color.Color) *Shape {
	s := &Shape{
		X:     x,
		Y:     y,
		color: color,
		layer: layer,
	}

	s.Cells = make([][]bool, width)
	for i := range s.Cells {
		s.Cells[i] = make([]bool, height)
	}

	return s
}

// Draw the shape to screen. It will be drawn starting at (X*scale,Y*scale). only x,y coordinates where Cells[x][y] is true are drawn.
func (s *Shape) Draw(screen *ebiten.Image, scale int) {
	cellImage, _ := ebiten.NewImage(scale, scale, ebiten.FilterDefault)
	defer cellImage.Dispose()
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

// AddToBoard stores the gameboard to the Shape, it doesn't set any positions on the board because shape is non physical
func (s *Shape) AddToBoard(gameboard Gameboard) {
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

func (s *Shape) Layer() int {
	return s.layer
}
