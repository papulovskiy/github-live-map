package main

import (
	"fmt"
	"encoding/json"
	"os"
)

type msg struct {
	Num int
}

type Configuration struct {
	EventsApiToken	string
	ProfilesApiToken	string
	WS	bool
	Port	string
	Uri	string
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
		// TODO: make better error handling
	}
	return configuration
}

func main() {
	conf := ReadConfig()
	fmt.Printf("%+v\n", conf)
	ch := make(chan Message)
	go GitHubLoop(ch)
	if conf.WS {
		wsLoop(conf.Port, conf.Uri, ch)
	} else {
		for {
			m := <- ch
			fmt.Printf("%+v\n", m)
		}
	}
}
