package nature

import (
	"image/color"

	"github.com/tannerhat/Cactus-Simulator/game"
)

// Plant is a Shape that takes in water from a root entity and grows bigger from the water.
type Plant struct {
	*game.Shape
	root             *Roots
	water            uint32
	speed            int
	ticks            int
	waterCostPerCell uint32
}

// NewPlant creates a plant that will start as a 1x1 Shape at x,y. It will SuckWater from root.
func NewPlant(x int, y int, root *Roots) *Plant {
	p := &Plant{
		Shape:            game.NewShape(x, y, 1, 1, 1, color.RGBA{0x00, 0xff, 0x00, 0xff}),
		speed:            2,
		ticks:            0,
		root:             root,
		waterCostPerCell: 20,
	}
	p.Cells[0][0] = true
	return p
}

// Update will take in water from roots once every "speed" ticks. Once it gets enough water to grow, it will expand it's shape.
func (p *Plant) Update() {
	p.ticks++
	if p.ticks%p.speed == 0 {
		p.water += p.root.SuckWater()
	}

	// if the cactus has gotten out of ratio, it gets wider
	if p.Width() < p.Height()/3 {
		// growing wider means adding Height cells
		waterCost := uint32(p.Height()) * p.waterCostPerCell
		if p.water >= waterCost {
			p.Cells = append(p.Cells, make([]bool, p.Height()))
			for y := range p.Cells[p.Width()-1] {
				p.Cells[p.Width()-1][y] = true
			}

			if p.Width()%2 != 0 {
				p.X--
			}
			p.water -= waterCost
		}
	} else {
		// just get taller (add Width cells)
		waterCost := uint32(p.Width()) * p.waterCostPerCell
		if p.water >= waterCost {
			for x := range p.Cells {
				p.Cells[x] = append(p.Cells[x], true)
			}
			p.Y--
			p.water -= waterCost
		}
	}

	return
}
