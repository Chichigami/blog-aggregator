package main

import (
	"fmt"

	"github.com/chichigami/blog-aggregator/internal/config"
)

func main() {
	configFile, err := config.GetConfigFilePath()
	if err != nil {
		fmt.Println("GetConfigFilePath error")
	}
	rcfg := config.ReadAndParse(configFile)
	rcfg.SetUser("chichigami")
}
