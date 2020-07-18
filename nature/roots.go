package nature

import (
	"image/color"
	"math/rand"

	"github.com/hajimehoshi/ebiten"
	"github.com/tannerhat/Cactus-Simulator/game"
)

const maxRootWetness uint32 = 3

type Roots struct {
	*game.Shape
	rootRoot *rootCell
	growRate int
	speed    int
	ticks    int
}

// rootCell is a single cell of the root. it is not a root in the computer sceince sense
type rootCell struct {
	children []*rootCell
	wetness  uint32
	x        int
	y        int
}

func NewRoots(x int, y int, width int, height int, startX int, startY int) *Roots {
	r := &Roots{
		Shape: game.NewShape(x, y, width, height, 2, color.RGBA{0xff, 0xff, 0xff, 0xff}),
		rootRoot: &rootCell{
			children: make([]*rootCell, 0),
			wetness:  0,
			x:        startX,
			y:        startY,
		},
		growRate: 1500,
		speed:    500,
		ticks:    0,
	}

	r.Cells[startX][startY] = true

	return r
}

func (rc *rootCell) absorbFromSoil(gameboard game.Gameboard, rootBox *Roots) {
	for _, child := range rc.children {
		child.absorbFromSoil(gameboard, rootBox)
	}

	if rc.wetness < maxRootWetness {
		boardX := rootBox.X + rc.x
		boardY := rootBox.Y + rc.y
		entity := gameboard.EntityAt(boardX, boardY)
		if soil, ok := entity.(*Soil); ok {
			// we we only grow into a wet cell
			waterRemoved, err := soil.RemoveWater(boardX, boardY)
			if err != nil {
				return
			}
			if waterRemoved {
				rc.wetness++
				return
			}
		}

		// look in surrounding cells for water if none found in current cell
		for dX := -1; dX < 2; dX++ {
			for dY := -1; dY < 2; dY++ {
				if rc.x+dX < 0 || rc.x+dX >= rootBox.Width() || rc.y+dY < 0 || rc.y+dY >= rootBox.Height() {
					continue
				}

				boardX := rootBox.X + rc.x + dX
				boardY := rootBox.Y + rc.y + dY

				entity := gameboard.EntityAt(boardX, boardY)
				if soil, ok := entity.(*Soil); ok {
					// we we only grow into a wet cell
					waterRemoved, err := soil.RemoveWater(boardX, boardY)
					if err != nil {
						return
					}
					if waterRemoved {
						rc.wetness++
						return
					}
				}
			}
		}
	}
}

func (rc *rootCell) getWaterFromChildren() uint32 {
	// we can only pass up water already in the root, this must be determined before
	// getting water from children
	var waterToPassUp uint32
	if rc.wetness > 0 {
		rc.wetness--
		waterToPassUp = 1
	}

	// now get water from children only if we have room, it is possible that
	// this operation will overload us, that's fine
	if rc.wetness < maxRootWetness {
		waterFromChildren := uint32(0)
		for _, child := range rc.children {
			waterFromChildren += child.getWaterFromChildren()
		}
		rc.wetness += waterFromChildren
	}

	// pass up original amount (bottlenecked at one)
	return waterToPassUp
}

// grow tells a root cell to grow. a call to grow will result in at most 1
// new root cell. grow will add the grown cells to the linked list as well
// as to the rootBox contaning it.
func (rc *rootCell) grow(gameboard game.Gameboard, rootBox *Roots) bool {
	// we favor deep root growth, give children a chance to grow first
	for _, child := range rc.children {
		if child.grow(gameboard, rootBox) {
			return true
		}
	}

	// no children grew. try and get this cell to grow
	if rand.Intn(rootBox.growRate) == 0 {
		xDir := -1 + rand.Intn(3)
		yDir := -1 + rand.Intn(3)
		if rootBox.AddRoot(rc.x+xDir, rc.y+yDir, gameboard) {
			rc.children = append(rc.children, &rootCell{
				children: make([]*rootCell, 0),
				wetness:  0,
				x:        rc.x + xDir,
				y:        rc.y + yDir,
			})
		}
	}

	return false
}

func (r *Roots) AddRoot(x int, y int, gameboard game.Gameboard) bool {
	if x < 0 || x >= r.Width() || y < 0 || y >= r.Height() {
		return false
	}
	// we can only add root if there's only 1 root cell in the 3x3
	// surrounding area.
	rootsFound := 0
	for dX := -1; dX < 2; dX++ {
		for dY := -1; dY < 2; dY++ {
			if x+dX < 0 || x+dX >= r.Width() || y+dY < 0 || y+dY >= r.Height() {
				continue
			}

			// we want to allow a continuous

			if r.Cells[x+dX][y+dY] {
				rootsFound++
			}
		}
	}

	if rootsFound > 1 {
		return false
	}

	boardX := r.X + x
	boardY := r.Y + y
	entity := gameboard.EntityAt(boardX, boardY)
	if soil, ok := entity.(*Soil); ok {
		// we we only grow into a wet cell
		wet, err := soil.IsWet(boardX, boardY)
		if err != nil {
			return false
		}

		if wet {
			r.Cells[x][y] = true
			return true
		}
	}

	return false
}

func (r *Roots) Draw(screen *ebiten.Image, scale int) {
	cellImage, _ := ebiten.NewImage(scale, scale, ebiten.FilterDefault)
	defer cellImage.Dispose()

	// draw dry roots
	cellImage.Fill(color.RGBA{0xff, 0xff, 0xff, 0xff})
	r.rootRoot.Draw(r, screen, scale, cellImage, 0)

	// draw all wet roots the same
	cellImage.Fill(color.RGBA{0x00, 0x00, 0xff, 0xff})
	for wetness := uint32(1); wetness <= maxRootWetness; wetness++ {
		r.rootRoot.Draw(r, screen, scale, cellImage, wetness)
	}
}

func (rc *rootCell) Draw(box *Roots, screen *ebiten.Image, scale int, image *ebiten.Image, wetness uint32) {
	rootWetness := rc.wetness
	if rootWetness > maxRootWetness {
		rootWetness = maxRootWetness
	}
	if rootWetness == wetness {
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(float64((box.X+rc.x)*scale), float64((box.Y+rc.y)*scale))
		screen.DrawImage(image, op)
	}

	for _, child := range rc.children {
		child.Draw(box, screen, scale, image, wetness)
	}
}

func (r *Roots) SuckWater() uint32 {
	return r.rootRoot.getWaterFromChildren()
}

func (r *Roots) Update() {
	r.rootRoot.grow(r.Gameboard, r)

	r.ticks++
	if r.ticks%r.speed == 0 {
		r.rootRoot.absorbFromSoil(r.Gameboard, r)
	}
}

func (r *Roots) AddToBoard(gameBoard game.Gameboard) {
	r.Shape.AddToBoard(gameBoard)
	// roots don't take up space on the board, they exist sort of on top
	// of soil. entities that interact with the cells that roots occupy
	// should treat the cells as containing soil. they must be in soil though
	// so we must check that.

	for x := range r.Cells {
		for y := range r.Cells[x] {
			if r.Cells[x][y] {
				boardX := r.X + x
				boardY := r.Y + y
				entity := r.Gameboard.EntityAt(boardX, boardY)
				if _, ok := entity.(*Soil); !ok {
					panic("creating roots in non soil location, idiot")
				}
			}
		}
	}

	return
}
