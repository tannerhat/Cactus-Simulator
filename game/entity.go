package game

import (
	"github.com/hajimehoshi/ebiten"
)

// Entity defines necessary functions for an entity that can be added to the game
type Entity interface {
	// Draw the entity to the given screen at the given scale
	Draw(screen *ebiten.Image, scale int)

	// Update the entity by one tick
	Update()

	// AddToBoard is called by game when an entity is added, the entity marks its initial positions on the game board
	AddToBoard(gameboard Gameboard)

	// Layer returns the layer of the entity for draw purposes
	Layer() int
}
