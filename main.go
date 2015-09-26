package main

import (
	"fmt"
)

func main() {
	ch := make(chan Event)
	go Loop(ch)
	for {
		e := <-ch
		fmt.Printf("%+v %+v\n", e.Actor, e.Type)
	}
}
