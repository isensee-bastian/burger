package burger

import (
	"bytes"
	"fmt"
	"github.com/hajimehoshi/ebiten/v2/audio"
	"github.com/hajimehoshi/ebiten/v2/audio/mp3"
	"log"
	"os"
)

const (
	soundBasePath = "resources/sounds/"
)

func newMp3Player(fileName string) *audio.Player {
	sellSoundRaw, err := os.ReadFile(fmt.Sprintf("%s/%s", soundBasePath, fileName))
	if err != nil {
		log.Fatalf("Failed to read sound file: %v", err)
	}

	sellSound, err := mp3.DecodeF32(bytes.NewReader(sellSoundRaw))
	if err != nil {
		log.Fatalf("Failed to decode raw sound as mp3: %v", err)
	}

	audioContext := audio.NewContext(24000)
	sellPlayer, err := audioContext.NewPlayerF32(sellSound)
	if err != nil {
		log.Fatalf("Failed to create mp3 audio player: %v", err)
	}

	return sellPlayer
}
