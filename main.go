package main

import (
	"log"
	"time"
	"user-server/server"
)

func main() {
	// server.Run()
	go server.Run()

	time.Sleep(time.Second)
	server.Test()

	time.Sleep(time.Second)
	log.Println("done")
}
