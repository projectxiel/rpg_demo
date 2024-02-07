package game

import (
	"fmt"
	"image/color"
	"rpg_demo/ability"
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
	if (!g.Player.Ability.Activated || g.Player.Ability.Type != ability.StopTime) && g.State == shared.TimeStopped {
		g.State = shared.PlayState
	}

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
			g.CutScene = Scene.Cutscenes["exampleCutscene"]
			g.processCutscene()
			fmt.Println(g.CutScene)
			g.CutScene.Start()
			g.State = shared.CutSceneState
		}
		g.KeyPressedLastFrame.KeyD = ebiten.IsKeyPressed(ebiten.KeyD)
		if ebiten.IsKeyPressed(ebiten.KeyV) && !g.KeyPressedLastFrame.KeyV {
			g.Player.Ability.CycleAbility()
		}
		g.KeyPressedLastFrame.KeyV = ebiten.IsKeyPressed(ebiten.KeyV)
		if g.Player.Ability.Type == ability.StopTime && g.Player.Ability.Activated {
			g.State = shared.TimeStopped
		}
	case shared.TimeStopped:
		g.Player.Update(Scene.Collisions, func(d *collisions.Door) { g.CurrentDoor = d }, func(newState shared.GameState) { g.State = newState })
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
			g.CutScene.Update(g.Transition, g.KeyPressedLastFrame)
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
	case shared.PlayState, shared.TimeStopped:
		Scene.Draw(screen, Scene.Background, g.Player)
		Scene.DrawNPCs(screen)
		g.Player.Draw(screen, Scene.Width, Scene.Height)
		Scene.Draw(screen, Scene.Foreground, g.Player)
		g.Dialogue.Draw(screen)
	case shared.TransitionState, shared.NewSceneState:
		Scene.Draw(screen, Scene.Background, g.Player)
		Scene.DrawNPCs(screen)
		g.Player.Draw(screen, Scene.Width, Scene.Height)
		Scene.Draw(screen, Scene.Foreground, g.Player)
		fadeImage := ebiten.NewImage(screen.Bounds().Dx(), screen.Bounds().Dy())
		fadeColor := color.RGBA{0, 0, 0, uint8(g.Transition.Alpha * 0xff)} // Black with variable Alpha
		fadeImage.Fill(fadeColor)
		screen.DrawImage(fadeImage, nil)
		fadeImage.Dispose()
	case shared.CutSceneState:
		Scene.Draw(screen, Scene.Background, g.Player)
		Scene.DrawNPCs(screen)
		g.Player.Draw(screen, Scene.Width, Scene.Height)
		Scene.Draw(screen, Scene.Foreground, g.Player)
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
	} else if g.Music.CurrentSong != "./assets/"+Scene.Music && !g.Transition.Music && g.State != shared.CutSceneState {
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

func (g *Game) processCutscene() {
	for i := range g.CutScene.Actions {
		target := g.resolveTarget(g.CutScene.Actions[i].Target)
		data := g.resolveData(g.CutScene.Actions[i].Data)
		// fmt.Println(target)
		g.CutScene.Actions[i].Target = target
		g.CutScene.Actions[i].Data = data
	}
}

func (g *Game) resolveTarget(target any) interface{} {
	if target != nil {
		id, ok := target.(string)
		if !ok {
			// Handle the case where target is not a string
			return target
		}
		// fmt.Println(id)
		switch id {
		case "player":
			return g.Player
		case "dialogue":
			return g.Dialogue
		case "music":
			return g.Music
		case "scene":
			return &g.CurrentScene
		default:
			npc1 := g.Scenes[g.CurrentScene].NPCs[id]
			// fmt.Println("Made it")
			// fmt.Println(npc1)
			return npc1
		}
	}
	return nil
}

func (g *Game) resolveData(actionData any) interface{} {
	if actionData != nil {
		switch t := actionData.(type) {
		case map[string]interface{}:
			fmt.Println("Map:", t)
			return cutscene.Vector2D{
				X: t["x"].(float64),
				Y: t["y"].(float64),
			}
		case []interface{}:
			var strings []string
			for _, item := range t {
				str, ok := item.(string)
				if !ok {
					// Handle the error or skip the item
					continue
				}
				strings = append(strings, str)
			}
			fmt.Println("Strings:", strings)
			return strings
		default:
			fmt.Println("Value:", t)
			return t
		}
	} else {
		fmt.Println("is Nil")
	}
	return nil
}
