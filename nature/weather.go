package nature

import (
	"image/color"
	"math/rand"

	"github.com/hajimehoshi/ebiten"
	"github.com/tannerhat/Cactus-Simulator/game"
)

const maxCloudDarkness = 5

type Weather struct {
	gameboard     game.Gameboard
	clouds        []*Cloud
	cloudSpawn    int
	skyImage      *ebiten.Image
	skyColor      color.Color
	sun           *Sun
	raining       bool
	rainStart     int
	rainStop      int
	rainIntensity int
}

func NewWeather(cloudSpawn int) *Weather {
	w := &Weather{
		clouds:        make([]*Cloud, 0),
		cloudSpawn:    cloudSpawn,
		skyColor:      color.RGBA{0x87, 0xce, 0xfa, 0xff},
		raining:       false,
		rainStart:     20000,
		rainStop:      3000,
		rainIntensity: 2,
	}

	return w
}

// Draw the entity to the given screen at the given scale
func (w *Weather) Draw(screen *ebiten.Image, scale int) {
	boardWidth, boardHeight := w.gameboard.Size()

	if w.skyImage == nil {
		w.skyImage, _ = ebiten.NewImage(boardWidth*scale, boardHeight*scale, ebiten.FilterDefault)
		w.skyImage.Fill(w.skyColor)
	}

	if w.sun == nil {
		w.sun = NewSun(2*boardWidth/3, boardHeight/15+rand.Intn(boardHeight/15), boardWidth/15, boardWidth/15, 0, color.RGBA{0xff, 0xde, 0x00, 0xff})

		for x := range w.sun.Cells {
			for y := range w.sun.Cells[x] {
				xEdge := (x == 0 || x == w.sun.Width()-1)
				yEdge := (y == 0 || y == w.sun.Height()-1)
				if !xEdge || !yEdge {
					w.sun.Cells[x][y] = true
				}
			}
		}
		w.gameboard.AddEntity(w.sun)
	}

	screen.DrawImage(w.skyImage, nil)
}

func (w *Weather) recalculateSky() {
	cloudCount := len(w.clouds)
	if cloudCount > maxCloudDarkness {
		cloudCount = maxCloudDarkness
	}
	if !w.raining {
		// only darken sky if raining
		cloudCount = 0
	}

	r, g, b, a := w.skyColor.RGBA()
	r &= 0xff
	g &= 0xff
	b &= 0xff
	a &= 0xff
	// max maxCloudDarkness / 2 prevents the sky from being too dark

	newColor := color.RGBA{
		uint8(((maxCloudDarkness + maxCloudDarkness) - uint32(cloudCount)) * r / (maxCloudDarkness + maxCloudDarkness)),
		uint8(((maxCloudDarkness + maxCloudDarkness) - uint32(cloudCount)) * g / (maxCloudDarkness + maxCloudDarkness)),
		uint8(((maxCloudDarkness + maxCloudDarkness) - uint32(cloudCount)) * b / (maxCloudDarkness + maxCloudDarkness)),
		uint8(a),
	}

	if w.skyImage != nil {
		w.skyImage.Fill(newColor)
	}

	if w.sun != nil {
		if cloudCount > 1 {
			w.sun.Hidden = true
		} else {
			w.sun.Hidden = false
		}
	}
}

// Update the entity by one tick
func (w *Weather) Update() {
	boardWidth, boardHeight := w.gameboard.Size()

	cloudCount := len(w.clouds)
	if cloudCount > maxCloudDarkness {
		cloudCount = maxCloudDarkness
	}

	if rand.Intn(w.cloudSpawn/(2*cloudCount+1)) == 0 {
		cloudWidth := boardWidth / 8
		c := NewCloud(0, boardHeight/15+rand.Intn(boardHeight/15), cloudWidth, 2*cloudWidth/3, 1)
		w.clouds = append(w.clouds, c)
		w.gameboard.AddEntity(c)
		c.SetStatus(w.raining, w.rainIntensity)
		w.recalculateSky()
	}

	for i := 0; i < len(w.clouds); {
		c := w.clouds[i]

		effectiveMoveRate := w.cloudSpawn / boardWidth
		if !w.raining {
			effectiveMoveRate /= 4
		}
		if rand.Intn(1+effectiveMoveRate) == 0 {
			c.X++
		}
		if c.X+c.Width() >= boardWidth {
			w.gameboard.RemoveEntity(c)
			w.clouds = append(w.clouds[:i], w.clouds[i+1:]...)
			i = 0 // removed an element, start over
			w.recalculateSky()
		} else {
			i++
		}
	}

	if len(w.clouds) > 0 {
		// there are clouds, determine if we should be raining
		if w.raining && rand.Intn(w.rainStop) == 0 {
			w.toggleRain(false)
		} else if !w.raining && len(w.clouds) > 0 && rand.Intn(w.rainStart/len(w.clouds)) == 0 {
			w.toggleRain(true)
		}
	}
	return
}

func (w *Weather) toggleRain(enable bool) {
	w.raining = enable
	for _, c := range w.clouds {
		c.SetStatus(enable, w.rainIntensity)
	}
	w.recalculateSky()
}

// AddToBoard is called by game when an entity is added, the entity marks its initial positions on the game board
func (w *Weather) AddToBoard(gameboard game.Gameboard) {
	w.gameboard = gameboard
}

// Layer returns the layer of the entity for draw purposes
func (w *Weather) Layer() int {
	return 0
}
