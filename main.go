package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

func main() {
	url := "https://api.github.com/events"
	res, err := http.Get(url)
	if err != nil {
		// TODO: error handling
		panic(err)
	}
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		// TODO: error handling
		panic(err)
	}
	var p interface{}
	err = json.Unmarshal(body, &p)
	if err != nil {
		panic(err)
	}

	fmt.Printf("%+v", p)
}
