package nature

import (
	"image/color"
	"math/rand"

	"github.com/hajimehoshi/ebiten"
)

const maxRootWetness uint32 = 3

type roots struct {
	*shape
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

func (c *roots) Name() string {
	return "roots"
}

func NewRoots(x int, y int, width int, height int, startX int, startY int) *roots {
	r := &roots{
		shape: NewShape(x, y, width, height, color.RGBA{0xff, 0xff, 0xff, 0xff}),
		rootRoot: &rootCell{
			children: make([]*rootCell, 0),
			wetness:  0,
			x:        startX,
			y:        startY,
		},
		growRate: 100,
		speed:    200,
		ticks:    0,
	}

	r.cells[startX][startY] = true

	return r
}

func (rc *rootCell) absorbFromSoil(gameboard GameBoard, rootBox *roots) {
	for _, child := range rc.children {
		child.absorbFromSoil(gameboard, rootBox)
	}

	if rc.wetness < maxRootWetness {
		boardX := rootBox.x + rc.x
		boardY := rootBox.y + rc.y
		entity := gameboard.EntityAt(boardX, boardY)
		if soil, ok := entity.(*soil); ok {
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

				boardX := rootBox.x + rc.x + dX
				boardY := rootBox.y + rc.y + dY

				entity := gameboard.EntityAt(boardX, boardY)
				if soil, ok := entity.(*soil); ok {
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
func (rc *rootCell) grow(gameboard GameBoard, rootBox *roots) bool {
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

func (r *roots) AddRoot(x int, y int, gameboard GameBoard) bool {
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

			if r.cells[x+dX][y+dY] {
				rootsFound++
			}
		}
	}

	if rootsFound > 1 {
		return false
	}

	boardX := r.x + x
	boardY := r.y + y
	entity := gameboard.EntityAt(boardX, boardY)
	if soil, ok := entity.(*soil); ok {
		// we we only grow into a wet cell
		wet, err := soil.IsWet(boardX, boardY)
		if err != nil {
			return false
		}

		if wet {
			err = soil.DigPartial(boardX, boardY)
			if err != nil {
				return false
			}
			r.cells[x][y] = true
			return true
		}
	}

	return false
}

func (r *roots) Draw(screen *ebiten.Image, scale int) {
	cellImage, _ := ebiten.NewImage(5, 5, ebiten.FilterDefault)

	// draw dry roots
	cellImage.Fill(color.RGBA{0xff, 0xff, 0xff, 0xff})
	r.rootRoot.Draw(r, screen, scale, cellImage, 0)

	// draw all wet roots the same
	cellImage.Fill(color.RGBA{0x00, 0x00, 0xff, 0xff})
	for wetness := uint32(1); wetness <= maxRootWetness; wetness++ {
		r.rootRoot.Draw(r, screen, scale, cellImage, wetness)
	}
}

func (rc *rootCell) Draw(box *roots, screen *ebiten.Image, scale int, image *ebiten.Image, wetness uint32) {
	rootWetness := rc.wetness
	if rootWetness > maxRootWetness {
		rootWetness = maxRootWetness
	}
	if rootWetness == wetness {
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(float64((box.x+rc.x)*5), float64((box.y+rc.y)*5))
		screen.DrawImage(image, op)
	}

	for _, child := range rc.children {
		child.Draw(box, screen, scale, image, wetness)
	}
}

func (r *roots) SuckWater() uint32 {
	return r.rootRoot.getWaterFromChildren()
}

func (r *roots) Update(gameBoard GameBoard) {
	r.rootRoot.grow(gameBoard, r)

	r.ticks++
	if r.ticks%r.speed == 0 {
		r.rootRoot.absorbFromSoil(gameBoard, r)
	}
}

func (r *roots) AddToBoard(gameBoard GameBoard) {
	// roots don't take up space on the board, they exist sort of on top
	// of soil. entities that interact with the cells that roots occupy
	// should treat the cells as containing soil. they must be in soil though
	// so we must check that.

	for x := range r.cells {
		for y := range r.cells[x] {
			if r.cells[x][y] {
				boardX := r.x + x
				boardY := r.y + y
				entity := gameBoard.EntityAt(boardX, boardY)
				if soil, ok := entity.(*soil); ok {
					err := soil.DigPartial(boardX, boardY)
					if err != nil {
						panic(err)
					}
				}
			}
		}
	}

	return
}
