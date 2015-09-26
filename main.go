package main

import (
	"fmt"
)

func main() {
	ch := make(chan Event)
	go Loop(ch)
	e := <-ch
	fmt.Printf("%+v\n", e.Type)
}
