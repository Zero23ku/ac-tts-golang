package main

import (
	"bytes"
	"log"

	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/faiface/beep/wav"

	"ac-tts/internal/animalese"
	"ac-tts/internal/assets"
	"ac-tts/internal/reproductor"
	"ac-tts/internal/twitch"
	"ac-tts/internal/web"
)

func main() {

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

	twitch.GetAuthorization()
	// Crear la aplicación
	myApp := app.New()
	myWindow := myApp.NewWindow("Hola Fyne")

	// Crear un botón
	button := widget.NewButton("Haz clic aquí", func() {
		myWindow.SetContent(widget.NewLabel("¡Botón presionado!"))
	})

	// Colocar el botón en la ventana
	myWindow.SetContent(container.NewVBox(
		widget.NewLabel("Bienvenido a tu primera app con Fyne"),
		button,
	))

	// Mostrar la ventana
	myWindow.ShowAndRun()

}
