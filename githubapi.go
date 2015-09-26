package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"
)

// TODO: Implement oAuth token here 
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

func ReadEvents(result *ApiResponse, remaining, reset *int64) error {
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
	*remaining, _ = strconv.ParseInt(res.Header["X-Ratelimit-Remaining"][0], 10, 64)
	*reset, _     = strconv.ParseInt(res.Header["X-Ratelimit-Reset"][0], 10, 64)
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
	var remaining int64
	var reset int64
	for {
		err := ReadEvents(&r, &remaining, &reset)
		if err != nil {
			panic(err)
		}
		ProcessEvents(&r)
		time_diff := reset - time.Now().Unix()
		if time_diff < 0 {
			time_diff = time_diff * -1
		}
		if remaining <= 0 {
			remaining = 1
			time_diff += 3 // just add a few seconds if we're close to rate limit
		}
		pause := int64(time_diff / remaining)
		//fmt.Printf("%d %d %d %d %d\n", remaining, reset, time.Now().Unix(), time_diff, pause)
		time.Sleep(time.Duration(pause) * time.Second)
	}
	return nil
}

