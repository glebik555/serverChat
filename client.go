package main

import (
	"bufio"
	"flag"
	"fmt"
	"net"
	"os"
	"strings"

)

const (
	connHost_ = "localhost"
	connPort_ = "8080"
	connType_ = "tcp"
)

var name = flag.String("name", "", "Advertised name of client")

func main() {
	flag.Parse()

	fmt.Println("Connecting to " + connType_ + " server " + connHost_ + ":" + connPort_)
	conn, err := net.Dial(connType_, connHost_+":"+connPort_)
	if err != nil {
		fmt.Println("Error connecting:", err.Error())
		os.Exit(1)
	}
	defer conn.Close()
	reader := bufio.NewReader(os.Stdin)

	conn.Write([]byte(*name + "\n"))

	go serverIncoming(conn)

	for {
		fmt.Print("Enter text: ")  // to0 (glebasta)1 sasha2 privet3 ... Сообщение серверу
		var sb strings.Builder
		sb.WriteString("to " + "(" + *name + ") ")
		text, _ := reader.ReadString('\n')
		sb.WriteString(text)
		fmt.Println("b: " + sb.String())
		conn.Write([]byte(sb.String() + "\n"))
		if strings.Contains(text, "/exit") {
			break
		}
	}
}

func serverIncoming(conn net.Conn) { // Сообщение от сервера
	for {
		buffer, err := bufio.NewReader(conn).ReadBytes('\n')
		splitted := strings.Split(string(buffer), " ")
		if err != nil {
			break
		}
		var sb strings.Builder
		for  i:=3; i < len(splitted);i++{
			sb.WriteString(splitted[i] + " ")
		}

		fmt.Println("(Server) From " + splitted[1] + ": " + sb.String())
	}
}
