package main

import (
	"flag"
	"fmt"
)

func main() {
	flag.Parse()
	fmt.Println("searching for inverter")
	s := newConnection()
	fmt.Println("connected to:", s.Samil.RemoteAddr())
	fmt.Printf("model: %+v\n", *s.Model())
	fmt.Printf("data: %+v\n", *s.Data())
	//fmt.Println("history")
	//s.History(17,17)
}
