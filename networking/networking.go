package networking

import (
	"fmt"
	"net"
	"strconv"
	"maxchain/config"
)

var peersConnections []net.Conn

func Init(config config.Configuration) {
	fmt.Println("Initializing networking package")
	fmt.Printf("Server Running on port %d...\n", config.ListeningPort)
	server, err := net.Listen("tcp", ":"+strconv.Itoa(config.ListeningPort))
	if err != nil {
		panic("Cannot listen on port: " + err.Error())
	}
	go networkingLoop(server)

	go InitPeersConnections(config)
}

func InitPeersConnections(config config.Configuration) {
	peersConnections = make([]net.Conn, len(config.Peers))
	for i := 0; i < len(config.Peers); i++ {
		connection, err := net.Dial("tcp", config.Peers[i].Ip + ":" + strconv.Itoa(config.Peers[i].Port))
		if err != nil {
			fmt.Println("Cannot connect to peer: " + err.Error())
		} else {
			peersConnections[i] = connection	// TODO: defer connection.Close()
		}
	}
}

func networkingLoop(server net.Listener) {
	defer server.Close()
	for {
		connection, err := server.Accept()
		if err != nil {
				panic("Error accepting: " + err.Error())
		}
		fmt.Println("New client request from " + connection.RemoteAddr().String() + " accepted")
		go processConnection(connection)
	}
}

func processConnection(connection net.Conn) {
	fmt.Println("Processing connection")
	buffer := make([]byte, 1024)
	defer connection.Close()
	for {
        mLen, err := connection.Read(buffer)
        if err != nil {
                fmt.Println("Error reading:", err.Error())
        }
        fmt.Println("Received: ", string(buffer[:mLen]))
        _, err = connection.Write([]byte("Thanks! Got your message:" + string(buffer[:mLen])))
	}
}

func Broadcast(message string) {
	for i := 0; i < len(peersConnections); i++ {
		if peersConnections[i] != nil {
			peersConnections[i].Write([]byte(message))
		}
	}
}



