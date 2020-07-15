package main

import (
	"fmt"
	"image/color"
	"log"
	"time"

	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/ebitenutil"
	"github.com/hajimehoshi/ebiten/inpututil"
	"github.com/paulbellamy/ratecounter"
	"github.com/tannerhat/Cactus-Simulator/game"
	"github.com/tannerhat/Cactus-Simulator/nature"
)

const (
	screenWidth  = 150
	screenHeight = 350
)

type Game struct {
	canvasImage *ebiten.Image
	gameBoard   game.GameBoard
	drawTime    *ratecounter.AvgRateCounter
	updateTime  *ratecounter.AvgRateCounter
	debug       bool
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
	return screenWidth, screenHeight
}

func NewGame(width, height int) *Game {
	g := Game{
		drawTime:   ratecounter.NewAvgRateCounter(time.Second),
		updateTime: ratecounter.NewAvgRateCounter(time.Second),
		debug:      false,
	}
	g.canvasImage, _ = ebiten.NewImage(width, height, ebiten.FilterDefault)
	g.canvasImage.Fill(color.RGBA{0x87, 0xce, 0xfa, 0xff})
	b := game.NewGameboard(width, height)
	g.gameBoard = b

	boardHeight := height / 5
	boardWidth := width / 5

	g.gameBoard.AddEntity(nature.NewCloud(2, 5, 10, 5, 1))
	g.gameBoard.AddEntity(nature.NewCloud(boardWidth-10-2, 10, 10, 5, 1))
	g.gameBoard.AddEntity(nature.NewSoil(0, boardHeight-30, boardWidth, 30))
	r := nature.NewRoots(0, boardHeight-30, boardWidth, 30, boardWidth/2, 0)
	g.gameBoard.AddEntity(r)
	g.gameBoard.AddEntity(nature.NewPlant(boardWidth/2, boardHeight-31, r))
	return &g
}

func main() {
	ebiten.SetWindowSize(screenWidth*2, screenHeight*2)
	ebiten.SetWindowTitle("Cactus Simulator")
	if err := ebiten.RunGame(NewGame(screenWidth, screenHeight)); err != nil {
		log.Fatal(err)
	}
}
