package networking

import (
	"encoding/json"
	"maxchain/cryptography"
	"maxchain/logging"
)

type Message struct {
	Tag       uint16
	Sender    string
	Signature string
	Content   string
}

func (message Message) ToBytes() []byte {
	serialized, err := json.Marshal(message)
	if err != nil {
		logging.Log("Error serializing message: "+err.Error(), "networking", "ERROR")
	}
	return []byte(string(serialized))
}

func (message Message) Validate() bool {
	return true // TODO: implement
}

type MsgHandshakeChallenge struct {
	Challenge cryptography.MInt
}

func (msg MsgHandshakeChallenge) ToMsg() Message {
	serialized, err := json.Marshal(&struct {
		Challenge string
	}{Challenge: msg.Challenge.ToString()})

	if err != nil {
		logging.PanicWithLog("Error serializing message: "+err.Error(), "networking")
	}

	return Message{
		Tag:       0,
		Sender:    "me",
		Signature: "123456789abcdef",
		Content:   string(serialized),
	}
}

func (msg Message) ToMsgHandshakeChallenge() (MsgHandshakeChallenge, error) {
	var content struct {
		Challenge string
	}
	err := json.Unmarshal([]byte(msg.Content), &content)
	if err != nil {
		logging.Log("Error deserializing message: "+err.Error(), "networking", "ERROR")
		return MsgHandshakeChallenge{}, err
	}

	return MsgHandshakeChallenge{
		Challenge: cryptography.MIntFromString(content.Challenge),
	}, nil
}

type MsgHandshakeResponse struct {
	Response cryptography.MInt
}

func (msg MsgHandshakeResponse) ToMsg() Message {
	serialized, err := json.Marshal(&struct {
		Response string
	}{Response: msg.Response.ToString()})

	if err != nil {
		logging.PanicWithLog("Error serializing message: "+err.Error(), "networking")
	}

	return Message{
		Tag:       1,
		Sender:    "me",
		Signature: "123456789abcdef",
		Content:   string(serialized),
	}
}

func (msg Message) ToMsgHandshakeResponse() (MsgHandshakeResponse, error) {
	var content struct {
		Response string
	}
	err := json.Unmarshal([]byte(msg.Content), &content)
	if err != nil {
		logging.Log("Error deserializing message: "+err.Error(), "networking", "ERROR")
		return MsgHandshakeResponse{}, err
	}

	return MsgHandshakeResponse{
		Response: cryptography.MIntFromString(content.Response),
	}, nil
}
