package models

import (
	"log"

	"github.com/hajimehoshi/ebiten/v2"

	"github.com/hajimehoshi/ebiten/v2/ebitenutil"

	_ "image/png"
)

type floorTile struct {
	gameObject
	image *ebiten.Image
}

type Floor struct {
	sprites []*floorTile
}

func CreateFloor() *Floor {

	var err error

	image, _, err := ebitenutil.NewImageFromFile("assets/floor.png")

	if err != nil {
		log.Fatal(err)
	}

	sprites := make([]*floorTile, f_spriteAmount)

	startingXPos := -f_width;
	startingYPos := f_YPos

	for i := 0; i < f_spriteAmount; i++ {
    	sprites[i] = &floorTile{gameObject{startingXPos + (float64(i) * f_width), startingYPos, f_velocity}, image}
	}

	return &Floor{sprites}
}

func (floor *Floor) Update(game *Game) {
	if (game.count % 2 == 0) {
		for i:= len(floor.sprites)-1; i >= 0; i-- {
			tile := floor.sprites[i]

			tile.xPos += tile.velocity * g_universalVelocity

			if (tile.xPos > g_screenWidth) {
				tile.xPos = -f_width + (tile.xPos - g_screenWidth)
			}
		}
	}
}

func (floor *Floor) Draw(game *Game, screen *ebiten.Image) {
	for i:= len(floor.sprites)-1; i >= 0; i-- {
		op := &ebiten.DrawImageOptions{}
		tile := floor.sprites[i]
		op.GeoM.Translate(tile.xPos, tile.yPos)
		screen.DrawImage(tile.image, op)
	}
}