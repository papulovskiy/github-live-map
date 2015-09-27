package main

import (
	"fmt"
	"encoding/json"
	"os"
)

type Configuration struct {
	EventsApiToken	string
	ProfilesApiToken	string
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

func ReadConfig() Configuration {
	file, _ := os.Open("config/app/app.json")
	decoder := json.NewDecoder(file)
	configuration := Configuration{}
	err := decoder.Decode(&configuration)
	if err != nil {
		fmt.Println("error:", err)
	}
	return configuration
}


func main() {
	conf := ReadConfig()
	fmt.Printf("%+v\n", conf)
//	ch := make(chan Message)
//	go GitHubLoop(ch)
//	for {
//		e := <-ch
//		fmt.Printf("%+v %+v\n", e.User.Login, e.Type)
//	}
}
