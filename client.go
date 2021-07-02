package main

import(
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)



// serveWs handles websocket requests from the peer.
func serveWs(w http.ResponseWriter, r *http.Request, roomId string) {
	fmt.Print(roomId)
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err.Error())
		return
	}
	c := &connection{send: make(chan []byte, 256), ws: ws}
	s := subscription{c, roomId}
	h.register <- s
	go s.writePump()
	go s.readPump()
}

// Upgrader specifies parameters for upgrading an HTTP connection to a WebSocket connection

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
 }


// Defining the connection type

type connection struct {
   ws *websocket.Conn
   send chan []byte
}

// The readPump and writePump are the simple functions to pump the message from socket to hub and from hub to socket

func (s subscription) readPump() {
	c := s.conn
	defer func() {
	   Hub.unregister <- s
	   c.ws.Close()
	}()
	c.ws.SetReadLimit(maxMessageSize)
	c.ws.SetReadDeadline(time.Now().Add(pongWait))
	c.ws.SetPongHandler(func(string) error { c.ws.SetReadDeadline(time.Now().Add(pongWait)); return nil })
	for {
	   _, msg, err := c.ws.ReadMessage()
	   if err != nil {
		  if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway) {
			 log.Printf("error: %v", err)
		  }
		  break
	   }
	   m := message{msg, s.room}
	   Hub.broadcast <- m
	}
 }
