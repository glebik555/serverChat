package main

import (
	"bufio"
	"crypto/tls"
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"strings"
	"time"
)

const (
	connHost = "localhost"
	connPort = "8080"
	connType = "tcp"
)

func main() {
	messages := make(map[string]chan string)

	fmt.Println("Starting " + connType + " server on " + connHost + ":" + connPort)
	certFile := flag.String("cert", "cert.pem", "certificate PEM file")
	keyFile := flag.String("key", "key.pem", "key PEM file")
	flag.Parse()
	cert, err := tls.LoadX509KeyPair(*certFile, *keyFile)
	if err != nil {
		log.Fatal(err)
	}
	config := &tls.Config{Certificates: []tls.Certificate{cert}}

	l, err := tls.Listen(connType, ":"+connPort, config)
	if err != nil {
		fmt.Println("Error listening:", err.Error())
		os.Exit(1)
	}

	defer l.Close()

	for {
		c, err := l.Accept()

		if err != nil {
			fmt.Println("Error connecting:", err.Error())
			return
		}
		fmt.Println("Client connected")
		fmt.Println("Client " + c.RemoteAddr().String() + " connected")
		fmt.Println(c.RemoteAddr().String() + " -#1 ")
		go handleConnection(c, c.RemoteAddr().String(), &messages)
	}
}

func handleConnection(conn net.Conn, id string, messages *map[string]chan string) {
	defer conn.Close()

	buffer, err := bufio.NewReader(conn).ReadBytes('\n') // Читает входящее сообщение
	if err != nil {
		fmt.Println("Client " + conn.RemoteAddr().String() + " left")
		panic(err)
	}
	name := string(buffer[:len(buffer)-1])

	fmt.Println("With name " + name)

	(*messages)[name] = make(chan string, 1)
	fmt.Println(conn.RemoteAddr().String() + " -#2 " + name)
	go serveOutcoming(name, conn, messages)
	for {
		buffer, err := bufio.NewReader(conn).ReadBytes('\n')
		if err != nil {
			fmt.Println("Client " + conn.RemoteAddr().String() + " left")
			break
		}

		message := string(buffer[:len(buffer)-1])

		if strings.HasPrefix(message, "to") {

			splitted := strings.Split(message, " ")
			if splitted[2] == "all" {
				for k, _ := range *messages {
					(*messages)[k] <- message
				}
			} else {
				(*messages)[splitted[2]] <- message
				fmt.Println("Wrote message to " + splitted[2] + " " + message)
			}
		}

		if strings.Contains(message, "/exit") {
			fmt.Println("Client " + conn.RemoteAddr().String() + " left")
			break
		} else {
			fmt.Println("Client "+conn.RemoteAddr().String()+" message:", message)
		}

	}
}

func serveOutcoming(name string, conn net.Conn, messages *map[string]chan string) {

	ourChan := (*messages)[name]

	for {
		if name == "all" {
			for key, _ := range *messages {
				ourChan = (*messages)[key]
				select {
				case message := <-ourChan:
					fmt.Println("Read message to " + name + " " + message)
					_, err := conn.Write([]byte(message + "\n"))
					if err != nil {
						panic(err)
					}
				case <-time.After(time.Millisecond):
				}
			}
		} else {
			select {
			case message := <-ourChan:
				fmt.Println("Read message to " + name + " " + message)
				fmt.Println(conn.RemoteAddr().String() + " - #3 " + name)
				_, err := conn.Write([]byte(message + "\n"))
				if err != nil {
					panic(err)
				}
			case <-time.After(time.Millisecond):
			}
		}
	}
}
