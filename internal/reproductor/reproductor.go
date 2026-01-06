package reproductor

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"sync"
	"time"

	"github.com/faiface/beep"
	"github.com/faiface/beep/effects"
	"github.com/faiface/beep/speaker"
	"github.com/faiface/beep/wav"
	"github.com/fiam/gounidecode/unidecode"

	"ac-tts/internal/animalese"
	"ac-tts/internal/assets"
	"ac-tts/internal/common"
	"ac-tts/internal/logging"
)

var UserMap = make(map[string]float64)
var once sync.Once

func SaveAsWav(text string, fileName string) {
	ani, err := animalese.NewAnimaleseFromBytes(assets.AnimaleseWav, func() {

	})

	if err != nil {
		logging.CreateLog("reproductor - initializing animaleses", err)
		log.Fatal(err)
		panic(err)
	}
	text = removeAccents(text)
	wave := ani.AnimaleseFunc(text, true, common.Pitch)

	streamer, format, err := wav.Decode(bytes.NewReader(wave))
	if err != nil {
		logging.CreateLog("reproductor - creating wav", err)
		log.Fatal(err)
	}
	defer streamer.Close()

	file, err := os.Create(fileName)
	if err != nil {
		logging.CreateLog("reproductor - error creating new file", err)
		log.Fatal(err)
	}
	defer file.Close()
	fmt.Println("Guardando archivo: " + fileName)
	if err := wav.Encode(file, streamer, format); err != nil {
		log.Fatal(err)
	}
}

func Reproduce(text string, user string, pitch float64, isRandomPitch bool) {
	ani, err := animalese.NewAnimaleseFromBytes(assets.AnimaleseWav, func() {

	})

	if err != nil {
		logging.CreateLog("reproductor - initializing animaleses", err)
		log.Fatal(err)
		panic(err)
	}
	text = removeAccents(text)
	var actualPitch float64
	if isRandomPitch {
		actualPitch = getPitchForUser(user, pitch)
	} else {
		actualPitch = pitch
	}

	wave := ani.AnimaleseFunc(text, true, actualPitch)

	streamer, format, err := wav.Decode(bytes.NewReader(wave))
	if err != nil {
		logging.CreateLog("reproductor - creating wav", err)
		log.Fatal(err)
	}
	defer streamer.Close()

	InitSpeaker(format)

	ctrl := &effects.Volume{
		Streamer: streamer,
		Base:     2,
		Volume:   -2,
		Silent:   false,
	}

	done := make(chan bool)
	streamerWithDone := beep.Seq(ctrl, beep.Callback(func() {
		done <- true
	}))

	speaker.Play(streamerWithDone)
	<-done
}

func getPitchForUser(user string, pitch float64) float64 {
	if value, ok := UserMap[user]; ok {

		return value
	}
	UserMap[user] = pitch
	return pitch
}

func InitSpeaker(format beep.Format) {
	once.Do(func() {
		speaker.Init(format.SampleRate, format.SampleRate.N(time.Second/10))
	})
}

func removeAccents(text string) string {
	return unidecode.Unidecode(text) // One-liner!
}
