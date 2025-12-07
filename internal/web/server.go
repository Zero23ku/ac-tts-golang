package web

import (
	"embed"
	"log"
	"net/http"

	"ac-tts/internal/twitch"
)

//go:embed static/*
var staticFiles embed.FS

func StartWebServer() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		data, err := staticFiles.ReadFile("static/index.html")
		if err != nil {
			http.Error(w, "index.html no encontrado", http.StatusInternalServerError)
			return
		}
		w.Write(data)
	})
	http.HandleFunc("/access-token", handler)
	if err := http.ListenAndServe(":9000", nil); err != nil {
		log.Fatal("Error")
	}
}

func handler(w http.ResponseWriter, r *http.Request) {
	token := r.URL.Query().Get("access_token")

	if token == "" {
		log.Fatal("No token")
	} else {
		twitch.SubscribeToChat(token)
	}
}
