package main

import (
	"fmt"
	"log"

	"github.com/Irooniam/dtail/internal/server"
	"github.com/lib/pq"
)

func problems(ev pq.ListenerEventType, err error) {
	log.Println("listen error ", err)
}

func main() {
	fmt.Println("vim-go")

	connstr := "postgres://graph:graph@127.0.0.1/graph?sslmode=disable"
	runner := server.NewRunner(connstr)
	runner.Run()

	/*
		db, err := internal.NewDB("postgres", connstr)
		if err != nil {
			panic(err)
		}
	*/

	/*
		minReconn := 10 * time.Second
		maxReconn := time.Minute
		listener := pq.NewListener(connstr, minReconn, maxReconn, problems)
		err := listener.Listen("dtail_table_update")
		if err != nil {
			panic(err)
		}

		for {
			select {
			case blah := <-listener.Notify:
				fmt.Println("received notification, new work available", blah.Extra)

			case <-time.After(90 * time.Second):
				go listener.Ping()
				// Check if there's more work available, just in case it takes
				// a while for the Listener to notice connection loss and
				// reconnect.
				fmt.Println("received no work for 90 seconds, checking for new work")
			}

		}
	*/
}
