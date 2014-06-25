package main

import (
	retiwe "../"
	"flag"
	"fmt"
	"io"
	"net"
	"strings"
)

//EXAMPLE:
//go run 01.go [port]
//go run 01.go 23
//telnet localhost 23

var hub = &retiwe.Connections{
	Clients:      make(map[chan string]bool),
	AddClient:    make(chan (chan string)),
	RemoveClient: make(chan (chan string)),
	Messages:     make(chan string),
}

func handleConnection(c net.Conn) {
	messageChannel := make(chan string)
	defer func() {
		hub.RemoveClient <- messageChannel
		c.Close()
	}()
	readbufferChannel := make(chan string)
	go func() {
		for {
			buf := make([]byte, 4096)
			n, err := c.Read(buf)
			if err == nil {
				readbufferChannel <- strings.TrimSpace(string(buf[0:n]))
			}
		}
	}()

	hub.AddClient <- messageChannel
	fmt.Printf("Connection from %v established.\n", c.RemoteAddr())
	io.WriteString(c, fmt.Sprintf("Welcome on %s\n", c.LocalAddr()))
	for {
		select {
		case msg := <-messageChannel:
			//receive message
			io.WriteString(c, msg+"\n")
		case msg := <-readbufferChannel:
			//send message
			if msg == ":q" {
				//hub.RemoveClient <- messageChannel
				c.Close()
				return
			}
			hub.Messages <- fmt.Sprintf("%s: %s", c.RemoteAddr(), msg)
		}
	}
	fmt.Printf("Connection from %v closed.\n", c.RemoteAddr())
	hub.RemoveClient <- messageChannel
	c.Close()
}

func main() {
	flag.Parse()
	port := ":" + flag.Arg(0)
	if port == ":" {
		port = ":1234"
	}
	fmt.Printf("Listening on port %v\n", port)
	hub.Init()
	ln, _ := net.Listen("tcp", port)
	for {
		conn, _ := ln.Accept()
		go handleConnection(conn)
	}
}
