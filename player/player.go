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
	X, Y        float64
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
		X: 100,
		Y: 100,
	}
}

func (p *Player) Draw(screen *ebiten.Image) {

	// Determine the x, y location of the current frame on the sprite sheet
	sx := p.Frame.Current * p.Frame.Width
	sy := 0 // This value always remains 0

	// Create a sub-image that represents the current frame
	frame := p.SpriteSheet.SubImage(image.Rect(sx, sy, sx+p.Frame.Width, sy+p.Frame.Height)).(*ebiten.Image)

	opts := &ebiten.DrawImageOptions{}
	opts.GeoM.Translate(p.X, p.Y)
	screen.DrawImage(frame, opts)
}

func (p *Player) Update() error {
	if ebiten.IsKeyPressed(ebiten.KeyLeft) {
		p.X -= 2 // Move left
	}
	if ebiten.IsKeyPressed(ebiten.KeyRight) {
		p.X += 2 // Move right
	}
	if ebiten.IsKeyPressed(ebiten.KeyUp) {
		p.Y -= 2 // Move up
	}
	if ebiten.IsKeyPressed(ebiten.KeyDown) {
		p.Y += 2 // Move down
	}
	return nil
}
