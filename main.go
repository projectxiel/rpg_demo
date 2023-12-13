package main

import (
	"log"
	g "rpg_demo/game"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/audio"
)

func main() {
	ebiten.SetWindowSize(640, 480)
	ebiten.SetWindowTitle("My Game")
	game := g.New()
	game.Music.SetCtx(audio.NewContext(44100))
	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
