package main

import (
	"fmt"
	"github.com/gorilla/websocket"
	"net/http"
	"io/ioutil"
)

type msg struct {
	Num int
}

type User struct {
	Id int64
	Login string
	Location string
	Avatar string
}

type Message struct {
	EventId string
	Type string
	User User
	Latitude float64
	Longitude float64
}

func rootHandler(w http.ResponseWriter, r *http.Request) {
	content, err := ioutil.ReadFile("index.html")
	if err != nil {
		fmt.Println("Could not open file.", err)
	}
	fmt.Fprintf(w, "%s", content)
}

func wsHandler(w http.ResponseWriter, r *http.Request) {
//	if r.Header.Get("Origin") != "http://"+r.Host {
//		http.Error(w, "Origin not allowed", 403)
//		return
//	}
	conn, err := websocket.Upgrade(w, r, w.Header(), 1024, 1024)
	if err != nil {
		http.Error(w, "Could not open websocket connection", http.StatusBadRequest)
		return
	}

	go echo(conn)
}

func echo(conn *websocket.Conn) {
	for {
		m := msg{}

//		err := conn.ReadJSON(&m)
//		if err != nil {
//			fmt.Println("Error reading json.", err)
//			conn.Close()
//			return
//		}

		fmt.Printf("Got message: %#v\n", m)

		if err := conn.WriteJSON(m); err != nil {
			fmt.Println(err)
			conn.Close()
			return
		}
	}
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
func main() {
	ch := make(chan Message)
	go GitHubLoop(ch)
	pool := make(map[int]*websocket.Conn)
	go mux(ch, pool)
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		conn, err := websocket.Upgrade(w, r, w.Header(), 1024, 1024)
		if err != nil {
			http.Error(w, "Could not open websocket connection", http.StatusBadRequest)
			return
		}
		client_id++
		pool[client_id] = conn
	})
	http.HandleFunc("/", rootHandler)
	panic(http.ListenAndServe(":8080", nil))
}
