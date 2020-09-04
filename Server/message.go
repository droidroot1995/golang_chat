package main

import "fmt"

type Message struct {
	msg string
	cli *Client
}

func NewMessage(msg string, client *Client) *Message {
	message := &Message{
		msg: msg,
		cli: client,
	}

	return message
}

func (msg *Message) AsString() string {
	strMsg := fmt.Sprintf("%s:%s", msg.cli.name, msg.msg)
	return strMsg
}
