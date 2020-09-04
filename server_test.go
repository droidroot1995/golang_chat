package main

import (
	"bufio"
	"net"
	"testing"
)

type Config struct {
	Host string
	Port string
	Username string
	Rooms []string
}

var config *Config
var conn net.Conn
var err error
var reader *bufio.Reader
var writer *bufio.Writer

func InitConfig() {
	config = &Config {
		Host: "127.0.0.1",
		Port: "3000",
		Username: "droidroot",
		Rooms: make([]string, 0),
	}

	config.Rooms = append(config.Rooms, "one")
	config.Rooms = append(config.Rooms, "two")
}

func TestConnection(t *testing.T) {

	InitConfig()
	conn, err = net.Dial("tcp", config.Host+":"+config.Port)
	if err != nil {
		t.Error(err)
	}

	reader = bufio.NewReader(conn)
	writer = bufio.NewWriter(conn)

	st, err := reader.ReadString('\n')

	if err != nil {
		t.Error(err)
	}

	t.Log(st)

	st, err = reader.ReadString('\n')

	if err != nil {
		t.Error(err)
	}

	t.Log(st)
}

func TestNameChange(t *testing.T) {

	_, err := writer.WriteString("/set_name " + config.Username + "\n")

	if err != nil {
		t.Error(err)
	}

	err = writer.Flush()
	if err != nil {
		t.Error(err)
	}


	st, err := reader.ReadString('\n')

	if err != nil {
		t.Error(err)
	}

	t.Log(st)
}

func TestSub(t *testing.T) {


	_, err = writer.WriteString("/sub " + config.Rooms[0] + "\n")

	if err != nil {
		t.Error(err)
	}

	err = writer.Flush()
	if err != nil {
		t.Error(err)
	}

	st, err := reader.ReadString('\n')

	if err != nil {
		t.Error(err)
	}

	t.Log(st)

}

func TestPub(t *testing.T) {

	_, err = writer.WriteString("/pub " + config.Rooms[0] + " hello\n")

	if err != nil {
		t.Error(err)
	}

	err = writer.Flush()
	if err != nil {
		t.Error(err)
	}

	st, err := reader.ReadString('\n')

	if err != nil {
		t.Error(err)
	}

	t.Log(st)
}

func TestUnsub(t *testing.T) {

	_, err = writer.WriteString("/unsub " + config.Rooms[0] + "\n")

	if err != nil {
		t.Error(err)
	}

	err = writer.Flush()
	if err != nil {
		t.Error(err)
	}

	st, err := reader.ReadString('\n')

	if err != nil {
		t.Error(err)
	}

	t.Log(st)
}

func TestHist(t *testing.T) {

	_, err = writer.WriteString("/hist " + config.Rooms[0] + "\n")

	if err != nil {
		t.Error(err)
	}

	err = writer.Flush()
	if err != nil {
		t.Error(err)
	}

	hist, err := reader.ReadString('\n')

	if err != nil {
		t.Error(err)
	}

	t.Log(hist)

	defer conn.Close()
}


