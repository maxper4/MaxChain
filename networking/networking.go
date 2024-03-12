package networking

import (
	"encoding/json"
	"errors"
	"maxchain/config"
	"maxchain/cryptography"
	"maxchain/logging"
	"net"
	"strconv"
)

var peersConnections []net.Conn
var configInstance config.Configuration

func Init(config config.Configuration) {
	logging.Log("Initializing networking package", "networking", "INFO")
	logging.Log("Server Running on port "+strconv.Itoa(config.ListeningPort), "networking", "INFO")
	configInstance = config
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
			connection.Close()
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
	if askHandshake(connection) {
		logging.Log("Handshake successful", "networking", "INFO")
		peersConnections = append(peersConnections, connection) // discovered peer
	} else {
		logging.Log("Handshake failed", "networking", "ERROR")
		connection.Close()
		return
	}

	defer connection.Close()
	countErrors := 0
	for {
		msg, err := readMsg(connection)
		if err != nil {
			countErrors++
			if countErrors > configInstance.DropConnectionAfterErrors && configInstance.DropConnectionAfterErrors > 0 {
				logging.Log("Too many errors, closing connection", "networking", "ERROR")
				connection.Close()
				return
			} else {
				logging.Log("Error reading message: "+err.Error(), "networking", "ERROR")
			}
			continue
		}

		switch msg.Tag {
		case 0:
			completeHandshake(connection)
		case 1:
			continue // should not happen?
		default:
			logging.Log("Unknown message tag: "+strconv.Itoa(int(msg.Tag)), "networking", "ERROR")
			countErrors++
			if countErrors > configInstance.DropConnectionAfterErrors && configInstance.DropConnectionAfterErrors > 0 {
				logging.Log("Too many errors, closing connection", "networking", "ERROR")
				connection.Close()
				return
			}
		}
	}
}

func askHandshake(connection net.Conn) bool {
	challenge := cryptography.Rand()
	sendMsg(connection, MsgHandshakeChallenge{
		Challenge: challenge,
	}.ToMsg().ToBytes())

	response, err := readMsg(connection)
	if err != nil {
		logging.Log("Error reading handshake response:"+err.Error(), "networking", "ERROR")
		return false
	}
	logging.Log("Received response: "+response.Content, "networking", "INFO")

	responseMsg, err := response.ToMsgHandshakeResponse()
	if err != nil {
		logging.Log("Error deserializing handshake response:"+err.Error(), "networking", "ERROR")
		return false
	}

	return cryptography.ECDSAcurve().Verify(challenge.ToString(), responseMsg.Response.ToString())
}

func completeHandshake(connection net.Conn) {
	challenge, err := readMsg(connection)
	if err != nil {
		logging.Log("Error reading handshake challenge:"+err.Error(), "networking", "ERROR")
		return
	}
	logging.Log("Received challenge: "+challenge.Content, "networking", "INFO")

	response := cryptography.ECDSAcurve().Sign(cryptography.MIntFromString(challenge.Content).ToString())
	sendMsg(connection, MsgHandshakeResponse{
		Response: cryptography.MIntFromString(response),
	}.ToMsg().ToBytes())
}

func readMsg(connection net.Conn) (*Message, error) {
	buffer := make([]byte, 1024)
	mLen, err := connection.Read(buffer)
	if err != nil {
		logging.Log("Error reading: "+err.Error(), "networking", "ERROR")
		return nil, err
	}
	msg := Message{}
	err = json.Unmarshal(buffer[:mLen], &msg)
	if err != nil {
		logging.Log("Error deserializing message: "+err.Error(), "networking", "ERROR")
		return nil, err
	}

	if !msg.Validate() {
		logging.Log("Invalid message received", "networking", "ERROR")
		return nil, errors.New("Invalid signature")
	}

	logging.Log("Received: "+msg.Content, "networking", "INFO")

	return &msg, nil
}

func sendMsg(connection net.Conn, msg []byte) {
	_, err := connection.Write(msg)
	if err != nil {
		logging.Log("Error writing: "+err.Error(), "networking", "ERROR")
	}
}

func Broadcast(message string) {
	for i := 0; i < len(peersConnections); i++ {
		if peersConnections[i] != nil {
			peersConnections[i].Write([]byte(message))
		}
	}
}
