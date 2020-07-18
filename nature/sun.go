package nature

import (
	"image/color"

	"github.com/hajimehoshi/ebiten"
	"github.com/tannerhat/Cactus-Simulator/game"
)

type Sun struct {
	*game.Shape
	Hidden bool
}

func NewSun(x int, y int, width int, height int, layer int, color color.Color) *Sun {
	s := &Sun{
		Shape:  game.NewShape(x, y, width, height, layer, color),
		Hidden: false,
	}
	return s
}

func (s *Sun) Draw(screen *ebiten.Image, scale int) {
	if !s.Hidden {
		s.Shape.Draw(screen, scale)
	}
}
