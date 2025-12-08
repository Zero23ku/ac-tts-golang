package main

import (
	"bytes"
	"log"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/layout"
	"github.com/faiface/beep/wav"

	"ac-tts/internal/animalese"
	"ac-tts/internal/assets"
	"ac-tts/internal/common"
	"ac-tts/internal/reproductor"
	"ac-tts/internal/twitch"
	"ac-tts/internal/web"
)

func main() {

	a := app.New()
	w := a.NewWindow("AC - Text to Speech for Twitch :)")

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
	w.Resize(fyne.NewSize(400, 200))

	res, _ := fyne.LoadResourceFromPath("cup-border.png")
	common.InitKofiButton(res)

	content := container.NewVBox(
		common.PitchRow,
		container.NewCenter(common.ConnectButton),
	)

	footer := container.NewCenter(common.KofiButton)

	w.SetContent(
		container.New(
			layout.NewBorderLayout(nil, footer, nil, nil),
			footer,
			container.NewStack(content),
		),
	)
	w.ShowAndRun()

}
