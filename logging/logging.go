package logging

import (
	"fmt"
	"maxchain/config"
	"os"
	"time"
)

var file *os.File
var level string

func Init(config config.Configuration) {
	level = config.LogLevel

	if compareLevels("INFO") {
		fmt.Println("Initializing logging package")
	}

	tempFile, err := os.OpenFile(config.LogFilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		panic("Cannot create log file: " + err.Error())
	}

	file = tempFile

	Log("Starting MaxChain", "logging", "INFO")
	Log("Initialized logging package", "logging", "INFO")
}

func compareLevels(lvl string) bool {
	switch level {
	case "DEBUG":
		return true
	case "INFO":
		if lvl == "INFO" || lvl == "DEBUG" {
			return true
		}
	case "ERROR":
		if lvl == "ERROR" {
			return true
		}
	}
	return false
}

func Log(message string, source string, level string) {
	if !compareLevels(level) {
		return
	}

	str := time.Now().String() + " [" + level + "] " + source + ": " + message + "\n"
	fmt.Print(str)
	file.WriteString(str)
}

func PanicWithLog(message string, source string) {
	Log(message, source, "ERROR")
	panic(message)
}
