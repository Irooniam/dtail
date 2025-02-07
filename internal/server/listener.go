package server

import (
	"fmt"
	"log"
	"time"

	"github.com/lib/pq"
)

type Runner struct {
	*pq.Listener
}

func (r *Runner) Errors(ev pq.ListenerEventType, err error) {
	if err != nil {
		log.Printf("Error type %d: %s", ev, err)
	}
}

func (r *Runner) Run() {
	err := r.Listen("dtail_table_update")
	if err != nil {
		panic(err)
	}

	for {
		select {
		case event := <-r.Notify:
			fmt.Println("received notification, new work available", event.Extra)

		//havent received events - ping db to make sure its all good
		case <-time.After(60 * time.Second):
			go r.Ping()
			fmt.Println("havent received any events - pinging db")
		}

	}

}

func runnerError(ev pq.ListenerEventType, err error) {
	if err != nil {
		log.Printf("Error type %d: %s", ev, err)
	}
}

func NewRunner(connstr string) *Runner {
	minReInterval := 5 * time.Second
	maxReInterval := 2 * time.Minute

	return &Runner{pq.NewListener(connstr, minReInterval, maxReInterval, runnerError)}
}
