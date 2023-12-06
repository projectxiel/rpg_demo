package game

import (
	"rpg_demo/player"
	"rpg_demo/scene"

	"github.com/hajimehoshi/ebiten/v2"
)

type Game struct {
	Player *player.Player
	Scene  *scene.Scene
}

func New() *Game {
	return &Game{
		Player: player.New(),
		Scene:  scene.New("assets/mainMap.png", "assets/mainMapFore.png"),
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
	g.Scene.DrawBg(screen)
	g.Player.Draw(screen)
	g.Scene.DrawFg(screen)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return 320, 240
}
