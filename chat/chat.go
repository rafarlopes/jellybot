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

// MessageSenderReceiver is a combination of both interfaces MessageReceiver and MessageSender
type MessageSenderReceiver interface {
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
