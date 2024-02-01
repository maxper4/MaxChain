package config

import (
    "encoding/json"
    "os"
	"fmt"
)

type Peer struct {
	Name string
	Ip string
	Port int
}

type Configuration struct {
	ListeningPort int
	Peers []Peer
}

func LoadConfiguration(file string) (Configuration, error) {
	var config Configuration
	configFile, err := os.Open(file)
	defer configFile.Close()
	if err != nil {
		fmt.Println(err.Error())
		return config, err
	}
	jsonParser := json.NewDecoder(configFile)
	jsonParser.Decode(&config)	
	return config, err
}