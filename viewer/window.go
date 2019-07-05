package main

import (
	"encoding/json"
	"fmt"
	"github.com/davecgh/go-spew/spew"
	"github.com/zserge/webview"
	"io/ioutil"
	"log"
	"os/user"
)

const (
	// windowWidth  = 960
	// windowHeight = 640
	windowWidth  = 550
	windowHeight = 550
)

type configOptions struct {
	// BOT_NAME   string `json:"bot_name"`
	BOT_ID     string `json:"bot_id"`
	USER_ID    string `json:"user_id"`
	BOT_TOKEN  string `json:"bot_token"`
	USER_TOKEN string `json:"user_token"`
}

func handleRPC(w webview.WebView, data string) {
	switch {
	case data == "yes":
		spew.Dump("yo")
	default:

		var firstChar = data[0:1]
		var lastChar = data[len(data)-1:]

		if firstChar == "{" && lastChar == "}" {
			var config configOptions
			json.Unmarshal([]byte(data), &config)

			usr, err := user.Current()
			if err != nil {
				log.Fatal(err)
			}

			// logger.Println(usr.HomeDir)

			connectionString := usr.HomeDir + "/Documents/MessageBridgeData/credentials.json"

			configJSON, _ := json.Marshal(config)
			err = ioutil.WriteFile(connectionString, configJSON, 0644)
			fmt.Printf("%+v", config)

			spew.Dump(connectionString)
			spew.Dump(config)
			w.Exit()
		} else {
			spew.Dump("Not JSON")
		}

		return
	}
}

func main() {
	// url := "https://slack.com/oauth/authorize?client_id=415460872373.416401049063&scope=bot,admin"
	// url := "https://slack.com/oauth/authorize?client_id=415460872373.416401049063&scope=bot,admi?redirect_uri=https://u8g4rrcigf.execute-api.us-east-1.amazonaws.com/default/slack-desktop-oauthna"
	url := "https://slack.com/oauth/authorize?client_id=415460872373.416401049063&scope=bot,admin"
	// url := "http://localhost:1313/desktop-setup"
	w := webview.New(webview.Settings{
		Width:                  windowWidth,
		Height:                 windowHeight,
		Title:                  "Message Bridge Sign In",
		Resizable:              true,
		URL:                    url,
		ExternalInvokeCallback: handleRPC,
	})
	w.SetColor(255, 255, 255, 255)
	defer w.Exit()
	w.Run()
}
