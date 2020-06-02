package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
	"github.com/streadway/amqp"
)

//Notifier asdff
type Notifier struct {
	Connections map[string]*websocket.Conn
	lock        sync.Mutex
}

//User ajsdlfkj
type User struct {
	Nickname  string
	SessionID string
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		// return r.Header.Get("Origin") == "https://website.me"
		return true
	},
}

//InsertConnection inserts ws connection
func (n *Notifier) InsertConnection(conn *websocket.Conn, id string) {
	n.lock.Lock()
	defer n.lock.Unlock()
	if len(n.Connections) == 0 {
		n.Connections = make(map[string]*websocket.Conn)
	}
	n.Connections[id] = conn
}

//RemoveConnection removes ws connection given sessionid
func (n *Notifier) RemoveConnection(id string) {
	n.lock.Lock()
	defer n.lock.Unlock()
	delete(n.Connections, id)
}

//WriteToConnections writes to all connections
func (n *Notifier) WriteToConnections(msgs <-chan amqp.Delivery) {
	for msg := range msgs {
		n.lock.Lock()
		byteMsg := []byte(msg.Body)
		for id, conn := range n.Connections {
			if err := conn.WriteMessage(websocket.TextMessage, byteMsg); err != nil {
				n.RemoveConnection(id)
				conn.Close()
			}
		}
		msg.Ack(false)
		n.lock.Unlock()
	}
}

//WsHandler fjdlskfj
func (ctx *HandlerContext) WsHandler(w http.ResponseWriter, r *http.Request) {
	//check origin?

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		http.Error(w, "Failed to open websocket connection", 401)
		return
	}

	//auth check?

	xuserJSON := &User{}
	xuserStr := r.URL.Query().Get("X-user")
	if err := json.Unmarshal([]byte(xuserStr), xuserJSON); err != nil {
		http.Error(w, "Error unmarshaling X-User into JSON", 500)
	}

	ctx.Notifier.InsertConnection(conn, xuserJSON.SessionID)

	go (func(conn *websocket.Conn, ctx *HandlerContext, id string) {
		defer conn.Close()
		defer ctx.Notifier.RemoveConnection(id)

		for {
			messageType, p, err := conn.ReadMessage()
			if messageType == websocket.TextMessage || messageType == websocket.BinaryMessage {
				xuserJSON := &User{}
				xuserStr := r.URL.Query().Get("X-user")
				if err := json.Unmarshal([]byte(xuserStr), xuserJSON); err != nil {
					http.Error(w, "Error unmarshaling X-User into JSON", 500)
				}
				fmt.Printf("test: %s: %s", xuserJSON.Nickname, string(p))

			} else if messageType == websocket.CloseMessage || err != nil {
				break
			}
		}
	})(conn, ctx, xuserJSON.SessionID)
}
