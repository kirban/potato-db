package main

import (
	"fmt"
	"github.com/kirban/potato-db/internal/config"
	"log"
	"os"
)

func main() {
	configPath := os.Getenv("CONFIG_PATH")

	if configPath == "" {
		configPath = "config.potato.yaml"
	}

	cfg, err := config.NewConfig(configPath)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("%+v", cfg)
}
