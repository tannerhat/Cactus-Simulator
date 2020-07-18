package nature

import (
	"image/color"
	"math/rand"

	"github.com/tannerhat/Cactus-Simulator/game"
)

// Cloud is a Shape that creates water entities that fall starting from the cloud's lower edge. Raining can be toggled by spacebar.
type Cloud struct {
	*game.Shape
	rate    int
	ticks   int
	raining bool
}

// NewCloud returns a cloud that will be at gameboard coordinates (x,y) once added to the game. Rate indicates
// how many ticks between each water entity creation.
func NewCloud(x int, y int, width int, height int, rate int) *Cloud {
	c := &Cloud{
		Shape:   game.NewShape(x, y, width, height, 0, color.RGBA{0xff, 0xff, 0xff, 0xff}),
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
	if c.raining {
		if c.ticks%c.rate == 0 {
			c.Gameboard.AddEntity(
				&Water{
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

func (c *Cloud) SetStatus(raining bool, rate int) {
	c.raining = raining
	c.rate = rate
}
