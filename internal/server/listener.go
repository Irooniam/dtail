package server

import (
	"log"
	"time"

	"github.com/lib/pq"
)

type Runner struct {
	pl *pq.Listener
	ch chan string
}

func (r *Runner) Errors(ev pq.ListenerEventType, err error) {
	if err != nil {
		log.Printf("Error type %d: %s", ev, err)
	}
}

func (r *Runner) Run() {
	err := r.pl.Listen("dtail_table_update")
	if err != nil {
		panic(err)
	}

	for {
		select {
		case event := <-r.pl.Notify:
			log.Println("DB notification: ", event.Extra)
			r.ch <- string(event.Extra)
		//havent received events - ping db to make sure its all good
		case <-time.After(60 * time.Second):
			go r.pl.Ping()
			log.Println("havent received any events - pinging db")
		}

	}

}

func runnerError(ev pq.ListenerEventType, err error) {
	if err != nil {
		log.Printf("Error type %d: %s", ev, err)
	}
}

func NewRunner(ch chan string, connstr string) *Runner {
	minReInterval := 5 * time.Second
	maxReInterval := 2 * time.Minute

	return &Runner{ch: ch, pl: pq.NewListener(connstr, minReInterval, maxReInterval, runnerError)}
}
