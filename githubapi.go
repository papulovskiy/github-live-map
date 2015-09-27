package main

import (
	"encoding/json"
	"fmt"
	"gopkg.in/redis.v3"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"
)

// TODO: Implement oAuth token here
const (
	ApiBase = "https://api.github.com/"
)

type CacheStat struct {
	Hit  int
	Miss int
}

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

func Reader(ch chan<- Event) error {
	var r ApiResponse
	var remaining, reset int64

	previousEvents := make(map[string]bool)
	for {
		currentEvents := make(map[string]bool)
		fmt.Printf("Reader iteration\n")
		err := ReadEvents(&r, &remaining, &reset)
		if err != nil {
			panic(err)
		}

		for _, e := range r.Events {
			if _, ok := previousEvents[e.Id]; !ok {
				ch <- e
				currentEvents[e.Id] = true
			}
		}
		previousEvents = currentEvents

		// Pause stuff to respect rate limits
		time_diff := reset - time.Now().Unix()
		if time_diff < 0 {
			time_diff = time_diff * -1
		}
		if remaining <= 0 {
			remaining = 1
			time_diff += 3 // just add a few seconds if we're close to rate limit
		}
		pause := int64(time_diff / remaining)
		if pause < 1 {
			// Let's not go faster than light, at least for now
			pause = 1
		}
		fmt.Printf("Reader pause %d seconds\n", pause)
		//fmt.Printf("%d %d %d %d %d\n", remaining, reset, time.Now().Unix(), time_diff, pause)
		time.Sleep(time.Duration(pause) * time.Second)
	}
	return nil
}

func ReadProfile(url string, actor *Actor) error {
	// TODO: implement Bloom filter here to avoid unnecessary Redis requests
	value, err := redisClient.Get("loc_" + actor.Login).Result()
	if err == nil {
		actor.Location = value
		redisStat.Hit++
		//		fmt.Printf("Cache hit: %+v, %+v\n", actor.Login, value)
		return nil
	}
	redisStat.Miss++
	res, err := http.Get(url)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return err
	}
	err = json.Unmarshal(body, &actor)
	if err != nil {
		return err
	}
	redisClient.Set("loc_"+actor.Login, actor.Location, 0)
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
		ReadProfile(event.Actor.Url, &event.Actor)
		fmt.Printf("Profile: %+v\n", event.Actor)
		//panic(event.Actor)
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

var redisClient *redis.Client
var redisStat CacheStat

func GitHubLoop(MessageCh chan<- Message) error {
	redisClient = redis.NewClient(&redis.Options{
		Addr:     conf.Redis,
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	EventCh := make(chan Event, 100) // Let's assume that if we're stuck on profile resolving, we shouldn't be reading a lot of new events
	go Reader(EventCh)
	go ProfileResolverLoop(EventCh, MessageCh)
	return nil
}
