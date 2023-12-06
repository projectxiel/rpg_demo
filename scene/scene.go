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

func (s *Scene) DrawBg(screen *ebiten.Image) {
	opts := getDrawOptions()
	screen.DrawImage(s.Background, opts)
}

func (s *Scene) DrawFg(screen *ebiten.Image) {
	opts := getDrawOptions()
	screen.DrawImage(s.Foreground, opts)
}

func getDrawOptions() *ebiten.DrawImageOptions {
	opts := &ebiten.DrawImageOptions{}
	return opts
}
