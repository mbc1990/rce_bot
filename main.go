package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/nlopes/slack"
	"log"
	"os"
	"os/exec"
	"strings"
)

type Rcebot struct {
	SlackAPI *slack.Client
}

type Configuration struct {
	Token string
	BotID string
}

type Message struct {
	ChannelID string
	Content   string
}

func (r *Rcebot) HandleMessage(ev *slack.MessageEvent) {
	spl := strings.Split(ev.Text, " ")
	if len(spl) < 2 {
		return
	}
	if spl[0] == "$" {
		var msg Message
		cmd := exec.Command("sh", "-c", spl[1])
		var outGood bytes.Buffer
		var outBad bytes.Buffer
		cmd.Stdout = &outGood
		cmd.Stderr = &outBad
		cmd.Run()
		if len(outGood.String()) > 0 {
			msg = Message{ChannelID: ev.Channel, Content: outGood.String()}
		}
		if len(outBad.String()) > 0 {
			msg = Message{ChannelID: ev.Channel, Content: outBad.String()}
		}
		fmt.Printf(outGood.String())
		fmt.Printf(outBad.String())
		params := slack.PostMessageParameters{Username: "rce_bot", IconEmoji: ":simple_smile:"}
		fmt.Println("Attempting to send message: " + msg.Content)
		_, _, err := r.SlackAPI.PostMessage(msg.ChannelID, msg.Content, params)
		if err != nil {
			fmt.Printf("failed to post message: %v\n", err)
		}
	}
}

func (r *Rcebot) Start() {
	rtm := r.SlackAPI.NewRTM()
	go rtm.ManageConnection()
	for msg := range rtm.IncomingEvents {
		switch ev := msg.Data.(type) {
		case *slack.MessageEvent:
			go r.HandleMessage(ev)
		case *slack.InvalidAuthEvent:
			log.Fatal("Invalid credentials")
		}
	}
}

func main() {
	fmt.Println("Starting bad idea bot")
	confPath := flag.String("conf", "conf.json", "Path to json configuration file")
	flag.Parse()

	file, err := os.Open(*confPath)
	if err != nil {
		log.Fatalf("failed to open config: %v", err)
	}

	var conf Configuration
	err = json.NewDecoder(file).Decode(&conf)
	if err != nil {
		log.Fatalf("failed to unmarshal config: %v", err)
	}
	bot := Rcebot{
		SlackAPI: slack.New(conf.Token),
	}
	bot.Start()

}
