package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/Mikkkkkkka/typoracer/internal/cli"
	"github.com/Mikkkkkkka/typoracer/internal/config"
)

var (
	command string
)

func initialize() *config.CliConfig {

	file, err := os.ReadFile("config.json")
	if err != nil {
		generateConfigFile()
		file, err = os.ReadFile("config.json")
		if err != nil {
			panic("cannot create config file")
		}
	}
	var cfg config.CliConfig
	json.Unmarshal(file, &cfg)

	if len(os.Args) < 2 {
		panic("The cli-client requires at least one parameter")
	}
	command = os.Args[1]

	return &cfg
}

func generateConfigFile() {
	config := config.CliConfig{
		Host: "localhost",
		Port: "8080",
	}
	file, err := os.Create("config.json")
	if err != nil {
		panic(err)
	}
	defer file.Close()
	json.NewEncoder(file).Encode(config)
}

func main() {
	cfg := initialize()

	var err error
	switch command {
	case "health":
		cli.Health(cfg)
	case "play":
		panic("Not implemented!") // TODO: implement typing
	case "stats":
		cli.Stats(cfg)
	case "register":
		err = cli.Register(cfg)
	default:
		fmt.Println("Sorry! Wrong command!")
		fmt.Println(command)
	}

	if err != nil {
		panic(err)
	}
}
