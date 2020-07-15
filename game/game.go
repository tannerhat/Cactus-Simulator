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

type Game struct {
	canvasImage  *ebiten.Image
	gameBoard    GameBoard
	drawTime     *ratecounter.AvgRateCounter
	updateTime   *ratecounter.AvgRateCounter
	debug        bool
	screenHeight int
	screenWidth  int
}

func (g *Game) Update(screen *ebiten.Image) error {
	if inpututil.IsKeyJustPressed(ebiten.KeyD) {
		g.debug = !g.debug
	}

	updateStart := time.Now()
	entityChan := g.gameBoard.Entities()

	for e := range entityChan {
		e.Update(g.gameBoard)
	}

	g.updateTime.Incr(int64(time.Since(updateStart)))

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	drawsStart := time.Now()
	screen.DrawImage(g.canvasImage, nil)

	entityChan := g.gameBoard.Entities()
	for e := range entityChan {
		e.Draw(screen, 5)
	}

	g.drawTime.Incr(int64(time.Since(drawsStart)))

	if g.debug {
		msg := fmt.Sprintf(`TPS: %0.2f
FPS: %0.2f
Draw Time: %0.2f ms
Update Time %0.2f ms`, ebiten.CurrentTPS(), ebiten.CurrentFPS(), g.drawTime.Rate()/float64(time.Millisecond), g.updateTime.Rate()/float64(time.Millisecond))
		ebitenutil.DebugPrint(screen, msg)
	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return g.screenWidth, g.screenHeight
}

func (g *Game) AddEntity(entity Entity) {
	g.gameBoard.AddEntity(entity)
}

func NewGame(width int, height int, background color.Color) *Game {
	g := Game{
		drawTime:     ratecounter.NewAvgRateCounter(time.Second),
		updateTime:   ratecounter.NewAvgRateCounter(time.Second),
		debug:        false,
		screenWidth:  width,
		screenHeight: height,
	}
	g.canvasImage, _ = ebiten.NewImage(width, height, ebiten.FilterDefault)
	g.canvasImage.Fill(background)
	b := NewGameboard(width, height)
	g.gameBoard = b

	return &g
}
