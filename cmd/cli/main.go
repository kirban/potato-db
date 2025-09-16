package main

import (
	"bufio"
	"fmt"
	"github.com/kirban/potato-db/internal/db"
	"os"
)

func main() {
	database := db.NewDbBuilder().
		InitLogger().
		InitStorage().
		InitCompute().
		Build()

	reader := bufio.NewReader(os.Stdin)

	fmt.Printf("Enter command and then press enter\n")

	for {
		fmt.Printf("> ")
		query, _ := reader.ReadString('\n')

		result := database.ExecuteQuery(query)

		fmt.Println(result)
	}
}
