package main

import (
	"fmt"
	"os"

	"github.com/rafarlopes/jellybot/chat"
	"github.com/rafarlopes/jellybot/cmd"
	log "github.com/sirupsen/logrus"
)

var (
	token, userID, chatID string
)

func init() {
	log.SetFormatter(&log.TextFormatter{})
	log.SetOutput(os.Stdout)
	log.SetLevel(log.InfoLevel)

	token = os.Getenv("CHAT_API_TOKEN")
	userID = os.Getenv("BOT_USER_ID")
	chatID = os.Getenv("DEFAULT_CHAT_ID")

	if token == "" {
		panic("set env variable CHAT_API_TOKEN")
	}

	if userID == "" {
		panic("set env variable BOT_USER_ID")
	}

	if chatID == "" {
		panic("set env variable DEFAULT_CHAT_ID")
	}

}

func main() {
	log.Infoln("starting jellybot...")

	chat := chat.New(token, userID, chatID)
	msgCh := chat.StartReceiving()

	for msg := range msgCh {
		log.WithField("msg", msg).Info("message received from channel")
		go func() {
			out, err := cmd.Parse(msg.Text)
			if err != nil {
				chat.Send(fmt.Sprintf("Sorry, I couldn't run you command. Here's the details: %s", err.Error()), msg.Chat)
				return
			}

			chat.Send(out, msg.Chat)
		}()

	}

	log.Info("stopping jellybot...")
}
