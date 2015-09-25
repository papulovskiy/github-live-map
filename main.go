package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

const (
	ApiUrl = "https://api.github.com/events"
)

type Event struct {
	Id	string			`json:"id",string`
	Type	string			`json:"type"`
	Actor	interface{}		`json:"actor"`
	Repo	interface{}		`json:"repo"`
	Payload	interface{}		`json:"payload"`
	Public	bool			`json:"public"`
	Created_at	string		`json:"created_at"`
	Org	interface{}		`json:"org,omitempty"`
}
type ApiResponse struct {
	Events []Event
}

func ReadEvents(result *ApiResponse) error {
	res, err := http.Get(ApiUrl)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return err
	}
	err = json.Unmarshal(body, &result.Events)
	if err != nil {
		return err
	}
	return nil
}

func ProcessEvents(result *ApiResponse) error {
	for index, e := range result.Events {
		fmt.Printf("%d %+v\n", index, e.Type)
	}
	return nil
}

func Loop() error {
	var r ApiResponse 
	for {
		err := ReadEvents(&r)
		if err != nil {
			panic(err)
		}
		ProcessEvents(&r)
	}
	return nil
}

func main() {
	Loop()
}
