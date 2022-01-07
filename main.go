package main

import (
	"log"
	"math/rand"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/rickread/runnergame/models"
)

func main() {
	rand.Seed(time.Now().UnixNano())
	game := &models.Game{}
	ebiten.SetRunnableOnUnfocused(true)
	ebiten.SetWindowResizable(false)
	ebiten.SetWindowTitle("Go Runner Go!")
	ebiten.SetWindowSize(game.GetDefaultWindowSize())

	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}