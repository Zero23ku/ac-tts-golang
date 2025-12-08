package github

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"ac-tts/internal/common"
	"ac-tts/internal/logging"
)

var releaseURL = "https://api.github.com/repos/Zero23ku/ac-tts-golang/tags"

func GetLatestReleaseVersion() {
	res, err := http.Get(releaseURL)

	if err != nil {
		logging.CreateLog(err)
		log.Fatal(err)
	}

	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)

	if err != nil {
		logging.CreateLog(err)
		log.Fatal(err)
	}

	var response []common.Release

	if err := json.Unmarshal(body, &response); err != nil {
		logging.CreateLog(err)
		log.Fatal(err)
	}

	if len(response) > 0 {
		fmt.Println(response[0].Name)
	}
}
