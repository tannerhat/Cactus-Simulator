package main

import (
	"fmt"
	"image/color"
	"log"
	"time"

	"github.com/hajimehoshi/ebiten"
	"github.com/tannerhat/Cactus-Simulator/nature"
)

const (
	screenWidth  = 150
	screenHeight = 350
)

type avgStruct struct {
	sum   time.Duration
	count int
}

type Game struct {
	canvasImage        *ebiten.Image
	gameBoard          nature.GameBoard
	updateCount        int
	averageProcessTime map[string]*avgStruct
	frames             int
	start              time.Time
}

func (g *Game) Update(screen *ebiten.Image) error {
	g.updateCount++
	if g.updateCount%1 == 0 {
		entityChan := g.gameBoard.Entities()

		for e := range entityChan {
			entStart := time.Now()
			entName := e.Name()
			e.Update(g.gameBoard)

			avg, ok := g.averageProcessTime[entName]
			if !ok {
				g.averageProcessTime[entName] = &avgStruct{
					sum:   0,
					count: 0,
				}
				avg = g.averageProcessTime[entName]
			}

			avg.sum += time.Since(entStart)
			avg.count++
		}
	}

	g.frames++
	if g.updateCount%60 == 0 {
		for k, v := range g.averageProcessTime {
			fmt.Printf("%s\t%d\t%d\n", k, int(v.sum)/v.count, v.count)
		}

		if g.updateCount%300 == 0 {
			fmt.Printf("framerate: %d\n", g.frames/int(time.Since(g.start)/time.Second))
			g.frames = 0
			g.start = time.Now()
		}

	}

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	drawsStart := time.Now()
	screen.DrawImage(g.canvasImage, nil)

	entityChan := g.gameBoard.Entities()
	for e := range entityChan {
		entStart := time.Now()
		entName := e.Name()
		e.Draw(screen, 5)
		fmt.Printf("drawl %s: %d\n", entName, time.Since(entStart))

	}
	draw2 := time.Now()
	fmt.Printf("drawtime: %d %d\n", time.Since(drawsStart), time.Since(draw2))
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

func NewGame(width, height int) *Game {
	g := Game{
		updateCount:        0,
		averageProcessTime: map[string]*avgStruct{},
	}
	g.canvasImage, _ = ebiten.NewImage(width, height, ebiten.FilterDefault)
	g.canvasImage.Fill(color.RGBA{0x87, 0xce, 0xfa, 0xff})
	b := nature.NewGameboard(width, height)
	g.gameBoard = b
	g.frames = 0
	g.start = time.Now()

	boardHeight := height / 5
	boardWidth := width / 5

	g.gameBoard.AddEntity(nature.NewCloud(10, 5, 10, 5, 1))
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
