package main

import (
	"bytes"
	"log"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"github.com/faiface/beep/wav"

	"ac-tts/internal/animalese"
	"ac-tts/internal/assets"
	"ac-tts/internal/common"
	"ac-tts/internal/reproductor"
	"ac-tts/internal/twitch"
	"ac-tts/internal/web"
)

func main() {

	// Crear la aplicación
	a := app.New()
	w := a.NewWindow("AC - Text to Speech :)")

	go func() {
		web.StartWebServer()
	}()

	ani, err := animalese.NewAnimaleseFromBytes(assets.AnimaleseWav, func() {

	})
	if err != nil {
		log.Fatal(err)
		panic(err)
	}

	wave := ani.AnimaleseFunc("test", true, 1.0)
	streamer, format, err := wav.Decode(bytes.NewReader(wave))
	if err != nil {
		log.Fatal(err)
	}
	defer streamer.Close()
	reproductor.InitSpeaker(format)

	pitchData := binding.BindFloat(&common.Pitch)
	common.InitPitchRow(pitchData)

	common.InitConnectButton(func() { twitch.GetAuthorization() })
	w.Resize(fyne.NewSize(400, 400))
	w.SetContent(
		container.NewVBox(
			common.PitchRow,
			container.NewCenter(common.ConnectButton),
		),
	)
	w.ShowAndRun()

}
