package main

import "fmt"

type ChatRoom struct {
	name string
	clients []*Client
	messages []string
}

func NewChatRoom(rname string) *ChatRoom {
	croom := &ChatRoom{
		name: rname,
		clients: make([]*Client, 0),
		messages: make([]string, 0, 128),
	}

	return croom
}

func (croom *ChatRoom) AddUser(client *Client) {
	croom.clients = append(croom.clients, client)

	hist := ""

	for _, ms := range croom.messages {
		hist += ms
		hist += "\n"
	}

	client.outChan <- hist
}

func (croom *ChatRoom) RemoveUser(client *Client) {
	for idx, cli := range croom.clients {
		if cli == client {
			croom.clients = append(croom.clients[:idx], croom.clients[idx+1:]...)
			break
		}
	}
}

func (croom *ChatRoom) AddMessage(msg string) {
	if len(croom.messages) < 128 {
		croom.messages = append(croom.messages, msg)
	} else {
		croom.messages = append(croom.messages[1:], msg)
	}

	for _, cli := range croom.clients {
		cli.outChan <- fmt.Sprintf("%s: %s", croom.name, msg)
	}
}

func (croom *ChatRoom) LastMessage() string {
	if len(croom.messages) == 0 {
		return ""
	} else {
		return croom.messages[len(croom.messages) - 1]
	}
}

func (croom *ChatRoom) History() []string {
	return croom.messages
}

func (croom *ChatRoom) ContainsUser(username string) bool {
	contains := false
	if len(croom.clients) == 0 {
		contains = false
	} else {
		for _, cli := range croom.clients {
			if cli.name == username {
				contains = true
				break
			}
		}
	}

	return contains
}