package main

import (
	"github.com/kirban/potato-db/internal/app/server"
	"log"
)

func main() {
	app, err := server.NewAppServer()

	if err != nil {
		log.Fatal(err)
	}

	app.Run()
}
