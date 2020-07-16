package nature

import (
	"image/color"

	"github.com/tannerhat/Cactus-Simulator/game"
)

type plant struct {
	*game.Shape
	root             *roots
	water            uint32
	speed            int
	ticks            int
	waterCostPerCell uint32
}

func NewPlant(x int, y int, root *roots) *plant {
	p := &plant{
		Shape:            game.NewShape(x, y, 1, 1, color.RGBA{0x00, 0xff, 0x00, 0xff}),
		speed:            2,
		ticks:            0,
		root:             root,
		waterCostPerCell: 20,
	}
	p.Cells[0][0] = true
	return p
}

func (p *plant) Update() {
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
