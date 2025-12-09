package web

import (
	"context"
	"embed"
	"log"
	"net/http"
	"time"

	"ac-tts/internal/logging"
	"ac-tts/internal/twitch"
)

//go:embed static/*
var staticFiles embed.FS

var srv *http.Server

func StartWebServer() {
	mux := http.NewServeMux()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		data, err := staticFiles.ReadFile("static/index.html")
		if err != nil {
			logging.CreateLog("server - couldn't find index.html", err)
			http.Error(w, "index.html not found", http.StatusNotFound)
			return
		}
		w.Write(data)
	})

	mux.HandleFunc("/access-token", handler)

	srv = &http.Server{
		Addr:    ":9000",
		Handler: mux,
	}

	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatal("Error initializing server...", err)
	}

}

func handler(w http.ResponseWriter, r *http.Request) {
	token := r.URL.Query().Get("access_token")

	if token == "" {
		log.Fatal("No token")
	} else {
		twitch.SubscribeToChat(token)
	}

	go func() {
		time.Sleep(500 * time.Millisecond)
		if err := srv.Shutdown(context.Background()); err != nil {
			logging.CreateLog("couldn't shutdown http server", err)
			log.Fatal("Error shutting down server..", err)
		}
	}()

}
