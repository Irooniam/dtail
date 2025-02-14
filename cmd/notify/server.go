package main

import (
	"log"
	"os"

	"github.com/Irooniam/sotailc/internal/server"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Cant find .env file")
	}

	var connstr string
	var notifych string

	if connstr = os.Getenv("DTAIL_DB_CONFIG"); connstr == "" {
		log.Fatal("env var DTAIL_DB_CONFIG cant be empty")
		return
	}

	if notifych = os.Getenv("DTAIL_LISTEN_CHANNEL"); connstr == "" {
		log.Fatal("env var DTAIL_LISTEN_CHANNEL cant be empty")
		return
	}

	//how notify service and ws server will communicate
	ch := make(chan string)
	runner := server.NewRunner(ch, connstr, notifych)
	log.Println("Starting to listen for events...")
	go runner.Run()

	ws := server.NewWS(ch)
	go ws.Heartbeat()
	go ws.ReceiveMsg()

	err = ws.Start()
	if err != nil {
		log.Fatal(err)
	}

	defer ws.Close()

}
