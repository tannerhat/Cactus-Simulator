package game

import (
	"sync"
)

type GameBoard interface {
	MoveEntity(px int, py int, x int, y int)
	AddEntity(e Entity)
	Entities() <-chan Entity
	SetEntity(e Entity, x int, y int)
	EntityAt(x int, y int) Entity
	Size() (width int, height int)
	RemoveEntity(e Entity)
}

type gameBoard struct {
	entityLock sync.RWMutex
	entities   []Entity
	board      [][]Entity
}

func NewGameboard(width int, height int) GameBoard {
	g := &gameBoard{
		board: make([][]Entity, width/5),
	}
	for i := range g.board {
		g.board[i] = make([]Entity, height/5)
	}
	return g
}

func (g *gameBoard) RemoveEntity(e Entity) {
	g.entityLock.Lock()
	defer g.entityLock.Unlock()
	for i, ent := range g.entities {
		if e == ent {
			g.entities = append(g.entities[:i], g.entities[i+1:]...)
		}
	}
}

func (g *gameBoard) MoveEntity(px int, py int, x int, y int) {
	g.board[x][y] = g.board[px][py]
	g.board[px][py] = nil
}

func (g *gameBoard) Entities() <-chan Entity {
	g.entityLock.Lock()
	defer g.entityLock.Unlock()

	c := make(chan Entity, len(g.entities))
	for _, value := range g.entities {
		c <- value
	}
	close(c)

	return c
}

func (g *gameBoard) AddEntity(e Entity) {
	g.entityLock.Lock()
	defer g.entityLock.Unlock()
	g.entities = append(g.entities, e)
	e.AddToBoard(g)
}

func (g *gameBoard) SetEntity(e Entity, x int, y int) {
	g.board[x][y] = e
}

func (g *gameBoard) EntityAt(x int, y int) Entity {
	return g.board[x][y]
}

func (g *gameBoard) Size() (int, int) {
	return len(g.board), len(g.board[0])
}
