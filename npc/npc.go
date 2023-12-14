package npc

import (
	"fmt"
	"image"
	"log"
	"math"
	"rpg_demo/data"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

type InteractionState int

const (
	NoInteraction InteractionState = iota
	PlayerInteracted
	WaitingForPlayerToResume
	CutSceneInteraction
)

type Behavior interface {
	Execute(*NPC)
}

type Frame struct {
	Height    int
	Width     int
	Count     int
	Current   int
	TickCount int // Counter to track the number of updates
}
type Timer struct {
	MoveTimer    int
	StopTimer    int
	IsStopped    bool
	StopDuration int
}

type Walker struct {
	Direction string
	Speed     float64
	Timer     *Timer
}
type Talker struct {
	Placeholder string
}

type NPC struct {
	Name             string
	SpriteSheets     map[string]*ebiten.Image // Map of sprite sheets for each direction
	Direction        string
	Frame            *Frame
	X, Y             float64
	Behaviors        []Behavior
	InteractionState InteractionState
}

func (n *NPC) Draw(screen *ebiten.Image, bgX, bgY float64) {
	currentSpriteSheet := n.SpriteSheets[n.Direction]
	// Determine the x, y location of the current frame on the sprite sheet
	sx := n.Frame.Current * n.Frame.Width
	sy := 0 // This value always remains 0

	// Create a sub-image that represents the current frame
	frame := currentSpriteSheet.SubImage(image.Rect(sx, sy, sx+n.Frame.Width, sy+n.Frame.Height)).(*ebiten.Image)

	opts := &ebiten.DrawImageOptions{}
	if n.Direction == "left" {
		opts.GeoM.Scale(-1, 1)                         // Flip horizontally
		opts.GeoM.Translate(float64(n.Frame.Width), 0) // Adjust the position after flipping
	}
	opts.GeoM.Translate(bgX+n.X, bgY+n.Y)
	screen.DrawImage(frame, opts)
}
func (npc *NPC) Update() {
	for _, behavior := range npc.Behaviors {
		behavior.Execute(npc)
	}
}
func New(data *data.NPCData) *NPC {
	sheets := loadSpriteSheets(data)
	direction, sheet := GetAnySpriteSheet(sheets)
	npc := &NPC{
		Name:         data.Name,
		SpriteSheets: sheets,
		Frame: &Frame{
			Height: sheet.Bounds().Dy(),
			Width:  sheet.Bounds().Dx() / data.FrameCount,
			Count:  data.FrameCount,
		},
		Direction: direction,
		X:         data.X,
		Y:         data.Y,
		Behaviors: loadBehaviors(data),
	}
	return npc
}

func loadBehaviors(data *data.NPCData) []Behavior {
	behaviors := make([]Behavior, 0)
	// Initialize behaviors
	for _, behaviorData := range data.Behaviors {
		switch behaviorData.Type {
		case "walker":
			if speed, ok := behaviorData.Details["speed"].(float64); ok {
				direction := behaviorData.Details["direction"].(string)
				timerData, ok := behaviorData.Details["timer"].(map[string]interface{})
				timer := &Timer{}
				if ok {

					if moveTimer, ok := timerData["moveTimer"].(float64); ok {
						timer.MoveTimer = int(moveTimer)
					}
					if stopTimer, ok := timerData["stopTimer"].(float64); ok {
						timer.StopTimer = int(stopTimer)
					}
					if isStopped, ok := timerData["isStopped"].(bool); ok {
						timer.IsStopped = isStopped
					}
					if stopDuration, ok := timerData["stopDuration"].(float64); ok {
						timer.StopDuration = int(stopDuration)
					}
				}

				fmt.Println(timer)
				behaviors = append(behaviors, &Walker{Direction: direction, Speed: speed, Timer: timer})
			}
			// Add cases for other behaviors like "talker", etc.
		}
	}
	return behaviors
}

func LoadNPCs(dataList []data.NPCData) map[string]*NPC {
	npcs := make(map[string]*NPC)
	for _, npc := range dataList {
		npcs[npc.Name] = New(&npc)
	}

	return npcs
}

func GetAnySpriteSheet(spriteSheets map[string]*ebiten.Image) (string, *ebiten.Image) {
	for key, sheet := range spriteSheets {
		return key, sheet // Return the first spritesheet found
	}
	return "", nil // Return nil if the map is empty
}

func loadSpriteSheets(data *data.NPCData) map[string]*ebiten.Image {
	sheets := make(map[string]*ebiten.Image)
	// Load sprite sheets
	for direction, path := range data.SpriteSheets {
		// Load the image for the given path and store it in SpriteSheets
		// Assuming a function LoadSpriteSheet(path string) (*ebiten.Image, error)
		spriteSheet, err := LoadSpriteSheet(path)
		if err != nil {
			log.Printf("Error loading sprite sheet: %s", err)
			continue
		}
		sheets[direction] = spriteSheet
	}
	return sheets
}

func LoadSpriteSheet(path string) (*ebiten.Image, error) {
	img, _, err := ebitenutil.NewImageFromFile("./assets/" + path)

	return img, err
}
func (t *Talker) Execute(npc *NPC) {
	fmt.Println(t.Placeholder)
}
func (w *Walker) Execute(npc *NPC) {
	// Check for interaction key press to change the NPC's state
	if ebiten.IsKeyPressed(ebiten.KeyZ) {
		if npc.InteractionState == PlayerInteracted {
			npc.InteractionState = WaitingForPlayerToResume
		}
	}
	// NPC movement logic
	if npc.InteractionState == NoInteraction {
		if w.Timer.IsStopped {
			// w.Timer is stopped, so we might count down the stop timer
			w.Timer.StopTimer--
			if w.Timer.StopTimer <= 0 {
				// Time to move again
				w.Timer.IsStopped = false
				// Reset the move timer to some value
				w.Timer.MoveTimer = 60
				// Change direction
				if w.Direction == "right" {
					w.Direction = "left"
				} else {
					w.Direction = "right"
				}
			}
		} else {
			w.Timer.MoveTimer--
			w.Move(npc)
			npc.Direction = w.Direction
			if w.Timer.MoveTimer <= 0 {
				// Time to stop
				w.Timer.IsStopped = true
				// Reset the stop timer to the duration of the stop
				w.Timer.StopTimer = w.Timer.StopDuration

			}
		}
	}
	// Update the current frame every 10 ticks
	if npc.Frame.TickCount >= 10 {
		npc.Frame.Current = (npc.Frame.Current + 1) % npc.Frame.Count
		npc.Frame.TickCount = 0 // Reset the tick count
	}
}
func (w *Walker) Move(npc *NPC) {
	switch w.Direction {
	case "left":
		npc.X -= w.Speed
	case "right":
		npc.X += w.Speed
	case "up":
		npc.Y -= w.Speed
	case "down":
		npc.Y += w.Speed
	}
	npc.Frame.TickCount++
}

func (npc *NPC) IsTalkerAndWalker() bool {
	hasWalker := false
	hasTalker := false

	for _, behavior := range npc.Behaviors {
		switch behavior.(type) {
		case *Walker:
			hasWalker = true
		case *Talker:
			hasTalker = true
		}
	}

	return hasWalker && hasTalker
}

func (npc *NPC) Near(playerX, playerY float64) bool {
	return math.Abs(playerX-npc.X) < 50 && math.Abs(playerY-npc.Y) < 50
}
