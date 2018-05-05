package chat

import (
	"fmt"
	"strings"
	"time"

	"github.com/nlopes/slack"
	log "github.com/sirupsen/logrus"
)

type slackChat struct {
	cli       *slack.Client
	rtm       *slack.RTM
	channelID string
}

// New creates an instance of SlackChat using a token
func New(token, channelID string) MessageSenderReceiver {
	log.WithFields(log.Fields{
		"token":     token,
		"channelID": channelID,
	}).Info("creating new slack chat")

	client := slack.New(token)
	rtm := client.NewRTM()

	//TODO should we move it to its own method?
	go rtm.ManageConnection()

	return &slackChat{
		cli:       client,
		rtm:       rtm,
		channelID: channelID,
	}
}

// Send a message to user or channel in slack
func (s *slackChat) Send(message, to string) {
	s.rtm.SendMessage(s.rtm.NewOutgoingMessage(message, to))
}

// Receive opens a channel to handler messages from slack to the specific bot user
func (s *slackChat) StartReceiving() <-chan Message {
	ch := make(chan Message)

	go func() {
		defer close(ch)
	Loop:
		for msg := range s.rtm.IncomingEvents {
			switch ev := msg.Data.(type) {
			case *slack.ConnectedEvent:
				time.AfterFunc(1*time.Second, func() {
					log.Info("bot connected, sending hello world")
					s.Send("Hi there, I'm here!", s.channelID)
				})

			case *slack.MessageEvent:
				log.WithFields(log.Fields{
					"channel": ev.Channel,
					"message": ev.Text,
					"from":    ev.User,
				}).Info("message received")

				info := s.rtm.GetInfo()
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
