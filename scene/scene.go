package scene

import (
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

type Scene struct {
	Background *ebiten.Image
	Foreground *ebiten.Image
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
	}
}

func (s *Scene) DrawBg(screen *ebiten.Image, playerX, playerY float64) {
	bgX, bgY := -playerX, -playerY
	opts := getDrawOptions(bgX, bgY)
	screen.DrawImage(s.Background, opts)
}

func (s *Scene) DrawFg(screen *ebiten.Image, playerX, playerY float64) {
	fgX, fgY := -playerX, -playerY
	opts := getDrawOptions(fgX, fgY)
	screen.DrawImage(s.Foreground, opts)
}

func getDrawOptions(X, Y float64) *ebiten.DrawImageOptions {
	opts := &ebiten.DrawImageOptions{}
	opts.GeoM.Translate(X, Y)
	return opts
}
