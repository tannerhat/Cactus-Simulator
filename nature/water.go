package nature

import (
	"image/color"
	"math/rand"

	"github.com/hajimehoshi/ebiten"
	"github.com/tannerhat/Cactus-Simulator/game"
)

// all water shares the same image. This is so that when drawing them, the draw operations can be combined into one big draw operation.
var waterImage *ebiten.Image

const maxDensity = 300

type Water struct {
	x         int
	y         int
	density   int
	settled   int
	gameboard game.Gameboard
}

func (c *Water) AddToBoard(gameBoard game.Gameboard) {
	c.gameboard = gameBoard
	c.gameboard.SetEntity(c, c.x, c.y)
}

func (w *Water) Draw(screen *ebiten.Image, scale int) {
	if waterImage == nil {
		// the first water to be drawn creates the image that all water will use
		// to draw itself
		waterImage, _ = ebiten.NewImage(scale, scale, ebiten.FilterDefault)

		waterImage.Fill(color.RGBA{
			0x00,
			0x00,
			0xff,
			0xff,
		})
	}

	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(float64(w.x*scale), float64(w.y*scale))
	screen.DrawImage(waterImage, op)
}

func (c *Water) flowTo(gameBoard game.Gameboard, x int, y int, force bool, dry bool) bool {
	width, height := gameBoard.Size()
	if x < 0 || x >= width {
		// off of left or right
		return false
	}

	if y < 0 || y >= height {
		// off of top or bottom
		return false
	}

	e := gameBoard.EntityAt(x, y)
	if e == nil {
		if dry {
			return true
		}
		// the spot is empty let's go!
		if c.density == 1 {
			// we are only 1 drop, just go there bab
			gameBoard.MoveEntity(c.x, c.y, x, y)
			c.x = x
			c.y = y
		} else {
			// we have more than one density, create a drop
			// in the flow to position
			gameBoard.AddEntity(&Water{
				x:       x,
				y:       y,
				density: 1,
			})
			c.density--
		}
		return true
	}

	if other, ok := e.(*Water); ok {
		if force && other.density == 1 && !other.underPressure(gameBoard) {
			if dry {
				return true
			}
			other.density++
			c.density--
			if c.density == 0 {
				gameBoard.SetEntity(nil, c.x, c.y)
				gameBoard.RemoveEntity(c)
			}
			return true
		}
	} else if other, ok := e.(*Soil); ok {
		if dry {
			return false
		}
		if other.Absorb(x, y) {
			c.density--
			if c.density == 0 {
				gameBoard.SetEntity(nil, c.x, c.y)
				gameBoard.RemoveEntity(c)
			}
			return true
		}
	}

	return false
}

func (c *Water) underPressure(gameBoard game.Gameboard) bool {
	return !(c.flowTo(gameBoard, c.x, c.y+1, false, true) ||
		c.flowTo(gameBoard, c.x-1, c.y, false, true) ||
		c.flowTo(gameBoard, c.x+1, c.y, false, true) ||
		c.flowTo(gameBoard, c.x, c.y-1, false, true))
}

func (c *Water) Update() {
	// try to flow down
	if c.flowTo(c.gameboard, c.x, c.y+1, true, false) {
		return
	}

	firstDir := -1 + 2*rand.Intn(2)
	// we couldn't go down, try flowing first dir
	if c.flowTo(c.gameboard, c.x+firstDir, c.y, false, false) {
		if c.density == 1 {
			// we flowed left and are now single density, flowing in another
			// direction will create a gap
			return
		}
	}

	firstDir = firstDir * -1 // opposite of first dir
	// okay now try other dir
	if c.flowTo(c.gameboard, c.x+firstDir, c.y, false, false) {
		if c.density == 1 {
			// we flowed right and are now single density, flowing in another
			// direction will create a gap
			return
		}
	}

	if c.density > 1 {
		// okay fine, if we are multi density, try flowing up
		if c.flowTo(c.gameboard, c.x, c.y-1, true, false) {
			if c.density == 1 {
				// we flowed left and are now single density, flowing in another
				// direction will create a gap
				return
			}

		}

		// we couldn't go down, try flowing firstDir
		if c.flowTo(c.gameboard, c.x+firstDir, c.y, true, false) {
			if c.density == 1 {
				// we flowed left and are now single density, flowing in another
				// direction will create a gap
				return
			}
		}

		firstDir = firstDir * -1 // opposite of first dir
		// okay now try other dir
		if c.flowTo(c.gameboard, c.x+firstDir, c.y, true, false) {
			if c.density == 1 {
				// we flowed right and are now single density, flowing in another
				// direction will create a gap
				return
			}
		}
	}
}

func (w *Water) Layer() int {
	return 2
}
