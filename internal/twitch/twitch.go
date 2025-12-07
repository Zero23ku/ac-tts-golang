package twitch

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net"
	"net/http"
	"os/exec"
	"runtime"
	"strings"

	"ac-tts/internal/common"
	"ac-tts/internal/reproductor"
)

const TWITCH_URL = "https://id.twitch.tv/oauth2/authorize" +
	"?response_type=token" +
	"&client_id=4u4v1h8d2yfvftoqtstu0pley1pooo" +
	"&redirect_uri=http://localhost:9000" +
	"&scope=chat:read+chat:edit" +
	"&state=c3ab8aa609ea11e793ae92361f002671"

const TWITCH_BROADCASTER_ID = "https://api.twitch.tv/helix/users"

const CLIENT_ID = "4u4v1h8d2yfvftoqtstu0pley1pooo"

const IRC_TWITCH_SERVER = "irc.chat.twitch.tv:6667"

func GetAuthorization() {
	fmt.Println("test")
	var err error
	switch runtime.GOOS {
	case "linux":
		err = exec.Command("xdg-open", TWITCH_URL).Start()
	case "windows":
		err = exec.Command("rundll32", "url.dll,FileProtocolHandler", TWITCH_URL).Start()
	default:
		err = fmt.Errorf("unsupported platform")
	}
	if err != nil {
		log.Fatal(err)
	}

}

var Active = false

func SubscribeToChat(token string) {
	broadcasterid, login, err := GetBroadcasterId(token)
	if err != nil {
		log.Fatal("Error retrieving broadcaster id", err)
	}

	conn, err := net.Dial("tcp", IRC_TWITCH_SERVER)
	if err != nil {
		log.Fatal("Error conectandose a IRC", err)
	}
	//defer conn.Close()

	fmt.Fprintf(conn, "PASS %s\r\n", "oauth:"+token)
	fmt.Fprintf(conn, "NICK %s\r\n", login)
	fmt.Fprintf(conn, "JOIN #%s\r\n", login)

	reader := bufio.NewReader(conn)
	fmt.Println(broadcasterid)
	go func() {
		for {
			line, err := reader.ReadString('\n')
			if err != nil {
				log.Fatal(err)
			}
			//fmt.Println(line)
			splitted := strings.Split(line, "#")
			if len(splitted) == 2 {
				fmt.Println(splitted[1])
				message := strings.Split(splitted[1], ":")
				if len(message) == 2 {
					if Active {
						reproductor.Reproduce(strip(message[1]), message[0])
					} else if strings.Compare(strings.TrimSpace(message[1]), "End of /NAMES list") == 0 && !Active {
						Active = true
					}

				}
			}
		}
	}()

}

func GetBroadcasterId(token string) (string, string, error) {

	client := &http.Client{}

	req, err := http.NewRequest("GET", TWITCH_BROADCASTER_ID, nil)
	if err != nil {
		log.Fatal("Error creating request", err)
		return "", "", err
	}

	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Client-Id", CLIENT_ID)

	resp, err := client.Do(req)
	if err != nil {
		log.Fatal("Error sending request", err)
		return "", "", err
	}
	defer resp.Body.Close()

	fmt.Println("Response status:", resp.Status)

	var result common.Response

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		log.Fatal("Error decoding JSON response", err)
		return "", "", err
	}

	if len(result.Data) > 0 {
		b := result.Data[0]
		broadcasterid := b.Id
		login := b.Login
		return broadcasterid, login, nil
	}

	return "", "", errors.New("No data in")
}

// Source - https://stackoverflow.com/a
// Posted by user5728991, modified by community. See post 'Timeline' for change history
// Retrieved 2025-12-06, License - CC BY-SA 4.0

func strip(s string) string {
	var result strings.Builder
	for i := 0; i < len(s); i++ {
		b := s[i]
		if ('a' <= b && b <= 'z') ||
			('A' <= b && b <= 'Z') ||
			('0' <= b && b <= '9') ||
			b == ' ' {
			result.WriteByte(b)
		}
	}
	return result.String()
}
