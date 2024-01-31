package cutscene

import (
	"fmt"
	"rpg_demo/data"
	"rpg_demo/dialogue"
	"rpg_demo/music"
	"rpg_demo/npc"
	"rpg_demo/player"
	"rpg_demo/shared"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
)

type CutsceneActionType int

const (
	MovePlayer CutsceneActionType = iota
	MoveNPC
	ShowDialogue
	TeleportNPC
	TeleportPlayer
	TurnNPC
	TurnPlayer
	FadeIn
	FadeOut
	ChangeScene
	StopMusic
	ChangeMusic
	Wait
)

// actionMap maps strings to CutsceneActionType constants
var actionMap = map[string]CutsceneActionType{
	"MovePlayer":     MovePlayer,
	"MoveNPC":        MoveNPC,
	"ShowDialogue":   ShowDialogue,
	"TeleportNPC":    TeleportNPC,
	"TeleportPlayer": TeleportPlayer,
	"TurnNPC":        TurnNPC,
	"TurnPlayer":     TurnPlayer,
	"FadeIn":         FadeIn,
	"FadeOut":        FadeOut,
	"ChangeScene":    ChangeScene,
	"StopMusic":      StopMusic,
	"ChangeMusic":    ChangeMusic,
	"Wait":           Wait,
}

type CutsceneAction struct {
	ActionType   CutsceneActionType
	Target       interface{}
	Data         interface{}
	WaitPrevious bool // Whether to wait for previous actions to complete
}
type Vector2D struct {
	X, Y float64
}
type Cutscene struct {
	Actions       []CutsceneAction
	Current       int
	ActiveActions map[int]bool // Tracks active actions by their index
	IsPlaying     bool
}

func LoadCutscenes(dataList []data.CutsceneData) map[string]*Cutscene {
	cutscenes := make(map[string]*Cutscene)
	for _, cutscene := range dataList {
		cutscenes[cutscene.ID] = New(&cutscene)
	}
	return cutscenes
}

func New(data *data.CutsceneData) *Cutscene {
	var actions []CutsceneAction

	for _, actionData := range data.Actions {
		action := CutsceneAction{
			ActionType:   getActionType(actionData.ActionType),
			WaitPrevious: actionData.WaitPrevious,
			Target:       actionData.TargetID,
			Data:         actionData.Data,
		}
		switch action.ActionType {
		case FadeOut:
			action.Data = actionData.Data
		}
		actions = append(actions, action)
	}
	cutscene := &Cutscene{
		Actions: actions,
	}
	return cutscene
}

func (c *Cutscene) Start() {
	c.Current = 0
	c.IsPlaying = true
	c.ActiveActions = make(map[int]bool)
}

func (c *Cutscene) Update(t *shared.Transition, k shared.KeyPressed) {
	if !c.IsPlaying {
		return
	}

	for i, action := range c.Actions {
		if i < c.Current && !c.ActiveActions[i] {
			// Skip completed actions
			continue
		}

		if action.WaitPrevious && i > c.Current {
			// If the action should wait for previous ones, don't start it yet
			continue
		}

		// Process the action
		completed := c.processAction(action, t, k)
		if completed {
			c.ActiveActions[i] = false // Mark action as completed
			if i == c.Current {
				c.Current++ // Move to the next action
			}
		} else {
			c.ActiveActions[i] = true // Mark action as active
		}
	}

	// Check if all actions are completed
	if c.Current >= len(c.Actions) {
		c.IsPlaying = false
		fmt.Println("DONE")
	}
}

