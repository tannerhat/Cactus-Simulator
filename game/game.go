package game

import (
	"fmt"
	"image/color"
	"log"
	"time"

	"github.com/golang/freetype/truetype"
	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/ebitenutil"
	"github.com/hajimehoshi/ebiten/examples/resources/fonts"
	"github.com/hajimehoshi/ebiten/inpututil"
	"github.com/hajimehoshi/ebiten/text"
	"github.com/paulbellamy/ratecounter"
	"golang.org/x/image/font"
)

var (
	arcadeFont font.Face
	fontSize   = 16
)

type Mode int

const (
	ModeTitle Mode = iota
	ModeGame
	ModeWin
)

// Game implements ebiten.Game and keeps track of the gameboard and entities.
type Game struct {
	gameboard    Gameboard
	drawTime     *ratecounter.AvgRateCounter
	updateTime   *ratecounter.AvgRateCounter
	debug        bool
	screenHeight int
	screenWidth  int
	scale        int
	speed        int
	ticks        int
	mode         Mode
}

func init() {
	tt, err := truetype.Parse(fonts.ArcadeN_ttf)
	if err != nil {
		log.Fatal(err)
	}
	const dpi = 72
	arcadeFont = truetype.NewFace(tt, &truetype.Options{
		Size:    float64(fontSize),
		DPI:     dpi,
		Hinting: font.HintingFull,
	})
}

// Update progresses the game one tick, updating all entities that have been added to the game's board.
func (g *Game) Update(screen *ebiten.Image) error {
	updateStart := time.Now()

	if g.mode == ModeGame {
		if inpututil.IsKeyJustPressed(ebiten.KeyGraveAccent) {
			g.speed = 0
		}
		if inpututil.IsKeyJustPressed(ebiten.Key1) {
			g.speed = 1
		}
		if inpututil.IsKeyJustPressed(ebiten.Key2) {
			g.speed = 10
		}
		if inpututil.IsKeyJustPressed(ebiten.Key3) {
			g.speed = 60
		}
		if inpututil.IsKeyJustPressed(ebiten.Key4) {
			g.speed = 300
		}
		if inpututil.IsKeyJustPressed(ebiten.KeyD) {
			g.debug = !g.debug
		}

		for i := 0; i < g.speed; i++ {
			g.ticks++

			entityChan := g.gameboard.Entities()

			for e := range entityChan {
				e.Update()
				if win, ok := e.(Winnable); ok {
					if win.Win() {
						g.mode = ModeWin
					}
				}
			}
		}
	} else if g.mode == ModeTitle {
		if inpututil.IsKeyJustPressed(ebiten.KeySpace) {
			g.mode = ModeGame
		}
	} else if g.mode == ModeWin {
		if inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
			return fmt.Errorf("game dones")
		}
		g.speed = 0
	}
	g.updateTime.Incr(int64(time.Since(updateStart)))
	return nil
}

// Draw writes the screen image to the given ebiten.Image. All entities in the gameboard are given the chance to draw.
// draw order is not guaranteed.
func (g *Game) Draw(screen *ebiten.Image) {
	drawsStart := time.Now()

	if g.mode == ModeGame || g.mode == ModeWin {

		entityChan := g.gameboard.Entities()
		entityList := []Entity{}
		maxLayer := 0
		for e := range entityChan {
			entityList = append(entityList, e)
			if e.Layer() > maxLayer {
				maxLayer = e.Layer()
			}
		}

		// reuse the entity list rather than get a new entities channel from gameboard because entities could've changed
		for layer := 0; layer <= maxLayer; layer++ {
			for _, e := range entityList {
				if e.Layer() == layer {
					e.Draw(screen, g.scale)
				}
			}
		}

		g.drawTime.Incr(int64(time.Since(drawsStart)))

		if g.debug {
			msg := fmt.Sprintf(`FPS: %0.2f
Draw Time: %0.2f ms
Update Time: %0.2f ms
Speed: %d
Game time: %0.2f hours`,
				ebiten.CurrentFPS(),
				g.drawTime.Rate()/float64(time.Millisecond),
				g.updateTime.Rate()/float64(time.Millisecond),
				g.speed,
				float64(g.ticks)/60.0/60.0/60.0)
			ebitenutil.DebugPrint(screen, msg)
		}
		if g.mode == ModeWin {
			texts := []string{"", "", "", "YOU GREW THE PERFECT:", "", "CACTUS", "", "Press Escape to Leave."}
			for i, l := range texts {
				x := (g.screenWidth - len(l)*fontSize) / 2
				text.Draw(screen, l, arcadeFont, x, (i+4)*fontSize, color.White)
			}
		}
	} else if g.mode == ModeTitle {
		texts := []string{"Welcome To Cactus Simulator", "", "Controls:", "~: pause", "1: 1x speed", "2: 10x speed", "3: 60x speed", "4: 300x speed", "space: abosorb water", "d: debug info", "", "", "Press spacebar to start"}
		for i, l := range texts {
			x := (g.screenWidth - len(l)*fontSize) / 2
			text.Draw(screen, l, arcadeFont, x, (i+4)*fontSize, color.White)
		}
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
func NewGame(width int, height int, scale int) *Game {
	g := Game{
		drawTime:     ratecounter.NewAvgRateCounter(time.Second),
		updateTime:   ratecounter.NewAvgRateCounter(time.Second),
		debug:        false,
		screenWidth:  width,
		screenHeight: height,
		scale:        scale,
		speed:        1,
		ticks:        0,
		mode:         ModeTitle,
	}
	b := NewGameboard(g.screenWidth/g.scale, g.screenHeight/g.scale)
	g.gameboard = b

	return &g
}
