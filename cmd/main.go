package main

import (
	"log"

	"github.com/hajimehoshi/ebiten"
	"github.com/tannerhat/Cactus-Simulator/game"
	"github.com/tannerhat/Cactus-Simulator/nature"
)

const (
	screenWidth  = 600
	screenHeight = 350
	scale        = 5
)

func main() {
	ebiten.SetRunnableOnUnfocused(true)
	ebiten.SetWindowSize(screenWidth*2, screenHeight*2)
	ebiten.SetWindowTitle("Cactus Simulator")

	game := game.NewGame(screenWidth, screenHeight, scale)

	boardHeight := screenHeight / scale
	boardWidth := screenWidth / scale
	//cloudWidth := boardWidth / 8

	game.AddEntity(nature.NewWeather(1000))
	//game.AddEntity(nature.NewCloud((boardWidth/4)/2-(cloudWidth/2)+0*boardWidth/4, boardHeight/15+rand.Intn(boardHeight/15), cloudWidth, 2*cloudWidth/3, 1))
	//game.AddEntity(nature.NewCloud((boardWidth/4)/2-(cloudWidth/2)+1*boardWidth/4, boardHeight/15+rand.Intn(boardHeight/15), cloudWidth, 2*cloudWidth/3, 1))
	//game.AddEntity(nature.NewCloud((boardWidth/4)/2-(cloudWidth/2)+2*boardWidth/4, boardHeight/15+rand.Intn(boardHeight/15), cloudWidth, 2*cloudWidth/3, 1))
	//game.AddEntity(nature.NewCloud((boardWidth/4)/2-(cloudWidth/2)+3*boardWidth/4, boardHeight/15+rand.Intn(boardHeight/15), cloudWidth, 2*cloudWidth/3, 1))
	game.AddEntity(nature.NewSoil(0, boardHeight-3*boardHeight/6, boardWidth, 3*boardHeight/6))
	r := nature.NewRoots(0, boardHeight-3*boardHeight/6, boardWidth, 3*boardHeight/6, boardWidth/2, 0)
	game.AddEntity(r)
	game.AddEntity(nature.NewPlant(boardWidth/2, boardHeight-3*boardHeight/6-1, r))

	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