func (c *Cutscene) processAction(action CutsceneAction, t *shared.Transition, k shared.KeyPressed) bool {
	switch action.ActionType {
	case MoveNPC:
		cnpc := action.Target.(*npc.NPC)
		destination := action.Data.(Vector2D)
		return moveTowards(cnpc, destination)
	case MovePlayer:
		p := action.Target.(*player.Player)
		destination := action.Data.(Vector2D)
		destination.X = destination.X + float64(p.Frame.Width)/2
		destination.Y = destination.Y + float64(p.Frame.Height)/2
		return moveTowards(p, destination)
	case FadeOut:
		t.Alpha += action.Data.(float64)
		f := false
		if t.Alpha >= 1.0 {
			t.Alpha = 1.0
			f = true
		}
		return f
	case FadeIn:
		t.Alpha -= action.Data.(float64)
		f := false
		if t.Alpha <= 0.0 {
			t.Alpha = 0.0
			f = true
			// The new scene is fully visible now, and game continues as normal
		}
		return f
	case TeleportPlayer:
		p := action.Target.(*player.Player)
		destination := action.Data.(Vector2D)
		p.X = destination.X + float64(p.Frame.Width)/2
		p.Y = destination.Y + float64(p.Frame.Height)/2
		return true
	case TeleportNPC:
		cnpc := action.Target.(*npc.NPC)
		destination := action.Data.(Vector2D)
		cnpc.X = destination.X
		cnpc.Y = destination.Y
		return true
	case TurnPlayer:
		p := action.Target.(*player.Player)
		dir := action.Data.(string)
		p.Direction = dir
		return true
	case TurnNPC:
		p := action.Target.(*npc.NPC)
		dir := action.Data.(string)
		p.Direction = dir
		return true
	case ShowDialogue:
		d := action.Target.(*dialogue.Dialogue)
		if !d.IsOpen {
			d.IsOpen = true
			d.CurrentLine = 0
			d.CharIndex = 0
			d.Finished = false
			d.TextLines = action.Data.([]string)
		} else {
			if ebiten.IsKeyPressed(ebiten.KeyZ) && !k.KeyZ {
				if d.Finished {
					d.NextLine()
					if !d.IsOpen {
						return true
					}
				} else {
					// Instantly display all characters in the current line
					d.CharIndex = len(d.TextLines[d.CurrentLine])
					d.Finished = true
				}

			}
			d.Update()
			// return d.Finished
		}
	case ChangeScene:
		s := action.Target.(*string)
		newScene := action.Data.(string)
		*s = newScene
		return true
	case StopMusic:
		m := action.Target.(*music.Music)
		ch := make(chan struct{})
		go m.FadeOut(time.Millisecond*500, ch)
		m.Paused = true
		return true
	case ChangeMusic:
		m := action.Target.(*music.Music)
		newSong := action.Data.(string)
		t.Music = true
		// Channel to signal when fade-out is complete
		doneChan := make(chan struct{})

		// Start fade-out in a goroutine
		go m.FadeOut(time.Second, doneChan)

		// Wait for the fade-out to complete in another goroutine
		go func() {
			<-doneChan // Wait for fade-out to complete
			// Load and play the new audio
			// Safe closure and loading of new audio
			if m.IsPlaying() || m.Paused {
				m.CloseAudio() // Ensure the current audio is closed
			}
			m.LoadAudio("./assets/" + newSong)
			m.PlayAudio()
			t.Music = false
		}()
		return true
	case Wait:
		t.Timer += 1
		target := action.Data.(int)
		if t.Timer == target {
			t.Timer = 0
			return true
		}
		fmt.Println(t.Timer)
	}
	return false
}

func moveTowards(entity interface{}, target Vector2D) bool {
	const speed = 5.0
	res := false
	switch e := entity.(type) {
	case *player.Player:
		fmt.Println(e.X, e.Y)
		if e.X > target.X {
			e.Direction = "left"
			e.X -= speed
		} else if e.X < target.X {
			e.Direction = "right"
			e.X += speed
		} else if e.Y < target.Y {
			e.Direction = "down"
			e.Y += speed
		} else if e.Y > target.Y {
			e.Direction = "up"
			e.Y -= speed
		} else {
			e.Frame.Count = 2
			res = true
			fmt.Println("Achieved")
		}
		e.Frame.TickCount++
		if e.Frame.TickCount >= 10 && !res {
			e.Frame.Current = (e.Frame.Current + 1) % e.Frame.Count
			e.Frame.TickCount = 0 // Reset the tick count
		}
		fmt.Println(e.Direction)
	case *npc.NPC:
		fmt.Println(e.X, e.Y)
		if e.X > target.X {
			e.Direction = "left"
			e.X -= speed
		} else if e.X < target.X {
			e.Direction = "right"
			e.X += speed
		} else if e.Y < target.Y {
			e.Direction = "down"
			e.Y += speed
		} else if e.Y > target.Y {
			e.Direction = "up"
			e.Y -= speed
		} else {
			// e.CurrentFrame = 2
			res = true
			fmt.Println("Achieved")
		}
		e.Frame.TickCount++
		if e.Frame.TickCount >= 10 && !res {
			e.Frame.Current = (e.Frame.Current + 1) % e.Frame.Count
			e.Frame.TickCount = 0 // Reset the tick count
		}
		fmt.Println(e.Direction)
	default:
		res = true
	}

	return res
}

// getActionType returns the CutsceneActionType for a given string
func getActionType(actionType string) CutsceneActionType {
	if val, ok := actionMap[actionType]; ok {
		return val
	}
	// Handle the case where the actionType is not found
	fmt.Println("Invalid action type:", actionType)
	return -1
}
