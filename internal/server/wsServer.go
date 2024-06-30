package server

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"nhooyr.io/websocket"
)

type WsServer struct {
	entering chan *websocket.Conn
	leaving  chan *websocket.Conn
	Messages chan interface{}
}

func New(messages chan interface{}) *WsServer {
	return &WsServer{
		entering: make(chan *websocket.Conn),
		leaving:  make(chan *websocket.Conn),
		Messages: messages,
	}
}

func (ws *WsServer) WsHandler(w http.ResponseWriter, r *http.Request) {
	c, err := websocket.Accept(w, r, &websocket.AcceptOptions{
		InsecureSkipVerify: true,
	})
	if err != nil {
		log.Println(err)
		return
	}
	defer c.Close(websocket.StatusInternalError, "StatusInternalError")

	log.Printf("entering %s", r.RemoteAddr)
	ws.entering <- c

	for {
		if websocket.CloseStatus(err) == websocket.StatusNormalClosure {
			log.Printf("leaving %s", r.RemoteAddr)
			ws.leaving <- c
			return
		}
	}
}

func (ws *WsServer) Broadcaster() {
	clients := make(map[*websocket.Conn]bool) // all connected clients
	for {
		select {
		case msg := <-ws.Messages:
			// Broadcast incoming messages to all
			// clients' outgoing message channels.
			for cli := range clients {
				json, _ := json.Marshal(msg)
				cli.Write(context.Background(), websocket.MessageText, json)
			}
		case cli := <-ws.entering:
			clients[cli] = true
			log.Printf("clients %d", len(clients))

			// every 5 seconds we try to send "ping" message, and if it fails
			// we remove that client
			go func(c *websocket.Conn) {
				for {
					err := c.Write(context.Background(), websocket.MessageText, []byte("ping"))
					if err != nil {
						log.Printf("ping error: %v", err)
						ws.leaving <- c
						return
					}
					time.Sleep(5 * time.Second)
				}
			}(cli)
		case cli := <-ws.leaving:
			log.Printf("leaving")
			delete(clients, cli)
			log.Printf("clients %d", len(clients))
		}
	}
}
