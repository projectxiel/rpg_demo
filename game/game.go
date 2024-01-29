package game

import (
	"image/color"
	"rpg_demo/collisions"
	"rpg_demo/cutscene"
	"rpg_demo/dialogue"
	"rpg_demo/music"
	"rpg_demo/player"
	"rpg_demo/scene"
	"rpg_demo/shared"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
)

type Game struct {
	Player              *player.Player
	Scenes              map[string]*scene.Scene
	CurrentScene        string
	CurrentDoor         *collisions.Door
	CutScene            *cutscene.Cutscene
	State               shared.GameState
	Transition          *shared.Transition
	Music               *music.Music
	KeyPressedLastFrame shared.KeyPressed
	Dialogue            *dialogue.Dialogue
}

func New() *Game {
	sceneMap := make(map[string]*scene.Scene)
	sceneMap["mainMap"] = scene.New("mainMap")
	return &Game{
		Player:       player.New(),
		CurrentScene: "mainMap",
		Scenes:       sceneMap,
		Transition: &shared.Transition{
			Alpha:     0.0,
			FadeSpeed: 0.05,
		},
		Music:    &music.Music{},
		Dialogue: dialogue.New(),
	}
}

func (g *Game) Update() error {
	Scene := g.Scenes[g.CurrentScene]
	g.HandleMusic()
	switch g.State {
	case shared.PlayState:
		err := g.Player.Update(Scene.Collisions, func(door *collisions.Door) {
			g.CurrentDoor = door
		}, func(state shared.GameState) {
			g.State = state
		})
		if err != nil {
			return err
		}
		Scene.Update()
		Scene.HandleNPCInteractions(g.Player, g.KeyPressedLastFrame.KeyZ, g.Dialogue)
		g.KeyPressedLastFrame.KeyZ = ebiten.IsKeyPressed(ebiten.KeyZ)
		if ebiten.IsKeyPressed(ebiten.KeyD) && !g.KeyPressedLastFrame.KeyD {
			c := createExampleCutscene(g)
			g.CutScene = &c
			g.CutScene.Start()
			g.State = shared.CutSceneState
		}
		g.KeyPressedLastFrame.KeyD = ebiten.IsKeyPressed(ebiten.KeyD)

	case shared.TransitionState:
		g.Transition.Alpha += g.Transition.FadeSpeed
		if g.Transition.Alpha >= 1.0 {
			g.Transition.Alpha = 1.0
			g.State = shared.NewSceneState
			g.CurrentScene = g.CurrentDoor.Destination
			g.Player.X = g.CurrentDoor.NewX
			g.Player.Y = g.CurrentDoor.NewY
		}
	case shared.NewSceneState:
		g.Transition.Alpha -= g.Transition.FadeSpeed
		if g.Transition.Alpha <= 0.0 {
			g.Transition.Alpha = 0.0
			g.State = shared.PlayState
			// The new scene is fully visible now, and game continues as normal
		}
	case shared.CutSceneState:
		if g.CutScene.IsPlaying {
			g.CutScene.Update((*cutscene.Transition)(g.Transition), cutscene.KeyPressed(g.KeyPressedLastFrame))
		} else {
			g.State = shared.PlayState
		}

	}
	if g.Dialogue.IsOpen && !g.Dialogue.Finished {
		g.Dialogue.Update()
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
	case shared.PlayState:
		Scene.Draw(screen, Scene.Background, g.Player.X, g.Player.Y)
		Scene.DrawNPCs(screen)
		g.Player.Draw(screen, Scene.Width, Scene.Height)
		Scene.Draw(screen, Scene.Foreground, g.Player.X, g.Player.Y)
		g.Dialogue.Draw(screen)
	case shared.TransitionState, shared.NewSceneState:
		Scene.Draw(screen, Scene.Background, g.Player.X, g.Player.Y)
		Scene.DrawNPCs(screen)
		g.Player.Draw(screen, Scene.Width, Scene.Height)
		Scene.Draw(screen, Scene.Foreground, g.Player.X, g.Player.Y)
		fadeImage := ebiten.NewImage(screen.Bounds().Dx(), screen.Bounds().Dy())
		fadeColor := color.RGBA{0, 0, 0, uint8(g.Transition.Alpha * 0xff)} // Black with variable Alpha
		fadeImage.Fill(fadeColor)
		screen.DrawImage(fadeImage, nil)
		fadeImage.Dispose()
	case shared.CutSceneState:
		Scene.Draw(screen, Scene.Background, g.Player.X, g.Player.Y)
		Scene.DrawNPCs(screen)
		g.Player.Draw(screen, Scene.Width, Scene.Height)
		Scene.Draw(screen, Scene.Foreground, g.Player.X, g.Player.Y)
		fadeImage := ebiten.NewImage(screen.Bounds().Dx(), screen.Bounds().Dy())
		fadeColor := color.RGBA{0, 0, 0, uint8(g.Transition.Alpha * 0xff)} // Black with variable Alpha
		fadeImage.Fill(fadeColor)
		g.Dialogue.Draw(screen)
		screen.DrawImage(fadeImage, nil)
		fadeImage.Dispose()
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
		g.Music.LoadAudio("./assets/" + Scene.Music)
		g.Music.PlayAudio()
	} else if !g.Music.IsPlaying() && !g.Music.Paused {
		g.Music.RewindMusic()
	} else if g.Music.CurrentSong != "./assets/"+Scene.Music && !g.Transition.Music {
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

func createExampleCutscene(g *Game) cutscene.Cutscene {
	return cutscene.Cutscene{
		Actions: []cutscene.CutsceneAction{
			{
				ActionType:   cutscene.FadeOut,
				Data:         0.01,
				WaitPrevious: false,
			},
			{
				ActionType:   cutscene.TeleportPlayer,
				Target:       g.Player,
				Data:         cutscene.Vector2D{X: 0 + float64(g.Player.Frame.Width)/2, Y: 0 + float64(g.Player.Frame.Height)/2},
				WaitPrevious: true,
			},
			{
				ActionType:   cutscene.TeleportNPC,
				Target:       g.Scenes[g.CurrentScene].NPCs["Bryan"],
				Data:         cutscene.Vector2D{X: 0, Y: 0},
				WaitPrevious: true,
			},
			{
				ActionType:   cutscene.FadeIn,
				Data:         0.01,
				WaitPrevious: true,
			},
			{
				ActionType:   cutscene.MovePlayer,
				Target:       g.Player,
				Data:         cutscene.Vector2D{X: 150 + float64(g.Player.Frame.Width)/2, Y: 200 + float64(g.Player.Frame.Height)/2}, // Target position for player
				WaitPrevious: true,
			},
			{
				ActionType:   cutscene.MoveNPC,
				Target:       g.Scenes[g.CurrentScene].NPCs["Bryan"],
				Data:         cutscene.Vector2D{X: 200, Y: 200}, // Target position for NPC
				WaitPrevious: false,
			},
			{
				ActionType:   cutscene.TurnNPC,
				Target:       g.Scenes[g.CurrentScene].NPCs["Bryan"],
				Data:         "left",
				WaitPrevious: true,
			},
			{
				ActionType:   cutscene.TurnPlayer,
				Target:       g.Player,
				Data:         "right",
				WaitPrevious: true,
			},
			{
				ActionType:   cutscene.ShowDialogue,
				Target:       g.Dialogue,
				Data:         []string{"This is our first Scene.", "Pretty Cool huh?"},
				WaitPrevious: true,
			},
		},
	}
}
