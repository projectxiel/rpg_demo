package music

import (
	"os"
	"strings"
	"time"

	"github.com/hajimehoshi/ebiten/v2/audio"
	"github.com/hajimehoshi/ebiten/v2/audio/mp3"
	"github.com/hajimehoshi/ebiten/v2/audio/wav"
)

type Music struct {
	audioContext *audio.Context
	player       *audio.Player
	audioFile    *os.File // Add a field to store the audio file
	CurrentSong  string
	Paused       bool
}

const sampleRate = 44100

func (m *Music) LoadAudio(filePath string) error {
	m.CurrentSong = filePath
	parts := strings.Split(filePath, ".")

	var err error
	m.audioFile, err = os.Open(filePath)
	if err != nil {
		return err
	}
	var d *mp3.Stream
	var dw *wav.Stream
	switch parts[2] {
	case "mp3":
		// Use DecodeWithSampleRate for decoding
		d, err = mp3.DecodeWithSampleRate(sampleRate, m.audioFile)
		if err != nil {
			m.audioFile.Close()
			return err
		}
		// Use the new method to create a player
		m.player, err = m.audioContext.NewPlayer(d)
		if err != nil {
			return err
		}
	case "wav":
		dw, err = wav.DecodeWithSampleRate(sampleRate, m.audioFile)
		if err != nil {
			m.audioFile.Close()
			return err
		}
		m.player, err = m.audioContext.NewPlayer(dw)
		if err != nil {
			return err
		}
	}

	return nil
}
func (m *Music) FadeIn(duration time.Duration, doneChan chan struct{}) {
	// Determine the amount of time to sleep between volume adjustments
	const steps = 30
	sleepDuration := duration / steps

	// Start with volume at 0
	m.player.SetVolume(0)
	// Gradually increase the volume
	for i := 0; i < steps; i++ {
		// Calculate the new volume (linearly increases)
		newVolume := float64(i+1) / float64(steps)
		m.player.SetVolume(newVolume)
		time.Sleep(sleepDuration)
	}

	// Ensure the volume is set to the maximum at the end
	m.player.SetVolume(1)
	// Resume playing if the player was paused
	m.player.Play()
	m.Paused = false
	// Signal that the fade-in is complete
	close(doneChan)
}
func (m *Music) FadeOut(duration time.Duration, doneChan chan struct{}) {
	// Determine the amount of time to sleep between volume adjustments
	const steps = 30
	sleepDuration := duration / steps

	// Gradually decrease the volume
	for i := 0; i < steps; i++ {
		// Calculate the new volume (linearly decreases)
		newVolume := float64(steps-i-1) / float64(steps)
		m.player.SetVolume(newVolume)
		time.Sleep(sleepDuration)
	}

	// Ensure the volume is set to 0 at the end
	m.player.SetVolume(0)
	// Stop the player if needed
	m.player.Pause()
	m.Paused = true
	// Signal that the fade-out is complete
	close(doneChan)
}
func (m *Music) GetPlayer() *audio.Player {
	return m.player
}
func (m *Music) IsPlaying() bool {
	return m.player.IsPlaying()
}
func (m *Music) Pause() {
	m.player.Pause()
	m.Paused = true
}
func (m *Music) RewindMusic() {
	m.player.Rewind()
	m.player.Play()
	m.Paused = false
}
func (m *Music) SetCtx(auctx *audio.Context) {
	m.audioContext = auctx
}
func (m *Music) IsEmpty() bool {
	return m.audioFile == nil
}

func (m *Music) PlayAudio() {
	if m.player != nil {
		m.player.Play()
	}
	m.Paused = false
}

func (m *Music) CloseAudio() {
	if m.player != nil {
		m.player.Close()
	}
	if m.audioFile != nil {
		m.audioFile.Close()
	}
}
