package reproductor

import (
	"bytes"
	"log"
	"math/rand"
	"sync"
	"time"

	"github.com/faiface/beep"
	"github.com/faiface/beep/effects"
	"github.com/faiface/beep/speaker"
	"github.com/faiface/beep/wav"

	"ac-tts/internal/animalese"
	"ac-tts/internal/assets"
	"ac-tts/internal/common"
)

var UserMap = make(map[string]float64)
var once sync.Once

func Reproduce(text string, user string) {
	ani, err := animalese.NewAnimaleseFromBytes(assets.AnimaleseWav, func() {

	})

	if err != nil {
		log.Fatal(err)
		panic(err)
	}
	wave := ani.AnimaleseFunc(text, true, common.Pitch)

	streamer, format, err := wav.Decode(bytes.NewReader(wave))
	if err != nil {
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

func GetPitchForUser(user string) float64 {
	if value, ok := UserMap[user]; ok {

		return value
	}
	min := 0.2
	max := 2.0
	step := 0.1
	steps := int((max - min) / step)
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	idx := r.Intn(steps + 1)
	pitch := min + float64(idx)*step
	UserMap[user] = pitch
	return pitch
}

func InitSpeaker(format beep.Format) {
	once.Do(func() {
		speaker.Init(format.SampleRate, format.SampleRate.N(time.Second/10))
	})
}
