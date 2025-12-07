package main

import (
	"bytes"
	"image/color"
	"log"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
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

	message := widget.NewLabel("Estado")
	message.SetText("TTS Iniciado")

	pitchData := binding.BindFloat(&common.Pitch)
	pitchSlider := widget.NewSliderWithData(0.2, 2.0, pitchData)
	pitchSlider.Step = 0.1
	pitchSlider.Resize(fyne.NewSize(500, 200))
	leftSpacer := canvas.NewRectangle(color.Transparent)
	leftSpacer.SetMinSize(fyne.NewSize(20, 0))

	pitchLabel := canvas.NewText("Voice Pitch", color.White)

	left := container.NewHBox(leftSpacer, pitchLabel)

	pitchRow := container.New(
		layout.NewBorderLayout(nil, nil, left, nil),
		left,
		pitchSlider,
	)
	common.InitConnectButton(func() { twitch.GetAuthorization() })
	w.Resize(fyne.NewSize(400, 400))
	w.SetContent(
		container.NewVBox(
			pitchRow,
			container.NewCenter(common.ConnectButton),
		),
	)
	w.ShowAndRun()

}
