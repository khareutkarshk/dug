package main

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func ws(w http.ResponseWriter, r *http.Request) {

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return
	}
	defer conn.Close()

	log.Println("client connected")

	for {
		mt, msg, err := conn.ReadMessage()
		if err != nil {
			break
		}

		log.Printf("received: %s\n", msg)

		conn.WriteMessage(mt, msg)
	}
}

func main() {

	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	http.HandleFunc("/ws", ws)

	log.Println("Backend listening :3003")

	log.Fatal(http.ListenAndServe(":3003", nil))
}
