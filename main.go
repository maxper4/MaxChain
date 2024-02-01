package main

import (
	"fmt"
	"maxchain/config"
	"maxchain/networking"
	"maxchain/cryptography"
)


func main() {
	fmt.Println("Starting MaxChain...")
	config, err := config.LoadConfiguration("config/config.json")
	if err != nil {
		fmt.Println(err.Error())
		panic("Cannot load configuration")
	}
	networking.Init(config)
	fmt.Println("MaxChain started")
	for true {
	}
}

func testMaths () {
	nb := cryptography.mIntFromString("123456789abcdef")
	fmt.Println(nb)
	fmt.Println(nb.ToString())
	fmt.Println(cryptography.mIntFromString("123456789abcdef").Add(cryptography.mIntFromString("123456789abcdef")).ToString())

	fmt.Println(cryptography.mIntFromString("123").Add(cryptography.mIntFromString("456")).ToString())
	fmt.Println(cryptography.mIntFromString("a").Add(cryptography.mIntFromString("a")).ToString())
	fmt.Println(cryptography.mIntFromString("1a").Add(cryptography.mIntFromString("a")).ToString())
	fmt.Println(cryptography.mIntFromString("a").Add(cryptography.mIntFromString("1a")).ToString())
}