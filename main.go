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

func ReadEvents(result *interface{}) error {
	res, err := http.Get(ApiUrl)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return err
	}
	err = json.Unmarshal(body, &result)
	if err != nil {
		return err
	}
	return nil
}

func ProcessEvents(result *interface{}) error {
	fmt.Printf("%+v", *result)
	return nil
}

func Loop() error {
	var r interface{}
	for {
		ReadEvents(&r)
		ProcessEvents(&r)
	}
	return nil
}

func main() {
	Loop()
}
