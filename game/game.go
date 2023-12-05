package game

import (
	"rpg_demo/player"

	"github.com/hajimehoshi/ebiten/v2"
)

type Game struct {
	Player *player.Player
}

func New() *Game {
	return &Game{
		Player: player.New(),
	}
}

func (g *Game) Update() error {
	err := g.Player.Update()
	if err != nil {
		return err
	}

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	g.Player.Draw(screen)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return 320, 240
}
