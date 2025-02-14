package server

import (
	"context"
	"fmt"
	"io"

	"log"

	"golang.org/x/time/rate"

	"net/http"

	"sync"
	"time"

	"github.com/coder/websocket"
)

type WS struct {
	conns  *sync.Map //key will be *websocket.Conn
	server *http.Server
	ch     chan string
}

func (ws *WS) hand(w http.ResponseWriter, r *http.Request) {
	wsKey := r.Header.Get("Sec-WebSocket-Key")
	rt := context.WithValue(r.Context(), "websocket-key", wsKey)
	c, err := websocket.Accept(w, r.WithContext(rt), &websocket.AcceptOptions{
		Subprotocols:       []string{"echo"},
		InsecureSkipVerify: true,
	})
	if err != nil {
		log.Println("error in handling websocket ", err)
		log.Println(err)
		return
	}

	ws.conns.Store(c, c)

	//rate limiter
	l := rate.NewLimiter(rate.Every(time.Millisecond*10), 2)

	//start loop for messages
	for {
		err = ws.pusher(rt, c, l)
		if err != nil {
			log.Printf("failed to echo with %s %s", r.RemoteAddr, err)
			ws.disconnect(c)
			return
		}
	}

}

func NewWS(ch chan string) *WS {
	conns := &sync.Map{}
	connuri := fmt.Sprintf("0.0.0.0:9999")

	ws := WS{}
	ws.conns = conns
	ws.ch = ch

	webserver := &http.Server{}
	webserver.Addr = connuri
	webserver.Handler = http.HandlerFunc(ws.hand)

	ws.server = webserver

	return &ws
}

func (ws *WS) Close() {
	log.Println("closing ws server...")
	ws.server.Close()
}

func (ws *WS) disconnect(conn *websocket.Conn) {
	log.Println("closing connection for client ", conn)

	err := conn.Close(websocket.StatusNormalClosure, "good bye")
	if err != nil {
		log.Println(err)
	}

	ws.conns.Delete(conn)
	//remove from global registry
	//Conns.Delete(client.Conn)

}
func (ws *WS) Start() error {
	log.Println("starting ws...")
	err := ws.server.ListenAndServe()
	if err != nil {
		return err
	}

	return nil
}

func (ws *WS) ReceiveMsg() {
	log.Println("Receive Msg running...")
	var err error
	for msg := range ws.ch {
		if err = ws.SendMsg(msg); err != nil {
			log.Println("tried sending messages to client but to ", err)
			continue
		}
	}
}

func (ws *WS) SendMsg(s string) error {
	var err error = nil
	var i int
	ws.conns.Range(func(k, v interface{}) bool {
		log.Printf("Sending msg to client %s with payload %s", v.(*websocket.Conn).Subprotocol(), s)
		err = v.(*websocket.Conn).Write(context.Background(), websocket.MessageText, []byte(s))
		if err != nil {
			log.Println(err)
			ws.disconnect(v.(*websocket.Conn))
			return false
		}
		i++
		return true
	})

	//only log if we have connections
	if i > 0 {
		log.Printf("sent message to %d clients\n", i)
	}
	return err
}

func (ws *WS) Heartbeat() error {
	for {
		log.Println("send heartbeat to all clients")
		if err := ws.SendMsg("h"); err != nil {
			log.Println(err)
		}

		time.Sleep(time.Second * 5)
	}

}

/*
this the routine responsible for sending / receiving from viewer client
*/
func (ws *WS) pusher(ctx context.Context, c *websocket.Conn, l *rate.Limiter) error {
	//keep connection open for 1 day
	ctx, cancel := context.WithTimeout(ctx, time.Second*86400)
	defer cancel()

	log.Println(" ctx in pusher ", ctx.Value("websocket-key"))
	err := l.Wait(ctx)
	if err != nil {
		return err
	}

	typ, r, err := c.Reader(ctx)
	if err != nil {
		return err
	}

	b, _ := io.ReadAll(r)
	log.Println("received message ", typ, string(b))

	//only thing we accept are heartbeats
	if string(b) == "h" {
		return nil
	}

	return nil
}
