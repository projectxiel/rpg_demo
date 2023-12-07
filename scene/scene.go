package scene

import (
	"log"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

type Scene struct {
	Background *ebiten.Image
	Foreground *ebiten.Image
	Width      float64
	Height     float64
}

func New(BgPath, FgPath string) *Scene {
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
}
