package main

import (
	"fmt"
	"github.com/gorilla/websocket"
	"io/ioutil"
	"net/http"
)

func rootHandler(w http.ResponseWriter, r *http.Request) {
	content, err := ioutil.ReadFile("index.html")
	if err != nil {
		fmt.Println("Could not open file.", err)
	}
	fmt.Fprintf(w, "%s", content)
}

func mux(ch chan Message, pool map[int]*websocket.Conn) {
	for {
		m := <-ch
		for _, conn := range pool {
			conn.WriteJSON(m)
		}
	}
}

var client_id int = 0

func wsLoop(port string, uri string, message chan Message) {
	pool := make(map[int]*websocket.Conn)
	go mux(message, pool)
	http.HandleFunc("/"+uri, func(w http.ResponseWriter, r *http.Request) {
		conn, err := websocket.Upgrade(w, r, w.Header(), 1024, 1024)
		if err != nil {
			http.Error(w, "Could not open websocket connection", http.StatusBadRequest)
			return
		}
		client_id++
		pool[client_id] = conn
	})
	http.HandleFunc("/", rootHandler)
	panic(http.ListenAndServe(":"+port, nil))

}
