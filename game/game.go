package game

import (
	"image/color"
	"rpg_demo/collisions"
	"rpg_demo/player"
	"rpg_demo/scene"

	"github.com/hajimehoshi/ebiten/v2"
)

type GameState int

const (
	PlayState GameState = iota
	TransitionState
	NewSceneState
)

type Transition struct {
	Alpha     float64
	FadeSpeed float64
}

type Game struct {
	Player       *player.Player
	Scenes       map[string]*scene.Scene
	CurrentScene string
	CurrentDoor  *collisions.Door
	State        GameState
	Transition   *Transition
}

func New() *Game {
	sceneMap := make(map[string]*scene.Scene)
	sceneMap["mainMap"] = scene.New("mainMap")
	return &Game{
		Player:       player.New(),
		CurrentScene: "mainMap",
		Scenes:       sceneMap,
		Transition: &Transition{
			Alpha:     0.0,
			FadeSpeed: 0.05,
		},
	}
}

func (g *Game) Update() error {
	switch g.State {
	case PlayState:
		err := g.Player.Update(g.Scenes[g.CurrentScene].Collisions, func(door *collisions.Door) {
			g.CurrentDoor = door
		}, func(state int) {
			g.State = GameState(state)
		})
		if err != nil {
			return err
		}
	case TransitionState:
		g.Transition.Alpha += g.Transition.FadeSpeed
		if g.Transition.Alpha >= 1.0 {
			g.Transition.Alpha = 1.0
			g.State = NewSceneState
			g.CurrentScene = g.CurrentDoor.Destination
			g.Player.X = g.CurrentDoor.NewX
			g.Player.Y = g.CurrentDoor.NewY
		}
	case NewSceneState:
		g.Transition.Alpha -= g.Transition.FadeSpeed
		if g.Transition.Alpha <= 0.0 {
			g.Transition.Alpha = 0.0
			g.State = PlayState
			// The new scene is fully visible now, and game continues as normal
		}
	}
	_, exists := g.Scenes[g.CurrentScene]
	if !exists {
		g.Scenes[g.CurrentScene] = scene.New(g.CurrentScene)
	}
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	switch g.State {
	case PlayState:
		g.Scenes[g.CurrentScene].Draw(screen, g.Scenes[g.CurrentScene].Background, g.Player.X, g.Player.Y)
		g.Player.Draw(screen, g.Scenes[g.CurrentScene].Width, g.Scenes[g.CurrentScene].Height)
		g.Scenes[g.CurrentScene].Draw(screen, g.Scenes[g.CurrentScene].Foreground, g.Player.X, g.Player.Y)
	case TransitionState, NewSceneState:
		g.Scenes[g.CurrentScene].Draw(screen, g.Scenes[g.CurrentScene].Background, g.Player.X, g.Player.Y)
		g.Player.Draw(screen, g.Scenes[g.CurrentScene].Width, g.Scenes[g.CurrentScene].Height)
		g.Scenes[g.CurrentScene].Draw(screen, g.Scenes[g.CurrentScene].Foreground, g.Player.X, g.Player.Y)
		fadeImage := ebiten.NewImage(screen.Bounds().Dx(), screen.Bounds().Dy())
		fadeColor := color.RGBA{0, 0, 0, uint8(g.Transition.Alpha * 0xff)} // Black with variable Alpha
		fadeImage.Fill(fadeColor)
		screen.DrawImage(fadeImage, nil)
	}

}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return 320 * 2.5, 240 * 2.5 //Mutiplied by 2.5
}
