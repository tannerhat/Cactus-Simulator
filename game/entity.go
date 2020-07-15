package game

import (
	"github.com/hajimehoshi/ebiten"
)

type Entity interface {
	Draw(screen *ebiten.Image, scale int)
	Update(gameBoard GameBoard)
	AddToBoard(gameBoard GameBoard)
	Name() string
}
