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
	"ac-tts/internal/github"
	"ac-tts/internal/logging"
	"ac-tts/internal/reproductor"
	"ac-tts/internal/twitch"
	"ac-tts/internal/web"
)

var Version = "dev"

func main() {
	github.GetLatestReleaseVersion()
	a := app.New()
	w := a.NewWindow("AC - Text to Speech for Twitch :)")

	go func() {
		web.StartWebServer()
	}()

	ani, err := animalese.NewAnimaleseFromBytes(assets.AnimaleseWav, func() {

	})
	if err != nil {
		logging.CreateLog(err)
		log.Fatal(err)
		panic(err)
	}

	wave := ani.AnimaleseFunc("test", true, 1.0)
	streamer, format, err := wav.Decode(bytes.NewReader(wave))
	if err != nil {
		logging.CreateLog(err)
		log.Fatal(err)
	}
	defer streamer.Close()
	reproductor.InitSpeaker(format)

	pitchData := binding.BindFloat(&common.Pitch)
	common.InitPitchRow(pitchData)
	common.InitTestPitchButton(func() {
		reproductor.Reproduce("Hola esto es una prueba de pitch :)", "")
	})

	common.InitConnectButton(func() { twitch.GetAuthorization() })
	w.Resize(fyne.NewSize(400, 200))

	icon := fyne.NewStaticResource("icon.png", assets.Icon)

	common.InitKofiButton()

	content := container.NewVBox(
		common.PitchRow,
		container.NewCenter(container.NewHBox(common.TestPitchButton, common.ConnectButton)),
	)

	footer := container.NewCenter(common.KofiButton)
	w.SetIcon(icon)
	w.SetContent(
		container.New(
			layout.NewBorderLayout(nil, footer, nil, nil),
			footer,
			container.NewStack(content),
		),
	)
	w.ShowAndRun()

}
