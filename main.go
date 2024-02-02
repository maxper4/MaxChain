package main

import (
	"fmt"
	//"maxchain/config"
	//"maxchain/networking"
	"maxchain/cryptography"
)


func main() {
	fmt.Println("Starting MaxChain...")
	// config, err := config.LoadConfiguration("config/config.json")
	// if err != nil {
	// 	fmt.Println(err.Error())
	// 	panic("Cannot load configuration")
	// }
	// networking.Init(config)
	// fmt.Println("MaxChain started")

	//testMaths()
	cryptography.TestCrypto()

	for true {
	}
}

func testMaths () {
	nb := cryptography.MIntFromString("123456789abcdef")
	fmt.Println(nb)
	fmt.Println(nb.ToString())
	fmt.Println(cryptography.MIntFromString("123456789abcdef").Add(cryptography.MIntFromString("123456789abcdef")).ToString())

	fmt.Println(cryptography.MIntFromString("123").Add(cryptography.MIntFromString("456")).ToString())
	fmt.Println(cryptography.MIntFromString("a").Add(cryptography.MIntFromString("a")).ToString())
	fmt.Println(cryptography.MIntFromString("1a").Add(cryptography.MIntFromString("a")).ToString())
	fmt.Println(cryptography.MIntFromString("a").Add(cryptography.MIntFromString("1a")).ToString())

	fmt.Println(cryptography.MIntFromString("123456789abcdef").Multi(2).ToString())
	fmt.Println(cryptography.MIntFromString("123456789abcdef").Mult(cryptography.MIntFromString("ae13d")).ToString())
	fmt.Println(cryptography.MIntFromString("8fd").Mult(cryptography.MIntFromString("d")).ToString())
	fmt.Println(cryptography.MIntFromString("8fd").Eq(cryptography.MIntFromString("d")))
	fmt.Println(cryptography.MIntFromString("123456789abcdef").Multi(2).Eq(cryptography.MIntFromString("2468acf13579bde")))
	fmt.Println(cryptography.MIntFromString("123456789abcdef").GreaterEq(cryptography.MIntFromString("123456789abcdef")))
	fmt.Println(cryptography.MIntFromString("123456789abcdef").GreaterEq(cryptography.MIntFromString("223456789abcdef")))
	fmt.Println(cryptography.MIntFromString("a123456789abcdef").GreaterEq(cryptography.MIntFromString("123456789abcdef")))
	fmt.Println(cryptography.MIntFromString("223456789abcdef").GreaterEq(cryptography.MIntFromString("123456789abcdef")))
	fmt.Println(cryptography.MIntFromString("123456789abcdef").Sub(cryptography.MIntFromString("123456789abcde0")).ToString())
	fmt.Println(cryptography.MIntFromString("123456789abcdef").Sub(cryptography.MIntFromString("123456789abcdef")).ToString())
	fmt.Println(cryptography.MIntFromString("123456789bbcdef").Sub(cryptography.MIntFromString("123456789abcdef")).ToString())
	fmt.Println(cryptography.MIntFromString("12345678aabcdef").Sub(cryptography.MIntFromString("123456789bbcdef")).ToString())
}