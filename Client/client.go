package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"os"
	"strings"
	"sync"
)

const (
	SUCCESS_SUB = "Successfully subscribed to room"
	SUCCESS_UNSUB = "Successfully unsubscribed from room"
	SUCCESS_NAME_CHANGE = "Name successfully changed"
)

var wg sync.WaitGroup

type Config struct {
	Host string
	Port string
	Username string
	Rooms []string
}

type Client struct {
	name string
	Rooms []string
	conn net.Conn
	config *Config
	first bool
	decoder *json.Decoder
}

func main() {
	file, _ := os.Open("./config.json")
	decoder := json.NewDecoder(file)
	config := new(Config)
	err := decoder.Decode(&config)

	if err != nil {
		log.Print(err)
	}

	client := &Client {
		name: config.Username,
		Rooms: make([]string, 0),
		config: config,
		first: true,
		decoder: decoder,
	}

	for _, room := range config.Rooms {
		client.Rooms = append(client.Rooms, room)
	}

	wg.Add(1)

	client.conn, err = net.Dial("tcp", config.Host + ":" + config.Port)

	if err != nil {
		log.Println(err)
	}

	go client.ReadMessages()
	go client.WriteMessages()

	wg.Wait()
}

func (cli *Client) ReadMessages() {
	reader := bufio.NewReader(cli.conn)

	for {
		st, err := reader.ReadString('\n')

		if strings.Contains(st, SUCCESS_NAME_CHANGE) {
			username := strings.TrimSpace(strings.TrimSuffix(strings.SplitN(st, ":", 2)[1], "\n"))
			cli.config.Username = username

			file, err := json.MarshalIndent(cli.config, "", " ")

			if err != nil {
				log.Print(err)
			}

			err = ioutil.WriteFile("config.json", file, 0644)

			if err != nil {
				log.Print(err)
			}
		}

		if strings.Contains(st, SUCCESS_SUB) {
			room := strings.TrimSpace(strings.TrimSuffix(strings.SplitN(st, ":", 2)[1], "\n"))

			confContains := false
			cliContains := false

			for _, r := range cli.config.Rooms {
				if r == room {
					confContains = true
					break
				}
			}

			if !confContains {
				cli.config.Rooms = append(cli.config.Rooms, room)
			}

			for _, r := range cli.Rooms {
				if r == room {
					cliContains = true
					break
				}
			}

			if !cliContains {
				cli.Rooms = append(cli.Rooms, room)
			}

			file, err := json.MarshalIndent(cli.config, "", " ")

			if err != nil {
				log.Print(err)
			}

			err = ioutil.WriteFile("config.json", file, 0644)

			if err != nil {
				log.Print(err)
			}
		}

		if strings.Contains(st, SUCCESS_UNSUB) {
			room := strings.TrimSpace(strings.TrimSuffix(strings.SplitN(st, ":", 2)[1], "\n"))

			for idx, cr := range cli.config.Rooms {
				if cr == room {
					cli.config.Rooms = append(cli.config.Rooms[:idx], cli.config.Rooms[idx+1:]...)
					break
				}
			}

			for idx, cr := range cli.Rooms {
				if cr == room {
					cli.Rooms = append(cli.Rooms[:idx], cli.Rooms[idx+1:]...)
					break
				}
			}

			file, err := json.MarshalIndent(cli.config, "", " ")

			if err != nil {
				log.Print(err)
			}

			err = ioutil.WriteFile("config.json", file, 0644)

			if err != nil {
				log.Print(err)
			}
		}

		if err != nil {
			fmt.Println("Disconnected from server")
			wg.Done()
			return
		}
		fmt.Print(st)
	}
}

func (cli *Client) WriteMessages() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriterSize(cli.conn, 254)

	for {
		if cli.first {
			_, err := writer.WriteString(fmt.Sprintf("/set_name %s\n", cli.config.Username))

			if err != nil {
				log.Println(err)
				os.Exit(1)
			}

			err = writer.Flush()
			if err != nil {
				log.Println(err)
				os.Exit(1)
			}

			rooms := cli.Rooms

			for _, room := range rooms {
				_, err := writer.WriteString(fmt.Sprintf("/unsub %s\n", room))

				if err != nil {
					log.Println(err)
					os.Exit(1)
				}

				err = writer.Flush()
				if err != nil {
					log.Println(err)
					os.Exit(1)
				}
			}

			for _, room := range rooms {
				_, err := writer.WriteString(fmt.Sprintf("/sub %s\n", room))

				if err != nil {
					log.Println(err)
					os.Exit(1)
				}

				err = writer.Flush()
				if err != nil {
					log.Println(err)
					os.Exit(1)
				}
			}


			cli.first = false
		} else {
			st, err := reader.ReadString('\n')

			if err != nil {
				log.Println(err)
				os.Exit(1)
			}

			_, err = writer.WriteString(st)

			if err != nil {
				log.Println(err)
				os.Exit(1)
			}

			err = writer.Flush()
			if err != nil {
				log.Println(err)
				os.Exit(1)
			}
		}
	}
}

