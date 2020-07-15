package nature

import (
	"image/color"
	"math/rand"

	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/inpututil"
	"github.com/tannerhat/Cactus-Simulator/game"
)

type cloud struct {
	*game.Solid
	rate       int
	ticks      int
	raining    bool
	waterImage *ebiten.Image
}

func (c *cloud) Name() string {
	return "cloud"
}

func NewCloud(x int, y int, width int, height int, rate int) *cloud {
	c := &cloud{
		Solid:   game.NewSolid(x, y, width, height, color.RGBA{0xff, 0xff, 0xff, 0xff}),
		rate:    rate,
		ticks:   0,
		raining: false,
	}

	c.waterImage, _ = ebiten.NewImage(5, 5, ebiten.FilterDefault)

	c.waterImage.Fill(color.RGBA{
		0x00,
		0x00,
		0xff,
		0xff,
	})

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

func (c *cloud) Update(gameboard game.GameBoard) {
	if inpututil.IsKeyJustPressed(ebiten.KeySpace) {
		c.raining = !c.raining
	}
	if c.raining {
		if c.ticks%c.rate == 0 {
			gameboard.AddEntity(
				&water{
					x:       rand.Intn(c.Width()-2) + c.X + 1, // because the edges are rounded
					y:       c.Y + c.Height(),
					density: 1,
					settled: 0,
					image:   c.waterImage,
				},
			)
		}
		c.ticks++
	}
}
