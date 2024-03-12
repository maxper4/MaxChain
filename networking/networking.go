package networking

import (
	"maxchain/config"
	"maxchain/logging"
	"net"
	"strconv"
)

var peersConnections []net.Conn

func Init(config config.Configuration) {
	logging.Log("Initializing networking package", "networking", "INFO")
	logging.Log("Server Running on port "+strconv.Itoa(config.ListeningPort), "networking", "INFO")
	server, err := net.Listen("tcp", ":"+strconv.Itoa(config.ListeningPort))
	if err != nil {
		logging.PanicWithLog("Cannot listen on port: "+err.Error(), "networking")
	}
	go networkingLoop(server)

	go InitPeersConnections(config)
}

func InitPeersConnections(config config.Configuration) {
	peersConnections = make([]net.Conn, len(config.Peers))
	for i := 0; i < len(config.Peers); i++ {
		connection, err := net.Dial("tcp", config.Peers[i].Ip+":"+strconv.Itoa(config.Peers[i].Port))
		if err != nil {
			logging.Log("Cannot connect to peer: "+err.Error(), "networking", "ERROR")
		} else {
			peersConnections[i] = connection // TODO: defer connection.Close()
		}
	}
}

func networkingLoop(server net.Listener) {
	defer server.Close()
	for {
		connection, err := server.Accept()
		if err != nil {
			logging.PanicWithLog("Error accepting: "+err.Error(), "networking")
		}
		logging.Log("New client request from "+connection.RemoteAddr().String()+" accepted", "networking", "INFO")
		go processConnection(connection)
	}
}

func processConnection(connection net.Conn) {
	logging.Log("Processing connection", "networking", "INFO")
	buffer := make([]byte, 1024)
	defer connection.Close()
	for {
		mLen, err := connection.Read(buffer)
		if err != nil {
			logging.Log("Error reading: "+err.Error(), "networking", "ERROR")
		}
		logging.Log("Received: "+string(buffer[:mLen]), "networking", "INFO")
		handshake(connection)
		_, err = connection.Write([]byte("Thanks! Got your message:" + string(buffer[:mLen])))
	}
}

func handshake(connection net.Conn) {

}

func Broadcast(message string) {
	for i := 0; i < len(peersConnections); i++ {
		if peersConnections[i] != nil {
			peersConnections[i].Write([]byte(message))
		}
	}
}
