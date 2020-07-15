package nature

import (
	"fmt"
	"image/color"
	"math/rand"
	"sync"

	"github.com/hajimehoshi/ebiten"
)

type soil struct {
	*solid
	wetness       [][]uint32
	absorbRate    int
	evaporateRate int
	colors        []color.Color
}

const maxWetness uint32 = 3

func (s *soil) Name() string {
	return "soil"
}

func NewSoil(x int, y int, width int, height int) *soil {
	s := &soil{
		solid:         NewSolid(x, y, width, height, color.RGBA{0xc2, 0xb2, 0x80, 0xff}),
		absorbRate:    3,
		evaporateRate: 300,
	}

	s.colors = make([]color.Color, maxWetness+1)
	for wetness := range s.colors {
		r, g, b, a := s.color.RGBA()
		r &= 0xff
		g &= 0xff
		b &= 0xff
		a &= 0xff
		// max wetness / 5 prevents the soil from being too dark
		s.colors[wetness] = color.RGBA{
			uint8(((maxWetness + maxWetness/2) - uint32(wetness)) * r / (maxWetness + maxWetness/2)),
			uint8(((maxWetness + maxWetness/2) - uint32(wetness)) * g / (maxWetness + maxWetness/2)),
			uint8(((maxWetness + maxWetness/2) - uint32(wetness)) * b / (maxWetness + maxWetness/2)),
			uint8(a),
		}

	}

	for x := range s.cells {
		for y := range s.cells[x] {
			// soil is fully solid
			s.cells[x][y] = true
		}
	}

	s.wetness = make([][]uint32, width)
	for i := range s.wetness {
		s.wetness[i] = make([]uint32, height)
	}

	return s
}

func (s *soil) Draw(screen *ebiten.Image, scale int) {
	cellImage, _ := ebiten.NewImage(5, 5, ebiten.FilterDefault)

	var wetness uint32
	for ; wetness <= maxWetness; wetness++ {
		c := s.getColor(wetness)
		cellImage.Fill(c)

		for x := range s.cells {
			for y := range s.cells[x] {
				if s.cells[x][y] {
					cellWetness := s.wetness[x][y]
					if cellWetness > maxWetness {
						cellWetness = maxWetness
					}
					if cellWetness == wetness {
						// scale color by wetness
						op := &ebiten.DrawImageOptions{}
						op.GeoM.Translate(float64((s.x+x)*5), float64((s.y+y)*5))
						screen.DrawImage(cellImage, op)
					}
				}
			}
		}
	}

}

// getColor takes the soil coordinates of a cell and returns the color to display the cell as
func (s *soil) getColor(wetness uint32) color.Color {
	if wetness > maxWetness {
		wetness = maxWetness
	}
	return s.colors[wetness]
}

func (s *soil) updateSubGroup(group int, subGroup int, wg *sync.WaitGroup) {
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

	subGroupWidth := s.Width() / 4
	subGroupStart := subGroupWidth*group + 2*subGroupWidth*subGroup
	subGroupEnd := subGroupStart + subGroupWidth
	if group == 1 && subGroup == 1 {
		subGroupEnd = s.Width()
	}
	for x := subGroupStart; x < subGroupEnd; x++ {
		for y := range s.cells[x] {
			if (s.wetness[x][y] == 1 || (s.wetness[x][y] > 1 && y == 0)) && rand.Intn((y+1)*s.evaporateRate) == 0 {
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
						}
					}
				}
			}
		}
	}
	wg.Done()
}

func (s *soil) Update(gameBoard GameBoard) {

	for group := 0; group < 2; group++ {
		wg := &sync.WaitGroup{}

		for subGroup := 0; subGroup < 2; subGroup++ {
			wg.Add(1)

			go s.updateSubGroup(group, subGroup, wg)
		}

		wg.Wait()
	}
}

// Absorb takes the gameboard coordinates of a soil cell and returns true if that cell successfully absorbs
func (s *soil) Absorb(x int, y int) bool {
	// convert x and y into soil position
	x -= s.x
	y -= s.y

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

// TODO implement absorber interface so soil can use the solid AddToBoard
func (s *soil) AddToBoard(gameboard GameBoard) {
	for x := range s.cells {
		for y := range s.cells[x] {
			if s.cells[x][y] {
				gameboard.SetEntity(s, s.x+x, s.y+y)
			}
		}
	}
}

// DigPartial takes gameboard coordinates of a soil cell and "makes room" for an
// entity to exist at those coordinates. Partial means that the change is cosmetic
// the soil will no longer draw the cell, otherwise no change.
func (s *soil) DigPartial(x int, y int) error {
	// convert x and y into soil position
	x -= s.x
	y -= s.y

	if x < 0 || y < 0 || x >= s.Width() || y >= s.Height() {
		return fmt.Errorf("cell being dug (%d,%d)(%d,%d) is not in the soil", x+s.x, y+s.y, x, y)
	}

	// don't display as part of the shape
	s.cells[x][y] = false

	return nil
}

func (s *soil) IsWet(x int, y int) (bool, error) {
	// convert x and y into soil position
	x -= s.x
	y -= s.y

	if x < 0 || y < 0 || x >= s.Width() || y >= s.Height() {
		return false, fmt.Errorf("cell being checked (%d,%d)(%d,%d) is not in the soil", x+s.x, y+s.y, x, y)
	}

	return s.wetness[x][y] > 0, nil
}

func (s *soil) RemoveWater(x int, y int) (bool, error) {
	// convert x and y into soil position
	x -= s.x
	y -= s.y

	if x < 0 || y < 0 || x >= s.Width() || y >= s.Height() {
		return false, fmt.Errorf("cell being checked (%d,%d)(%d,%d) is not in the soil", x+s.x, y+s.y, x, y)
	}

	if s.wetness[x][y] > 0 {
		s.wetness[x][y] = s.wetness[x][y] - 1
		return true, nil
	}
	return false, nil
}
