package main

import (
	"fmt"
	"io"
	"log"
	"net/http"

	"golang.org/x/net/websocket"
)

var hub []*websocket.Conn

func main() {
	http.HandleFunc("/", index)
	http.Handle("/connws/", websocket.Handler(connWs))

	log.Fatal(http.ListenAndServe(":8080", nil))
}

func index(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "index.html")
}

func connWs(ws *websocket.Conn) {
	data := map[string]string{}
	hub = append(hub, ws)

	for {
		err := websocket.JSON.Receive(ws, &data)
		if err != nil && err != io.EOF {
			log.Println(err)
			ws.Close()
			break
		}

		message := fmt.Sprintf("%s: %s<br>", data["name"], data["text"])

		for i := range hub {
			err = websocket.Message.Send(hub[i], message)
			if err != nil {
				log.Println(err)
			}
		}
	}
}
