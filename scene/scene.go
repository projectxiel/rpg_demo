package scene

import (
	"image/color"
	"log"
	"math"
	"rpg_demo/collisions"
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
	}
}

func (s *Scene) Draw(screen, img *ebiten.Image, playerX, playerY float64) {
	ScreenWidth := float64(screen.Bounds().Dx())
	ScreenHeight := float64(screen.Bounds().Dy())

	bgX, bgY := -playerX+ScreenWidth/2, -playerY+ScreenHeight/2

	// Constrain background position to world boundaries
	bgX = math.Min(math.Max(bgX, -s.Width+ScreenWidth), 0)
	bgY = math.Min(math.Max(bgY, -s.Height+ScreenHeight), 0)

	opts := &ebiten.DrawImageOptions{}
	opts.GeoM.Translate(bgX, bgY)
	screen.DrawImage(img, opts)
	// drawCollisions(s, screen, bgX, bgY)
	// drawDoors(s, screen, bgX, bgY)
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
			if ebiten.IsKeyPressed(ebiten.KeyZ) && !PressedLastFrame {
				if npc1.InteractionState == npc.NoInteraction {
					npc1.InteractionState = npc.PlayerInteracted
					player.CanMove = false // Disallow player movement
				} else if npc1.InteractionState == npc.WaitingForPlayerToResume && dial.Finished && dial.IsLastLine() {
					npc1.InteractionState = npc.NoInteraction
					dial.Image = nil
					player.CanMove = true // Allow player movement
				}

				if !dial.IsOpen {
					dial.Image = npc1.Image
					dial.IsOpen = true
					dial.CurrentLine = 0
					dial.CharIndex = 0
					dial.Finished = false
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
func DrawCollisions(s *Scene, screen *ebiten.Image, bgX, bgY float64) {
	for _, obstacle := range s.Collisions.Obstacles {
		// Translate the obstacle's position based on the background position
		obstacleOpts := &ebiten.DrawImageOptions{}
		obstacleImage := ebiten.NewImage(obstacle.Dx(), obstacle.Dy())
		obstacleColor := color.RGBA{255, 0, 0, 80} // Semi-transparent red color
		obstacleOpts.GeoM.Translate(bgX+float64(obstacle.Min.X), bgY+float64(obstacle.Min.Y))
		// Create a colored rectangle to represent the obstacle

		obstacleImage.Fill(obstacleColor)

		// Draw the obstacle image
		screen.DrawImage(obstacleImage, obstacleOpts)

		// Dispose of the obstacle image to avoid memory leaks if you're done with it
		obstacleImage.Dispose()
	}
}

func DrawDoors(s *Scene, screen *ebiten.Image, bgX, bgY float64) {
	for _, door := range s.Collisions.Doors {
		// Translate the obstacle's position based on the background position
		doorOpts := &ebiten.DrawImageOptions{}
		doorImage := ebiten.NewImage(door.Rect.Dx(), door.Rect.Dy())
		doorColor := color.RGBA{0, 0, 255, 80} // Semi-transparent blue color
		doorOpts.GeoM.Translate(bgX+float64(door.Rect.Min.X), bgY+float64(door.Rect.Min.Y))
		// Create a colored rectangle to represent the door

		doorImage.Fill(doorColor)

		// Draw the door image
		screen.DrawImage(doorImage, doorOpts)

		// Dispose of the door image to avoid memory leaks if you're done with it
		doorImage.Dispose()
	}
}
