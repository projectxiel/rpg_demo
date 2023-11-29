package player

import (
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

type Player struct {
	SpriteSheet *ebiten.Image
}

func New() *Player {
	img, _, err := ebitenutil.NewImageFromFile("assets/playerDownBlack.png")
	if err != nil {
		log.Fatal(err)
	}
	return &Player{
		SpriteSheet: img,
	}
}

func (p *Player) Draw(screen *ebiten.Image) {
	opts := &ebiten.DrawImageOptions{}
	screen.DrawImage(p.SpriteSheet, opts)
}
