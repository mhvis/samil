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
	fmt.Println("model info:", s.ModelInfo())
	fmt.Printf("data: %+v\n", s.Data())
}
