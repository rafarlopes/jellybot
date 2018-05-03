package chat

import (
	"fmt"
	"strings"
	"time"

	"github.com/nlopes/slack"
	log "github.com/sirupsen/logrus"
)

// New creates an instance of SlackChat using a token
func New(token, userID, channelID string) *Chat {
	log.WithFields(log.Fields{
		"token":     token,
		"userID":    userID,
		"channelID": channelID,
	}).Info("creating new slack chat")

	client := slack.New(token)
	rtm := client.NewRTM()

	go rtm.ManageConnection()
	return &Chat{
		client: rtm,
		userID: userID,
		chatID: channelID,
	}
}

// Send a message to user or channel in slack
func (c *Chat) Send(message, to string) {
	rtm := c.client.(*slack.RTM)
	rtm.SendMessage(rtm.NewOutgoingMessage(message, to))
}

// Receive opens a channel to handler messages from slack to the specific bot user
func (c *Chat) StartReceiving() <-chan Message {
	rtm := c.client.(*slack.RTM)
	ch := make(chan Message)

	go func() {
		defer close(ch)
	Loop:
		for msg := range rtm.IncomingEvents {
			switch ev := msg.Data.(type) {
			case *slack.ConnectedEvent:
				time.AfterFunc(1*time.Second, func() {
					log.Info("bot connected, sending hello world")
					c.Send("Hi there, I'm here!", c.chatID)
				})

			case *slack.MessageEvent:
				log.WithFields(log.Fields{
					"channel": ev.Channel,
					"message": ev.Text,
					"from":    ev.User,
				}).Info("message received")

				info := rtm.GetInfo()
				prefix := fmt.Sprintf("<@%s> ", info.User.ID)
				if ev.User != info.User.ID && strings.HasPrefix(ev.Text, prefix) {
					ch <- Message{
						Text: ev.Text,
						From: ev.User,
						To:   info.User.ID,
						Chat: ev.Channel,
					}
				}

			case *slack.LatencyReport:
				//fmt.Printf("Current latency: %v\n", ev.Value)
				//TODO handle later as heartbeat

			case *slack.RTMError:
				log.Errorf("Error: %s\n", ev.Error())
				break Loop

			case *slack.InvalidAuthEvent:
				log.Error(ErrInvalidCredentials)
				break Loop
			}
		}
	}()

	return ch
}
