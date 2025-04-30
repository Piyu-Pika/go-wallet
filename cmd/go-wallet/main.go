package main

import (
	"github.com/Piyu-Pika/go-wallet/internal/database"
	"github.com/Piyu-Pika/go-wallet/internal/routers"
)

func main() {
	database.InitDB()
	defer database.CloseDB()

	r := routers.SetupRouter()
	r.Run(":8080")
}
