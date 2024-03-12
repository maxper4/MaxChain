package config

import (
	"encoding/json"
	"fmt"
	"os"
)

type Peer struct {
	Name string
	Ip   string
	Port int
}

type Configuration struct {
	ListeningPort int
	Peers         []Peer
	LogFilePath   string
	LogLevel      string
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

	valid, errmsg := validateConfig(config)
	if !valid {
		return config, fmt.Errorf("Invalid configuration: %s", errmsg)
	}

	return config, err
}

func validateConfig(config Configuration) (bool, string) {
	if config.ListeningPort < 1024 || config.ListeningPort > 65535 {
		return false, "Invalid listening port"
	}
	for _, peer := range config.Peers {
		if peer.Port < 1024 || peer.Port > 65535 {
			return false, "Invalid peer port for peer " + peer.Name
		}
	}
	if config.LogLevel != "DEBUG" && config.LogLevel != "INFO" && config.LogLevel != "ERROR" {
		return false, "Invalid log level: only DEBUG, INFO and ERROR are allowed"
	}
	return true, ""
}
