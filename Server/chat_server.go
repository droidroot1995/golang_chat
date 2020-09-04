package main

import (
	"fmt"
	"strings"
)

type ChatServer struct {
	users []*Client
	chatRooms map[string]*ChatRoom

	message chan *Message
	join chan *Client
	leave chan *Client

}

func NewChatServer(config *Config) *ChatServer {
	cs := &ChatServer {
		users: make([]*Client, 0),
		chatRooms: make(map[string]*ChatRoom),
		message: make(chan *Message),
		join: make(chan *Client),
		leave: make(chan *Client),
	}

	for _, room := range config.Rooms {
		cr := NewChatRoom(room)
		cs.chatRooms[room] = cr
	}

	cs.Listen()

	return cs
}

func (cs *ChatServer) Listen() {
	go func() {
		for {
			select {
			case msg := <-cs.message:
				cs.ParseMsg(msg)
			case client := <-cs.join:
				cs.CliJoin(client)
			case client := <- cs.leave:
				cs.CliLeave(client)
			}
		}
	}()
}

func (cs *ChatServer) CliJoin(client *Client) {
	cs.users = append(cs.users, client)

	client.outChan <- MSG_WELCOME
	go func() {
		for msg := range client.inChan {
			cs.message <- msg
		}
		cs.leave <- client
	}()
}

func (cs *ChatServer) CliLeave(client *Client) {
	if len(client.ChatRooms) != 0 {
		for _, room := range client.ChatRooms {
			room.RemoveUser(client)
		}
	}

	for idx, cli := range cs.users {
		if client == cli {
			cs.users = append(cs.users[:idx], cs.users[idx+1:]...)
			break
		}
	}

	client.conn.Close()
}

func (cs *ChatServer) ParseMsg(msg *Message) {
	switch {
	case strings.HasPrefix(msg.msg, CMD_SET_NAME):
		cs.ChangeClientName(msg)
	case strings.HasPrefix(msg.msg, CMD_SUB):
		cs.SubClientToRoom(msg)
	case strings.HasPrefix(msg.msg, CMD_UNSUB):
		cs.UnsubClientFromRoom(msg)
	case strings.HasPrefix(msg.msg, CMD_PUB):
		cs.PubMessageToRoom(msg)
	case strings.HasPrefix(msg.msg, CMD_HIST):
		cs.SendRoomHist(msg)
	case strings.HasPrefix(msg.msg, CMD_LIST):
		cs.GetRoomsList(msg.cli)
	case strings.HasPrefix(msg.msg, CMD_HELP):
		cs.ShowHelp(msg.cli)
	case strings.HasPrefix(msg.msg, CMD_QUIT):
		cs.CliLeave(msg.cli)
	}
}

func (cs *ChatServer) ChangeClientName(msg *Message) {

	name := strings.TrimSpace(strings.TrimSuffix(strings.SplitN(msg.msg, " ", 2)[1], "\n"))

	for _, cli := range cs.users {
		if cli == msg.cli {
			cli.name = name
			cli.outChan <- fmt.Sprintf("%s%s\n", SUCCESS_NAME_CHANGE, name)
			break
		}
	}
}

func (cs *ChatServer) SubClientToRoom(msg *Message) {

	room := strings.TrimSpace(strings.TrimSuffix(strings.SplitN(msg.msg, " ", 2)[1], "\n"))

	client := msg.cli
	for _, cli := range cs.users {
		if cli == msg.cli {
			client = cli
			break
		}
	}

	if _, ok := cs.chatRooms[room]; !ok {
		client.outChan <- FAIL_SUB_RNE
		return
	}

	contains := cs.chatRooms[room].ContainsUser(client.name)

	if !contains {
		client.outChan <- fmt.Sprintf("%s%s\n", SUCCESS_SUB, room)
		cs.chatRooms[room].AddUser(client)
	} else {
		client.outChan <- FAIL_SUB_AE
	}

}

func (cs *ChatServer) UnsubClientFromRoom(msg *Message) {
	room := strings.TrimSpace(strings.TrimSuffix(strings.SplitN(msg.msg, " ", 2)[1], "\n"))

	if _, ok := cs.chatRooms[room]; !ok {
		msg.cli.outChan <- FAIL_UNSUB_RNE
		return
	}

	client := msg.cli
	for _, cli := range cs.users {
		if cli == msg.cli {
			client = cli
			break
		}
	}

	contains := cs.chatRooms[room].ContainsUser(client.name)

	if !contains {
		client.outChan <- FAIL_UNSUB_NE
	} else {
		cs.chatRooms[room].RemoveUser(msg.cli)

		for idx, r := range client.ChatRooms {
			if r == cs.chatRooms[room] {
				client.ChatRooms = append(client.ChatRooms[:idx], client.ChatRooms[idx+1:]...)
			}
		}

		client.outChan <- fmt.Sprintf("%s%s\n", SUCCESS_UNSUB, room)
	}
}

func (cs *ChatServer) PubMessageToRoom(msg *Message) {
	msgInfo := strings.SplitN(msg.msg, " ", 3)
	room := strings.TrimSpace(strings.TrimSuffix(msgInfo[1], "\n"))
	msgTxt := msgInfo[2]

	if _, ok := cs.chatRooms[room]; !ok {
		msg.cli.outChan <- FAIL_PUB_RNE
		return
	}

	client := msg.cli
	for _, cli := range cs.users {
		if cli == msg.cli {
			client = cli
			break
		}
	}

	contains := cs.chatRooms[room].ContainsUser(client.name)

	if contains {
		cs.chatRooms[room].AddMessage(fmt.Sprintf("%s:%s\n", client.name, msgTxt))

	} else {
		client.outChan <- FAIL_PUB_NE
	}
}

func (cs *ChatServer) SendRoomHist(msg *Message) {

	room := strings.TrimSpace(strings.TrimSuffix(strings.SplitN(msg.msg, " ", 2)[1], "\n"))

	if _, ok := cs.chatRooms[room]; !ok {
		msg.cli.outChan <- FAIL_HIST_GET_NE
		return
	}

	client := msg.cli
	for _, cli := range cs.users {
		if cli == msg.cli {
			client = cli
			break
		}
	}

	contains := cs.chatRooms[room].ContainsUser(client.name)

	if contains {
		hist := cs.chatRooms[room].History()

		st := ""

		for _, ms := range hist {
			st += ms
			st += "\n"
		}

		client.outChan <- st

	} else {
		client.outChan <- FAIL_HIST_GET_NS
	}
}

func (cs *ChatServer) GetRoomsList(cli *Client) {
	lst := ""

	for _, room := range cs.chatRooms {
		lst += room.name
		lst += "\n"
	}

	cli.outChan <- lst
}

func (cs *ChatServer) ShowHelp(cli *Client) {
	cli.outChan <- fmt.Sprintf("%s\n", MSG_HELP)
}