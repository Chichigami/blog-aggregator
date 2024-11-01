package main

import (
	"log"
	"os"

	"github.com/chichigami/blog-aggregator/internal/config"
)

func main() {
	configFile, err := config.GetConfigFilePath()
	if err != nil {
		log.Fatalf("error reading config: %v", err)
	}
	cfg := config.ReadAndParse(configFile)
	userCfg := state{
		cfg: &cfg,
	}
	cmds := commands{
		handlers: make(map[string]func(*state, command) error),
	}
	cmds.register("login", handlerLogin)
	input := os.Args
	if len(input) < 2 {
		log.Fatal("not enough arguments")
	}
	commandInput := command{
		name: input[1],
		args: input[2:],
	}
	err = cmds.run(&userCfg, commandInput)
	if err != nil {
		log.Fatal(err)
	}
}
