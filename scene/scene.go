package scene

import (
	"log"
	"math"
	"rpg_demo/ability"
	"rpg_demo/collisions"
	"rpg_demo/cutscene"
	"rpg_demo/data"
	"rpg_demo/dialogue"
	"rpg_demo/npc"
	"rpg_demo/player"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

type Scene struct {
	Background *ebiten.Image
	Foreground *ebiten.Image
	Width      float64
	Height     float64
	Collisions collisions.Collisions
	Music      string
	NPCs       map[string]*npc.NPC
	Cutscenes  map[string]*cutscene.Cutscene
	X, Y       float64
}

func New(name string) *Scene {
	prefix := "assets/"
	JsonPath := prefix + name + ".json"
	data := data.LoadJsonFile(JsonPath)

	BgPath := prefix + name + ".png"
	FgPath := prefix + name + "Fore.png"

	// Load the background
	Bg, _, err := ebitenutil.NewImageFromFile(BgPath)
	if err != nil {
		log.Fatal(err)
	}
	// Load the foreground
	Fg, _, err := ebitenutil.NewImageFromFile(FgPath)
	if err != nil {
		log.Fatal(err)
	}
	return &Scene{
		Background: Bg,
		Foreground: Fg,
		Width:      float64(Bg.Bounds().Dx()),
		Height:     float64(Bg.Bounds().Dy()),
		Collisions: collisions.New(data),
		Music:      data.Music,
		NPCs:       npc.LoadNPCs(data.NPCs),
		Cutscenes:  cutscene.LoadCutscenes(data.Cutscenes),
	}
}

func (s *Scene) Draw(screen, img *ebiten.Image, p *player.Player) {
	ScreenWidth := float64(screen.Bounds().Dx())
	ScreenHeight := float64(screen.Bounds().Dy())

	bgX, bgY := -p.X+ScreenWidth/2, -p.Y+ScreenHeight/2

	// Constrain background position to world boundaries
	bgX = math.Min(math.Max(bgX, -s.Width+ScreenWidth), 0)
	bgY = math.Min(math.Max(bgY, -s.Height+ScreenHeight), 0)

	opts := &ebiten.DrawImageOptions{}
	opts.GeoM.Translate(bgX, bgY)
	if p.Ability.Type == ability.StopTime && p.Ability.Activated {
		opts.ColorScale.Scale(.5, .5, .5, 1)
	}
	screen.DrawImage(img, opts)
	s.X, s.Y = bgX, bgY
}
func (s *Scene) Update() {
	for _, npc := range s.NPCs {
		npc.Update()
	}
}
func (s *Scene) DrawNPCs(screen *ebiten.Image) {
	for _, npc := range s.NPCs {
		npc.Draw(screen, s.X, s.Y)
	}
}
func (s *Scene) HandleNPCInteractions(player *player.Player, PressedLastFrame bool, dial *dialogue.Dialogue) {
	playerX, playerY := player.X-float64(player.Frame.Width)/2, player.Y-float64(player.Frame.Height)/2
	for _, npc1 := range s.NPCs {
		if npc1.Near(playerX, playerY) {
			if ebiten.IsKeyPressed(ebiten.KeyZ) && !PressedLastFrame && npc1.IsTalker() {
				if npc1.InteractionState == npc.NoInteraction {
					npc1.ChangeDirection(playerX, playerY)
					npc1.InteractionState = npc.PlayerInteracted
					player.CanMove = false // Disallow player movement
				} else if npc1.InteractionState == npc.WaitingForPlayerToResume && dial.Finished && dial.IsLastLine() {
					npc1.InteractionState = npc.NoInteraction
					dial.Image = nil
					player.CanMove = true // Allow player movement
				}

				if !dial.IsOpen {
					dial.Image = npc1.Image
					dial.Speaker = npc1.Name
					dial.OpenAndReset()
					dial.TextLines = npc1.Behaviors["talker"].Value()
				} else {
					if dial.Finished {
						dial.NextLine()
					} else {
						// Instantly display all characters in the current line
						dial.CharIndex = len(dial.TextLines[dial.CurrentLine])
						dial.Finished = true
					}
				}
			}
		}
	}
}
