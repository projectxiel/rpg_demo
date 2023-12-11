package scene

import (
	"image/color"
	"log"
	"math"
	"rpg_demo/collisions"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

type Scene struct {
	Background *ebiten.Image
	Foreground *ebiten.Image
	Width      float64
	Height     float64
	Collisions collisions.Collisions
}

func New(BgPath, FgPath, ColPath string) *Scene {
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
		Collisions: collisions.New(ColPath),
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
	drawDoors(s, screen, bgX, bgY)

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

func drawDoors(s *Scene, screen *ebiten.Image, bgX, bgY float64) {
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
