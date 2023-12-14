package game

import (
	"fmt"
	"image/color"
	"rpg_demo/collisions"
	"rpg_demo/music"
	"rpg_demo/player"
	"rpg_demo/scene"
	"time"

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
	Music     bool
}

type KeyPressed struct {
	KeyP bool
	KeyZ bool
}

type Game struct {
	Player              *player.Player
	Scenes              map[string]*scene.Scene
	CurrentScene        string
	CurrentDoor         *collisions.Door
	State               GameState
	Transition          *Transition
	Music               *music.Music
	KeyPressedLastFrame KeyPressed
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
		Music: &music.Music{},
	}
}

func (g *Game) Update() error {
	Scene := g.Scenes[g.CurrentScene]
	g.HandleMusic()
	switch g.State {
	case PlayState:
		err := g.Player.Update(Scene.Collisions, func(door *collisions.Door) {
			g.CurrentDoor = door
		}, func(state int) {
			g.State = GameState(state)
		})
		if err != nil {
			return err
		}
		Scene.Update()
		Scene.HandleNPCInteractions(g.Player, g.KeyPressedLastFrame.KeyZ)
		g.KeyPressedLastFrame.KeyZ = ebiten.IsKeyPressed(ebiten.KeyZ)

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
	Scene := g.Scenes[g.CurrentScene]
	switch g.State {
	case PlayState:
		Scene.Draw(screen, Scene.Background, g.Player.X, g.Player.Y)
		Scene.DrawNPCs(screen)
		g.Player.Draw(screen, Scene.Width, Scene.Height)
		Scene.Draw(screen, Scene.Foreground, g.Player.X, g.Player.Y)
	case TransitionState, NewSceneState:
		Scene.Draw(screen, Scene.Background, g.Player.X, g.Player.Y)
		Scene.DrawNPCs(screen)
		g.Player.Draw(screen, Scene.Width, Scene.Height)
		Scene.Draw(screen, Scene.Foreground, g.Player.X, g.Player.Y)
		fadeImage := ebiten.NewImage(screen.Bounds().Dx(), screen.Bounds().Dy())
		fadeColor := color.RGBA{0, 0, 0, uint8(g.Transition.Alpha * 0xff)} // Black with variable Alpha
		fadeImage.Fill(fadeColor)
		screen.DrawImage(fadeImage, nil)
	}

}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return 320 * 2.5, 240 * 2.5 //Mutiplied by 2.5
}

func (g *Game) HandleMusic() {
	Scene := g.Scenes[g.CurrentScene]
	if ebiten.IsKeyPressed(ebiten.KeyP) && !g.KeyPressedLastFrame.KeyP {
		if g.Music.Paused && !g.Music.IsPlaying() {
			ch := make(chan struct{})
			go g.Music.FadeIn(time.Millisecond*500, ch)
			g.Music.PlayAudio()
		} else {
			ch := make(chan struct{})
			go g.Music.FadeOut(time.Millisecond*500, ch)
			g.Music.Paused = true
		}

	}
	g.KeyPressedLastFrame.KeyP = ebiten.IsKeyPressed(ebiten.KeyP)
	if g.Music.IsEmpty() {
		fmt.Println("Empty")
		g.Music.LoadAudio("./assets/" + Scene.Music)
		g.Music.PlayAudio()
	} else if !g.Music.IsPlaying() && !g.Music.Paused {
		g.Music.RewindMusic()
		fmt.Println("looped")
	} else if g.Music.CurrentSong != "./assets/"+Scene.Music && !g.Transition.Music {
		fmt.Println("New Song")
		g.Transition.Music = true
		// Channel to signal when fade-out is complete
		doneChan := make(chan struct{})

		// Start fade-out in a goroutine
		go g.Music.FadeOut(time.Second, doneChan)

		// Wait for the fade-out to complete in another goroutine
		go func() {
			<-doneChan // Wait for fade-out to complete
			// Load and play the new audio
			// Safe closure and loading of new audio
			if g.Music.IsPlaying() || g.Music.Paused {
				g.Music.CloseAudio() // Ensure the current audio is closed
			}
			g.Music.LoadAudio("./assets/" + Scene.Music)
			g.Music.PlayAudio()
			g.Transition.Music = false
		}()
	}

}
