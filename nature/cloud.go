package nature

import (
	"image/color"
	"math/rand"

	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/inpututil"
	"github.com/tannerhat/Cactus-Simulator/game"
)

// Cloud is a Solid that creates water entities that fall starting from the cloud's lower edge. Raining can be toggled by spacebar.
type Cloud struct {
	*game.Solid
	rate    int
	ticks   int
	raining bool
}

// NewCloud returns a cloud that will be at gameboard coordinates (x,y) once added to the game. Rate indicates
// how many ticks between each water entity creation.
func NewCloud(x int, y int, width int, height int, rate int) *Cloud {
	c := &Cloud{
		Solid:   game.NewSolid(x, y, width, height, color.RGBA{0xff, 0xff, 0xff, 0xff}),
		rate:    rate,
		ticks:   0,
		raining: false,
	}

	for x := range c.Cells {
		for y := range c.Cells[x] {
			xEdge := (x == 0 || x == width-1)
			yEdge := (y == 0 || y == height-1)
			if !xEdge || !yEdge {
				c.Cells[x][y] = true
			}
		}
	}

	return c
}

// Update the cloud, if it causes a water entity to be created, the cloud will add the entity to the gameboard directly.
func (c *Cloud) Update() {
	if inpututil.IsKeyJustPressed(ebiten.KeySpace) {
		c.raining = !c.raining
	}
	if c.raining {
		if c.ticks%c.rate == 0 {
			c.Gameboard.AddEntity(
				&water{
					x:       rand.Intn(c.Width()-2) + c.X + 1, // because the edges are rounded
					y:       c.Y + c.Height(),
					density: 1,
					settled: 0,
				},
			)
		}
		c.ticks++
	}
}
