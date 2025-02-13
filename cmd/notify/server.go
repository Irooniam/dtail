package main

import (
	"log"

	"github.com/Irooniam/sotailc/internal/server"
)

func main() {
	/*
		connstr := "postgres://graph:graph@127.0.0.1/graph?sslmode=disable"
		runner := server.NewRunner(connstr)
		log.Println("Starting to listen for events...")
		runner.Run()
	*/

	ws := server.NewWS()
	go ws.Heartbeat()
	err := ws.Start()
	if err != nil {
		log.Fatal(err)
	}

	defer ws.Close()

}
