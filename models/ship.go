package models

import (
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"

	_ "image/png"
)

type Ship struct {
	gameObject
	image *ebiten.Image
}

func CreateShip() *Ship {

	var err error

	image, _, err := ebitenutil.NewImageFromFile("assets/ship.png")

	if err != nil {
		log.Fatal(err)
	}

	startingXPos := -f_width;
	startingYPos := float64(f_YPos - 20)

	return &Ship{gameObject{startingXPos, startingYPos, 2}, image}
}

func (ship *Ship) Update(game *Game) {
	if (game.count % 4 == 0) {
		ship.xPos += ship.velocity
	}

	if ship.xPos > g_screenWidth {
		ship.xPos = -f_width
	}
}

func (ship *Ship) Draw(game *Game, screen *ebiten.Image) {
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(ship.xPos, ship.yPos)
	screen.DrawImage(ship.image, op)
}