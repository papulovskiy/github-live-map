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
	ApiBase = "https://api.github.com/"
)

type Actor struct {
	Id          int64  `json:"id"`
	Login       string `json:"login"`
	Gravatar_id string `json:"gravatar_id"`
	Url         string `json:"url"`
	Avatar_url  string `json:"avatar_url"`
	Location    string `json:"location",omitempty`
}

type Event struct {
	Id         string      `json:"id",string`
	Type       string      `json:"type"`
	Actor      Actor       `json:"actor"`
	Repo       interface{} `json:"repo"`
	Payload    interface{} `json:"payload"`
	Public     bool        `json:"public"`
	Created_at string      `json:"created_at"`
	Org        interface{} `json:"org,omitempty"`
}
type ApiResponse struct {
	Events []Event
}

func ReadEvents(result *ApiResponse, remaining, reset *int64) error {
	res, err := http.Get(ApiBase + "events")
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
	*reset, _ = strconv.ParseInt(res.Header["X-Ratelimit-Reset"][0], 10, 64)
	return nil
}

func ProcessEvents(result *ApiResponse, ch chan<- Event) error {
	for _, e := range result.Events {
		ch <- e
	}
	return nil
}

func Reader(ch chan<- Event) error {
	var r ApiResponse
	var remaining, reset int64

	for {
		fmt.Printf("Reader iteration\n")
		err := ReadEvents(&r, &remaining, &reset)
		if err != nil {
			panic(err)
		}
		ProcessEvents(&r, ch)
		time_diff := reset - time.Now().Unix()
		if time_diff < 0 {
			time_diff = time_diff * -1
		}
		if remaining <= 0 {
			remaining = 1
			time_diff += 3 // just add a few seconds if we're close to rate limit
		}
		pause := int64(time_diff / remaining)
		fmt.Printf("Reader pause %d seconds\n", pause)
		//fmt.Printf("%d %d %d %d %d\n", remaining, reset, time.Now().Unix(), time_diff, pause)
		time.Sleep(time.Duration(pause) * time.Second)
	}
	return nil
}

// So... Here we need to fetch user's location or check in the cache
// I'd prefer to make two layers of cache:
//	* small in-memory to handle most recent active users
//	* bigger cache based on redis to reduce amount of API calls
func ProfileResolverLoop(ProfileCh <-chan Event, MessageCh chan<- Message) error {
	//	localUsers := make(map[string]*User)
	for {
		event := <-ProfileCh
		fmt.Printf("Profile Resolver loop: %+v %+v\n", event.Actor.Login, event.Type)
		var m Message
		var u User
		u.Id = event.Actor.Id
		u.Login = event.Actor.Login
		u.Avatar = "" // TODO: load avatar and send it as base64 or just skip it
		m.EventId = event.Id
		m.Type = event.Type
		m.User = u
		MessageCh <- m
		time.Sleep(time.Duration(3) * time.Second) // sleep here is only for testing purpose
	}
	return nil
}

/*
	This loop is to buffer events from Reader and pass them one by one to profile resolver
	because channels operations in go are blocking
*/
/*
	This function was removed because I found great Go's feature: buffered channels!
*/


func GitHubLoop(MessageCh chan<- Message) error {
	EventCh := make(chan Event, 100) // Let's assume that if we're stuck on profile resolving, we shouldn't be reading a lot of new events
	go Reader(EventCh)
	go ProfileResolverLoop(EventCh, MessageCh)
	return nil
}
