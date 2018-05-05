package main

import (
	"fmt"
	"os"

	"github.com/rafarlopes/jellybot/chat"
	"github.com/rafarlopes/jellybot/cmd"
	log "github.com/sirupsen/logrus"
)

var (
	token, chatID string
)

func init() {
	log.SetFormatter(&log.TextFormatter{})
	log.SetOutput(os.Stdout)
	log.SetLevel(log.InfoLevel)

	token = os.Getenv("CHAT_API_TOKEN")
	if token == "" {
		panic("set env variable CHAT_API_TOKEN")
	}

	chatID = os.Getenv("DEFAULT_CHAT_ID")
	if chatID == "" {
		panic("set env variable DEFAULT_CHAT_ID")
	}

}

func main() {
	log.Infoln("starting jellybot...")
	c := chat.New(token, chatID)
	msgCh := c.StartReceiving()
	docker, err := cmd.NewDockerClient()

	if err != nil {
		panic(err)
	}

	for msg := range msgCh {
		log.WithField("msg", msg).Info("message received from channel")

		go func(msg chat.Message) {
			out, err := docker.RunCommand(msg.Text)
			if err != nil {
				c.Send(fmt.Sprintf("Sorry, I couldn't run you command. Here's the details: %s", err.Error()), msg.Chat)
				return
			}

			c.Send(out, msg.Chat)
		}(msg)

	}

	log.Info("stopping jellybot...")
}
