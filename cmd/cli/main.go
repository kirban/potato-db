package main

import (
	"github.com/kirban/potato-db/internal/app/cli"
	"log"
)

func main() {
	app, err := cli.NewCliApp()

	if err != nil {
		log.Fatal(err)
	}

	if err := app.Run(); err != nil {
		log.Fatal(err)
	}
}
