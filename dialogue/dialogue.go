package dialogue

import (
	"image/color"
	"log"
	"strings"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text"
	"golang.org/x/image/font"
	"golang.org/x/image/font/gofont/goregular"
	"golang.org/x/image/font/opentype"
	"golang.org/x/image/math/fixed"
)

type Dialogue struct {
	TextLines         []string
	CurrentLine       int
	CharIndex         int
	FramesPerChar     int // Number of frames to wait before showing the next character
	AccumulatedFrames int // Frame counter for the typewriter effect
	IsOpen            bool
	Finished          bool
	Font              font.Face
	Image             *ebiten.Image
}

func New() *Dialogue {
	font, err := loadFontFace()
	if err != nil {
		log.Fatal(err)
	}
	d := &Dialogue{
		TextLines:     []string{""},
		FramesPerChar: 2, // For example, one character every 2 frames
		IsOpen:        false,
		CurrentLine:   0,
		CharIndex:     0,
		Font:          font,
	}
	return d

}

func (d *Dialogue) Update() {
	if !d.IsOpen || d.Finished {
		return
	}

	d.AccumulatedFrames++
	if d.AccumulatedFrames >= d.FramesPerChar {
		d.AccumulatedFrames = 0
		d.CharIndex++
		if d.CharIndex > len(d.TextLines[d.CurrentLine]) {
			d.CharIndex = len(d.TextLines[d.CurrentLine])
			d.Finished = true
		}
	}
}

func (d *Dialogue) NextLine() {
	if d.CurrentLine < len(d.TextLines)-1 {
		d.CurrentLine++
		d.CharIndex = 0
		d.Finished = false
	} else {
		// No more lines, close the dialogue
		d.IsOpen = false
	}
}

func (d *Dialogue) Draw(screen *ebiten.Image) {
	if !d.IsOpen {
		return
	}

	// Set up the dialogue box dimensions
	boxWidth := screen.Bounds().Dx() - 40         // 10 pixels padding on each side
	boxHeight := 170                              // Fixed height for the dialogue box
	boxX := 20                                    // X position of the box
	boxY := screen.Bounds().Dy() - boxHeight - 20 // Y position of the box, 10 pixels above the bottom of the screen

	// Draw the dialogue box background
	dialogueBox := ebiten.NewImage(boxWidth, boxHeight)
	dialogueBox.Fill(color.Black)
	opts := &ebiten.DrawImageOptions{}
	opts.GeoM.Translate(float64(boxX), float64(boxY))
	screen.DrawImage(dialogueBox, opts)
	var err error
	if d.Image != nil {
		ImageOpts := &ebiten.DrawImageOptions{}

		ImageOpts.GeoM.Scale(0.165, 0.165)
		ImageOpts.GeoM.Translate(float64(boxX), float64(boxY))
		screen.DrawImage(d.Image, ImageOpts)
		if err != nil {
			log.Fatal(err)
		}
	}

	fontFace := d.Font
	if err != nil {
		log.Fatal(err)
	}
	var textToDisplay string

	if d.Image != nil {
		textToDisplay = wrapText(d.TextLines[d.CurrentLine][:d.CharIndex], 540, fontFace)
	} else {
		textToDisplay = wrapText(d.TextLines[d.CurrentLine][:d.CharIndex], 630, fontFace)
	}
	// Calculate the number of lines and the height of each line
	numLines := countLines(textToDisplay)
	lineHeight := fontFace.Metrics().Height.Ceil() // Or a custom line height if you prefer

	// Calculate the starting Y position for vertical centering
	totalTextHeight := numLines * lineHeight
	startY := boxY + 20 + (boxHeight-totalTextHeight)/2 // This centers the text

	// If the text height exceeds the box height, adjust the start Y position
	if totalTextHeight > boxHeight {
		startY = boxY + (totalTextHeight - boxHeight) // Adjust to move text up as it grows
	}

	// Draw the text
	if d.Image != nil {
		text.Draw(screen, textToDisplay, fontFace, boxX+200, startY, color.White)
	} else {
		text.Draw(screen, textToDisplay, fontFace, boxX+70, startY+5, color.White)
	}

}
func (d *Dialogue) IsLastLine() bool {
	return d.CurrentLine == len(d.TextLines)-1
}

func wrapText(text string, maxWidth int, face font.Face) string {
	var wrapped string
	var lineWidth fixed.Int26_6
	spaceWidth := font.MeasureString(face, " ")

	for _, word := range strings.Fields(text) {
		wordWidth := font.MeasureString(face, word)

		// If adding the new word exceeds the max width, then insert a new liFace
		if lineWidth > 0 && lineWidth+wordWidth+spaceWidth > fixed.I(maxWidth) {
			wrapped += "\n"
			lineWidth = 0
		}

		if lineWidth > 0 {
			wrapped += " "
			lineWidth += spaceWidth
		}

		wrapped += word
		lineWidth += wordWidth
	}

	return wrapped
}
func countLines(text string) int {
	return strings.Count(text, "\n") + 1
}
func loadFontFace() (font.Face, error) {
	// Read the font data
	fontBytes := goregular.TTF

	// Parse the font data
	fontParsed, err := opentype.Parse(fontBytes)
	if err != nil {
		log.Fatal(err)
	}

	// Specify the font size
	const dpi = 72
	face, err := opentype.NewFace(fontParsed, &opentype.FaceOptions{
		Size:    25,
		DPI:     dpi,
		Hinting: font.HintingNone,
	})
	if err != nil {
		log.Fatal(err)
	}

	return face, nil
}
