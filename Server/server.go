package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"os"
)

type Config struct {
	Host string
	Port string
	Rooms []string
}

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	file, _ := os.Open("./config.json")
	decoder := json.NewDecoder(file)
	config := new(Config)
	err := decoder.Decode(&config)

	if err != nil {
		log.Println(err)
	}

	cs := NewChatServer(config)

	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%s", config.Host, config.Port))

	if err != nil {
		log.Println("Error: ", err)
		os.Exit(1)
	}

	defer listener.Close()
	log.Println("Server started at address: ", fmt.Sprintf("%s:%s", config.Host, config.Port))

	for {
		conn, err := listener.Accept()

		if err != nil {
			log.Println("Error: ", err)
			continue
		}

		cs.CliJoin(NewClient(conn))
	}
}