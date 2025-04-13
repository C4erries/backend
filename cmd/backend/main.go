package main

import (
	"backend/internal/database"
	"backend/internal/server"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	if err := database.Connect(); err != nil {
		panic(err)
	}
	go server.Start()
	//graceful shutdown, когда система вызовет системный сигнал SIGTERM, он запишется в канал,
	//затем с канала будет прочитан сигнал и код функции main продолжится
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)

	<-stop

	database.MustClose()
}
