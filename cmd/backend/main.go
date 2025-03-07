package main

import (
	"backend/internal/database"
	"backend/internal/server"
)

func main() {
	database.DatabaseInit()
	server.Start()
}
