package game

import (
	"fmt"
	"time"

	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/ebitenutil"
	"github.com/hajimehoshi/ebiten/inpututil"
	"github.com/paulbellamy/ratecounter"
)

// Game implements ebiten.Game and keeps track of the gameboard and entities.
type Game struct {
	gameboard    Gameboard
	drawTime     *ratecounter.AvgRateCounter
	updateTime   *ratecounter.AvgRateCounter
	debug        bool
	screenHeight int
	screenWidth  int
	scale        int
	speed        int
	ticks        int
}

// Update progresses the game one tick, updating all entities that have been added to the game's board.
func (g *Game) Update(screen *ebiten.Image) error {
	updateStart := time.Now()

	if inpututil.IsKeyJustPressed(ebiten.Key1) {
		g.speed = 0
	}
	if inpututil.IsKeyJustPressed(ebiten.Key2) {
		g.speed = 1
	}
	if inpututil.IsKeyJustPressed(ebiten.Key3) {
		g.speed = 10
	}
	if inpututil.IsKeyJustPressed(ebiten.Key4) {
		g.speed = 60
	}
	if inpututil.IsKeyJustPressed(ebiten.Key5) {
		g.speed = 300
	}

	for i := 0; i < g.speed; i++ {
		g.ticks++
		if inpututil.IsKeyJustPressed(ebiten.KeyD) {
			g.debug = !g.debug
		}

		entityChan := g.gameboard.Entities()

		for e := range entityChan {
			e.Update()
		}

	}
	g.updateTime.Incr(int64(time.Since(updateStart)))
	return nil
}

// Draw writes the screen image to the given ebiten.Image. All entities in the gameboard are given the chance to draw.
// draw order is not guaranteed.
func (g *Game) Draw(screen *ebiten.Image) {
	drawsStart := time.Now()

	entityChan := g.gameboard.Entities()
	entityList := []Entity{}
	maxLayer := 0
	for e := range entityChan {
		entityList = append(entityList, e)
		if e.Layer() > maxLayer {
			maxLayer = e.Layer()
		}
	}

	// reuse the entity list rather than get a new entities channel from gameboard because entities could've changed
	for layer := 0; layer <= maxLayer; layer++ {
		for _, e := range entityList {
			if e.Layer() == layer {
				e.Draw(screen, g.scale)
			}
		}
	}

	g.drawTime.Incr(int64(time.Since(drawsStart)))

	if g.debug {
		msg := fmt.Sprintf(`FPS: %0.2f
Draw Time: %0.2f ms
Update Time: %0.2f ms
Speed: %d
Game time: %0.2f hours`,
			ebiten.CurrentFPS(),
			g.drawTime.Rate()/float64(time.Millisecond),
			g.updateTime.Rate()/float64(time.Millisecond),
			g.speed,
			float64(g.ticks)/60.0/60.0/60.0)
		ebitenutil.DebugPrint(screen, msg)
	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return g.screenWidth, g.screenHeight
}

// AddEntity adds the given entity to the game's board.
func (g *Game) AddEntity(entity Entity) {
	g.gameboard.AddEntity(entity)
}

// NewGame creates a game with the given screen width and height. Scale indicates how many pixels per cell in the gameboard.
func NewGame(width int, height int, scale int) *Game {
	g := Game{
		drawTime:     ratecounter.NewAvgRateCounter(time.Second),
		updateTime:   ratecounter.NewAvgRateCounter(time.Second),
		debug:        true,
		screenWidth:  width,
		screenHeight: height,
		scale:        scale,
		speed:        1,
		ticks:        0,
	}
	b := NewGameboard(g.screenWidth/g.scale, g.screenHeight/g.scale)
	g.gameboard = b

	return &g
}
