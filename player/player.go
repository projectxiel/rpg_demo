package player

import (
	"image"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

type Frame struct {
	Height  int
	Width   int
	Count   int
	Current int
}

type Player struct {
	SpriteSheet *ebiten.Image
	Frame       *Frame
}

func New() *Player {
	img, _, err := ebitenutil.NewImageFromFile("assets/playerDownBlack.png")
	if err != nil {
		log.Fatal(err)
	}
	return &Player{
		SpriteSheet: img,
		Frame: &Frame{
			Height: 68,
			Width:  192 / 4,
			Count:  4,
		},
	}
}

func (p *Player) Draw(screen *ebiten.Image) {

	// Determine the x, y location of the current frame on the sprite sheet
	sx := p.Frame.Current * p.Frame.Width
	sy := 0 // This value always remains 0

	// Create a sub-image that represents the current frame
	frame := p.SpriteSheet.SubImage(image.Rect(sx, sy, sx+p.Frame.Width, sy+p.Frame.Height)).(*ebiten.Image)

	opts := &ebiten.DrawImageOptions{}
	screen.DrawImage(frame, opts)
}
