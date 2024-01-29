package player

import (
	"image"
	"log"
	"rpg_demo/collisions"
	"rpg_demo/shared"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

type Frame struct {
	Height    int
	Width     int
	Count     int
	Current   int
	TickCount int // Counter to track the number of updates
}

type Player struct {
	SpriteSheets map[string]*ebiten.Image // Map of sprite sheets for each direction
	Direction    string
	Frame        *Frame
	Speed        float64
	X, Y         float64
	CanMove      bool
}

func New() *Player {
	return &Player{
		SpriteSheets: loadSpriteSheets(),
		Frame: &Frame{
			Height: 68,
			Width:  192 / 4,
			Count:  4,
		},
		Speed:     5,
		Direction: "down",
		X:         1000,
		Y:         1000,
		CanMove:   true,
	}
}

func (p *Player) Draw(screen *ebiten.Image, WorldWidth, WorldHeight float64) {
	currentSpriteSheet := p.SpriteSheets[p.Direction]
	// Determine the x, y location of the current frame on the sprite sheet
	sx := p.Frame.Current * p.Frame.Width
	sy := 0 // This value always remains 0

	// Create a sub-image that represents the current frame
	frame := currentSpriteSheet.SubImage(image.Rect(sx, sy, sx+p.Frame.Width, sy+p.Frame.Height)).(*ebiten.Image)

	opts := &ebiten.DrawImageOptions{}
	if p.Direction == "left" {
		opts.GeoM.Scale(-1, 1)                         // Flip horizontally
		opts.GeoM.Translate(float64(p.Frame.Width), 0) // Adjust the position after flipping
	}

	//Draw Character at the center of the screen
	var charX, charY float64
	ScreenWidth, ScreenHeight := float64(screen.Bounds().Dx()), float64(screen.Bounds().Dy())

	// Keep player centered unless near world boundaries
	if p.X > ScreenWidth/2 && p.X < WorldWidth-ScreenWidth/2 {
		charX = ScreenWidth / 2
	} else {
		// Player is near or at the horizontal boundary
		if p.X <= ScreenWidth/2 {
			charX = p.X
		} else {
			charX = ScreenWidth - (WorldWidth - p.X)
		}
	}

	if p.Y > ScreenHeight/2 && p.Y < WorldHeight-ScreenHeight/2 {
		charY = ScreenHeight / 2
	} else {
		// Player is near or at the vertical boundary
		if p.Y <= ScreenHeight/2 {
			charY = p.Y
		} else {
			charY = ScreenHeight - (WorldHeight - p.Y)
		}
	}
	opts.GeoM.Translate(charX-float64(p.Frame.Width)/2, charY-float64(p.Frame.Height)/2)
	screen.DrawImage(frame, opts)
}

func (p *Player) Update(sceneCollisions collisions.Collisions, onDoorChange func(*collisions.Door), onStateChange func(newState shared.GameState)) error {
	var newX, newY float64
	moving := false
	if p.CanMove {
		if ebiten.IsKeyPressed(ebiten.KeyLeft) {
			p.Direction = "left"
			moving = true
		}
		if ebiten.IsKeyPressed(ebiten.KeyRight) {
			p.Direction = "right"
			moving = true
		}
		if ebiten.IsKeyPressed(ebiten.KeyUp) {
			p.Direction = "up"
			moving = true
		}
		if ebiten.IsKeyPressed(ebiten.KeyDown) {
			p.Direction = "down"
			moving = true
		}
	}

	if ebiten.IsKeyPressed(ebiten.KeyR) {
		p.X = 1000
		p.Y = 1000
	}
	if moving {
		newX, newY = p.FuturePosition(p.Direction)
		if !p.Colliding(sceneCollisions.Obstacles, newX, newY) {
			p.X = newX
			p.Y = newY
		}
		colliding, door := p.CollidingWithDoor(sceneCollisions.Doors, newX, newY)
		if p.Colliding(sceneCollisions.Obstacles, newX, newY) && colliding {
			onDoorChange(door)
			onStateChange(shared.TransitionState)
		}
		// Increment the tick count
		p.Frame.TickCount++
	}

	// Update the current frame every 10 ticks
	if p.Frame.TickCount >= 10 {
		p.Frame.Current = (p.Frame.Current + 1) % p.Frame.Count
		p.Frame.TickCount = 0 // Reset the tick count
	}
	return nil
}

func (p *Player) Colliding(obstacles []*image.Rectangle, newX, newY float64) bool {
	playerRect := image.Rect(int(newX)-p.Frame.Width/2, int(newY)-p.Frame.Height/2, int(newX)+p.Frame.Width-p.Frame.Width/2, int(newY)+p.Frame.Height-p.Frame.Height/2)
	for _, obstacle := range obstacles {
		if !playerRect.Intersect(*obstacle).Empty() {
			// Collision detected
			return true
		}
	}
	// No collision
	return false
}

func (p *Player) CollidingWithDoor(doors []*collisions.Door, newX, newY float64) (bool, *collisions.Door) {
	playerRect := image.Rect(int(newX)-p.Frame.Width/2, int(newY)-p.Frame.Height/2, int(newX)+p.Frame.Width-p.Frame.Width/2, int(newY)+p.Frame.Height-p.Frame.Height/2)
	for _, door := range doors {
		if !playerRect.Intersect(*door.Rect).Empty() {
			//Collision dectected
			return true, door
		}
	}
	// No collision
	var door *collisions.Door
	return false, door
}

func (p Player) FuturePosition(dir string) (float64, float64) {
	switch dir {
	case "left":
		p.X -= p.Speed // Move left
	case "right":
		p.X += p.Speed // Move right
	case "up":
		p.Y -= p.Speed // Move up
	case "down":
		p.Y += p.Speed // Move down
	}
	return p.X, p.Y
}

func loadSpriteSheets() map[string]*ebiten.Image {
	// Create a map to hold the sprite sheets
	spriteSheets := make(map[string]*ebiten.Image)

	// List of directions
	directions := []string{"up", "down", "right", "left"}

	// Loop over the directions and load the corresponding sprite sheet
	for _, direction := range directions {
		// Construct the file path for the sprite sheet
		// This assumes you have files named like "playerUp.png", "playerDown.png", etc.
		if direction == "left" {
			spriteSheets["left"] = spriteSheets["right"]
			break
		}
		c := cases.Upper(language.English)
		path := "assets/player" + c.String(direction) + "Black.png"

		// Load the image
		img, _, err := ebitenutil.NewImageFromFile(path)
		if err != nil {
			log.Fatalf("failed to load '%s' sprite sheet: %v", direction, err)
		}

		// Store the loaded image in the map
		spriteSheets[direction] = img
	}

	return spriteSheets
}
