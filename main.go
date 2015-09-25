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

func ReadEvents(result interface{}) error {
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

func main() {
	var p interface{}
	ReadEvents(&p)
	fmt.Printf("%+v", p)
}
