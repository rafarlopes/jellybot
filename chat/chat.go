package chat

import "errors"

// ErrInvalidCredentials is raised when a invalid token is provided
var ErrInvalidCredentials = errors.New("could not connect to chat with provided credentials")

// MessageReceiver returns a unbuffered channel of all messages send to the chat
type MessageReceiver interface {
	StartReceiving() <-chan Message
}

// MessageSender sends a message to chat
type MessageSender interface {
	Send(message, to string)
}

// Chat combination of both interfaces to Send and Receive messages
type Chat struct {
	client interface{}
	userID string
	chatID string
	MessageSender
	MessageReceiver
}

// Message struct is used to represent bot incomming messages
type Message struct {
	Text string
	From string
	Chat string
	To   string
}
