package main

import (
	"log"
	g "rpg_demo/game"

	"github.com/hajimehoshi/ebiten/v2"
)

func main() {
	ebiten.SetWindowSize(640, 480)
	ebiten.SetWindowTitle("My Game")
	game := g.New()
	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
