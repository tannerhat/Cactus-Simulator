package main

import (
	"image/color"
	"log"

	"github.com/hajimehoshi/ebiten"
	"github.com/tannerhat/Cactus-Simulator/game"
	"github.com/tannerhat/Cactus-Simulator/nature"
)

const (
	screenWidth  = 300
	screenHeight = 700
	cellSize     = 10
)

func main() {
	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("Cactus Simulator")

	game := game.NewGame(screenWidth, screenHeight, color.RGBA{0x87, 0xce, 0xfa, 0xff}, cellSize)

	boardHeight := screenHeight / cellSize
	boardWidth := screenWidth / cellSize

	game.AddEntity(nature.NewCloud(2, 5, 10, 5, 1))
	game.AddEntity(nature.NewCloud(boardWidth-10-2, 10, 10, 5, 1))
	game.AddEntity(nature.NewSoil(0, boardHeight-30, boardWidth, 30))
	r := nature.NewRoots(0, boardHeight-30, boardWidth, 30, boardWidth/2, 0)
	game.AddEntity(r)
	game.AddEntity(nature.NewPlant(boardWidth/2, boardHeight-31, r))

	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
