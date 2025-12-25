package tiktok

import (
	"fmt"
	"log"

	"github.com/steampoweredtaco/gotiktoklive"
)

func InitTiktok(username string) {
	fmt.Println("Iniciando tiktok...")
	tiktok, err := gotiktoklive.NewTikTok()
	if err != nil {
		fmt.Println("Error 1: ", err)
		log.Fatal(err)
	}
	fmt.Println("Capturando chat: ", username)
	live, err := tiktok.TrackUser(username)

	if err != nil {
		fmt.Println("Error 2: ", err)
		log.Fatal(err)
	}

	fmt.Println("leyendo eventos del chat")
	for event := range live.Events {
		switch e := event.(type) {
		case gotiktoklive.ChatEvent:
			fmt.Println(gotiktoklive.ChatEvent(e).Comment)
		}

	}

}
