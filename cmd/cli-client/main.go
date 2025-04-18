package main

import (
	"fmt"
	"os"

	"github.com/Mikkkkkkka/typoracer/internal/cli"
)

var (
	command string
)

func initialise() {

	// TODO: load client-config

	if len(os.Args) < 2 {
		panic("The cli-client requires at least one parameter")
	}
	command = os.Args[1]
}

func main() {
	initialise()

	var err error
	switch command {
	case "health":
		cli.Health()
	case "type":
		panic("Not implemented!") // TODO: implement typing
	case "stats":
		panic("Not implemented!") // TODO: implement stats fetching
	case "register":
		err = cli.Register()
	default:
		fmt.Println("Sorry! Wrong command!")
		fmt.Println(command)
	}

	if err != nil {
		panic(err)
	}
}
