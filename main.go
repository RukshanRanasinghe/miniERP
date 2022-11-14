package main

import (
	"fmt"

	"examples.com/grocery/api"
	"examples.com/grocery/config"
)

func main() {
	fmt.Println("Welcome to Peter's Grocery")
	config.Connection()
	api.InitiateServer()
}
