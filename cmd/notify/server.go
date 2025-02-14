package main

import (
	"log"

	"github.com/Irooniam/sotailc/internal/server"
)

func main() {
	//how notify service and ws server will communicate
	ch := make(chan string)

	connstr := "postgres://graph:graph@127.0.0.1/graph?sslmode=disable"
	runner := server.NewRunner(ch, connstr)
	log.Println("Starting to listen for events...")
	go runner.Run()

	ws := server.NewWS(ch)
	go ws.Heartbeat()
	go ws.ReceiveMsg()

	err := ws.Start()
	if err != nil {
		log.Fatal(err)
	}

	defer ws.Close()

}
