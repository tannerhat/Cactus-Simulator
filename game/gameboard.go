package game

import (
	"sync"
)

// Gameboard tracks the entities in play and the game locations of any solid entities. A single entity may exist at multiple locations. An entity may also not have any game location.
type Gameboard interface {
	// MoveEntity moves the entity at (px,py) to (x,y), (px,py) will be empty after this
	MoveEntity(px int, py int, x int, y int)

	// AddEntity add the entity to the list that Gameboard tracks. It will also call the entity's AddToBoard
	AddEntity(e Entity)

	// Entities returns a channel of all entities in the board at time of calling. Changes to the entity list after calling will not affect the channel's contents.
	Entities() <-chan Entity

	// SetEntity puts e at the game location (x,y)
	SetEntity(e Entity, x int, y int)

	// EntityAt returns the entity at the game location (x,y)
	EntityAt(x int, y int) Entity

	// Size returns the width and height of the game board
	Size() (width int, height int)

	// RemoveEntity takes the given entity out of the entity list
	RemoveEntity(e Entity)
}

type gameboard struct {
	entityLock sync.RWMutex
	entities   []Entity
	board      [][]Entity
}

// NewGameboard gives a simple implementation of Gameboard with the given width and height
func NewGameboard(width int, height int) Gameboard {
	g := &gameboard{
		board: make([][]Entity, width),
	}
	for i := range g.board {
		g.board[i] = make([]Entity, height)
	}
	return g
}

// RemoveEntity takes the given entity out of the entity list, the caller is expected to have
// already removed the entity's game locations using SetEntity(nil,x,y) for all locations it occupied.
func (g *gameboard) RemoveEntity(e Entity) {
	g.entityLock.Lock()
	defer g.entityLock.Unlock()
	for i, ent := range g.entities {
		if e == ent {
			g.entities = append(g.entities[:i], g.entities[i+1:]...)
		}
	}
}

func (g *gameboard) MoveEntity(px int, py int, x int, y int) {
	g.board[x][y] = g.board[px][py]
	g.board[px][py] = nil
}

func (g *gameboard) Entities() <-chan Entity {
	g.entityLock.Lock()
	defer g.entityLock.Unlock()

	c := make(chan Entity, len(g.entities))
	for _, value := range g.entities {
		c <- value
	}
	close(c)

	return c
}

func (g *gameboard) AddEntity(e Entity) {
	g.entityLock.Lock()
	defer g.entityLock.Unlock()
	g.entities = append(g.entities, e)
	e.AddToBoard(g)
}

func (g *gameboard) SetEntity(e Entity, x int, y int) {
	g.board[x][y] = e
}

func (g *gameboard) EntityAt(x int, y int) Entity {
	return g.board[x][y]
}

func (g *gameboard) Size() (int, int) {
	return len(g.board), len(g.board[0])
}
