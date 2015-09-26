package main

import (
	"fmt"
)

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

func main() {
	ch := make(chan Message)
	go GitHubLoop(ch)
	for {
		e := <-ch
		fmt.Printf("%+v %+v\n", e.User.Login, e.Type)
	}
}
