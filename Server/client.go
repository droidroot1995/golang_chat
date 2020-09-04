package main

import (
	"bufio"
	"log"
	"net"
	"strings"
	"sync"
)

var wg sync.WaitGroup

type Client struct {
	name string
	ChatRooms []*ChatRoom

	conn net.Conn

	inChan chan *Message
	outChan chan string

	reader *bufio.Reader
	writer *bufio.Writer
}

func NewClient(conn net.Conn) *Client {
	reader := bufio.NewReader(conn)
	writer := bufio.NewWriter(conn)

	cli := &Client{
		name: "",
		ChatRooms: make([]*ChatRoom, 0),
		conn: conn,

		inChan: make(chan *Message),
		outChan: make(chan string),

		reader: reader,
		writer: writer,
	}

	cli.Listen()

	return cli
}

func (cli *Client) Listen() {
	go cli.ReadMessages()
	go cli.WriteMessages()
}

func (cli *Client) ReadMessages() {
	for {
		msg, err := cli.reader.ReadString('\n')

		if err != nil {
			log.Println(err)
			break
		}

		mess := NewMessage(strings.TrimSuffix(msg, "\n"), cli)
		cli.inChan <- mess
	}

	close(cli.inChan)
	log.Println("Closed client's read channel")
}

func (cli *Client) WriteMessages() {
	for msg := range cli.outChan{
		_, err := cli.writer.WriteString(msg)
		if err != nil {
			log.Println(err)
			break
		}

		err = cli.writer.Flush()

		if err != nil {
			log.Println(err)
			break
		}
	}
	log.Println("Closed client's write channel")
}

