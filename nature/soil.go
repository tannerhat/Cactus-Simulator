package nature

import (
	"fmt"
	"image/color"
	"math/rand"
	"sync"

	"github.com/hajimehoshi/ebiten"
	"github.com/tannerhat/Cactus-Simulator/game"
)

type Soil struct {
	*game.Solid
	wetness       [][]uint32
	absorbRate    int
	evaporateRate int
	colors        []color.Color
}

const maxWetness uint32 = 3

func NewSoil(x int, y int, width int, height int) *Soil {
	s := &Soil{
		Solid:         game.NewSolid(x, y, width, height, 1, color.RGBA{0xc2, 0xb2, 0x80, 0xff}),
		absorbRate:    3,
		evaporateRate: 200,
	}

	s.colors = make([]color.Color, maxWetness+1)
	for wetness := range s.colors {
		r, g, b, a := color.RGBA{0xc2, 0xb2, 0x80, 0xff}.RGBA()
		r &= 0xff
		g &= 0xff
		b &= 0xff
		a &= 0xff
		// max wetness / 2 prevents the soil from being too dark
		s.colors[wetness] = color.RGBA{
			uint8(((maxWetness + maxWetness/2) - uint32(wetness)) * r / (maxWetness + maxWetness/2)),
			uint8(((maxWetness + maxWetness/2) - uint32(wetness)) * g / (maxWetness + maxWetness/2)),
			uint8(((maxWetness + maxWetness/2) - uint32(wetness)) * b / (maxWetness + maxWetness/2)),
			uint8(a),
		}

	}

	for x := range s.Cells {
		for y := range s.Cells[x] {
			// soil is fully solid
			s.Cells[x][y] = true
		}
	}

	s.wetness = make([][]uint32, width)
	for i := range s.wetness {
		s.wetness[i] = make([]uint32, height)
	}

	return s
}

func (s *Soil) Draw(screen *ebiten.Image, scale int) {
	cellImage, _ := ebiten.NewImage(scale, scale, ebiten.FilterDefault)
	defer cellImage.Dispose()

	var wetness uint32
	for ; wetness <= maxWetness; wetness++ {
		c := s.getColor(wetness)
		cellImage.Fill(c)

		for x := range s.Cells {
			for y := range s.Cells[x] {
				if s.Cells[x][y] {
					cellWetness := s.wetness[x][y]
					if cellWetness > maxWetness {
						cellWetness = maxWetness
					}
					if cellWetness == wetness {
						// scale color by wetness
						op := &ebiten.DrawImageOptions{}
						op.GeoM.Translate(float64((s.X+x)*scale), float64((s.Y+y)*scale))
						screen.DrawImage(cellImage, op)
					}
				}
			}
		}
	}

}

// getColor takes the soil coordinates of a cell and returns the color to display the cell as
func (s *Soil) getColor(wetness uint32) color.Color {
	if wetness > maxWetness {
		wetness = maxWetness
	}
	return s.colors[wetness]
}

func (s *Soil) updateSubGroup(group int, subGroup int, wg *sync.WaitGroup) {

}

func (s *Soil) Update() {
	directions := [][]int{
		[]int{0, 1},
		[]int{0, -1},
		[]int{1, 0},
		[]int{-1, 0},
		[]int{-1, -1},
		[]int{-1, 1},
		[]int{1, 1},
		[]int{1, -1},
	}

	for x := 0; x < s.Width(); x++ {
		for y := 0; y < s.Height(); y++ {
			if (s.wetness[x][y] == 1 || (s.wetness[x][y] > 1 && y == 0)) && rand.Intn((y/2+1)*s.evaporateRate) == 0 {
				s.wetness[x][y]--
			}
			if s.wetness[x][y] > 1 {
				for _, modifier := range directions {
					if rand.Intn(s.absorbRate) == 0 {
						otherX := x + modifier[0]
						otherY := y + modifier[1]
						if otherX >= 0 && otherX < s.Width() &&
							otherY >= 0 && otherY < s.Height() {
							if (s.wetness[x][y]-1 > s.wetness[otherX][otherY]) ||
								(s.wetness[x][y] > maxWetness) {
								s.wetness[otherX][otherY]++
								s.wetness[x][y]--
							}
						} else if s.wetness[x][y] > maxWetness {
							// otherX/otherY is off the screen. transfer wetness if we are > max to prevent
							// soil oversaturation
							s.wetness[x][y]--
						}
					}
				}
			}
		}
	}
}

// Absorb takes the gameboard coordinates of a soil cell and returns true if that cell successfully absorbs
func (s *Soil) Absorb(x int, y int) bool {
	// convert x and y into soil position
	x -= s.X
	y -= s.Y

	// we can absorb upt to max wetness + 1. an oversaturated soil cell is capable
	// of sending its extra water to neighbors. If we don't allow oversaturation,
	// our absorbtion algorithm doesn't give a way to get non topsoil cells to reach
	// maxWetness
	if s.wetness[x][y] < (maxWetness+1) && rand.Intn(s.absorbRate) == 0 {
		s.wetness[x][y]++
		return true
	}

	return false
}

// TODO implement absorber interface so soil can use the solid AddToBoard, if we use soil add to board,
// the roots/water won't be able to tell.
func (s *Soil) AddToBoard(gameboard game.Gameboard) {
	s.Solid.AddToBoard(gameboard)
	for x := range s.Cells {
		for y := range s.Cells[x] {
			if s.Cells[x][y] {
				gameboard.SetEntity(s, s.X+x, s.Y+y)
			}
		}
	}
}

func (s *Soil) IsWet(x int, y int) (bool, error) {
	// convert x and y into soil position
	x -= s.X
	y -= s.Y

	if x < 0 || y < 0 || x >= s.Width() || y >= s.Height() {
		return false, fmt.Errorf("cell being checked (%d,%d)(%d,%d) is not in the soil", x+s.X, y+s.Y, x, y)
	}

	return s.wetness[x][y] > 0, nil
}

func (s *Soil) RemoveWater(x int, y int) (bool, error) {
	// convert x and y into soil position
	x -= s.X
	y -= s.Y

	if x < 0 || y < 0 || x >= s.Width() || y >= s.Height() {
		return false, fmt.Errorf("cell being checked (%d,%d)(%d,%d) is not in the soil", x+s.X, y+s.Y, x, y)
	}

	if s.wetness[x][y] > 0 {
		s.wetness[x][y] = s.wetness[x][y] - 1
		return true, nil
	}
	return false, nil
}
