package game

import (
	"fmt"
	"image/color"
	"time"

	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/ebitenutil"
	"github.com/hajimehoshi/ebiten/inpututil"
	"github.com/paulbellamy/ratecounter"
)

// Game implements ebiten.Game and keeps track of the gameboard and entities.
type Game struct {
	canvasImage  *ebiten.Image
	gameboard    Gameboard
	drawTime     *ratecounter.AvgRateCounter
	updateTime   *ratecounter.AvgRateCounter
	debug        bool
	screenHeight int
	screenWidth  int
	scale        int
}

// Update progresses the game one tick, updating all entities that have been added to the game's board.
func (g *Game) Update(screen *ebiten.Image) error {
	if inpututil.IsKeyJustPressed(ebiten.KeyD) {
		g.debug = !g.debug
	}

	updateStart := time.Now()
	entityChan := g.gameboard.Entities()

	for e := range entityChan {
		e.Update()
	}

	g.updateTime.Incr(int64(time.Since(updateStart)))

	return nil
}

// Draw writes the screen image to the given ebiten.Image. All entities in the gameboard are given the chance to draw.
// draw order is not guaranteed.
func (g *Game) Draw(screen *ebiten.Image) {
	drawsStart := time.Now()
	screen.DrawImage(g.canvasImage, nil)

	entityChan := g.gameboard.Entities()
	for e := range entityChan {
		e.Draw(screen, g.scale)
	}

	g.drawTime.Incr(int64(time.Since(drawsStart)))

	if g.debug {
		msg := fmt.Sprintf(`FPS: %0.2f
Draw Time: %0.2f ms
Update Time %0.2f ms`, ebiten.CurrentFPS(), g.drawTime.Rate()/float64(time.Millisecond), g.updateTime.Rate()/float64(time.Millisecond))
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
func NewGame(width int, height int, background color.Color, scale int) *Game {
	g := Game{
		drawTime:     ratecounter.NewAvgRateCounter(time.Second),
		updateTime:   ratecounter.NewAvgRateCounter(time.Second),
		debug:        false,
		screenWidth:  width,
		screenHeight: height,
		scale:        scale,
	}
	g.canvasImage, _ = ebiten.NewImage(width, height, ebiten.FilterDefault)
	g.canvasImage.Fill(background)
	b := NewGameboard(g.screenWidth/g.scale, g.screenHeight/g.scale)
	g.gameboard = b

	return &g
}
