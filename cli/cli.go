package main

import (
	"flag"
	"fmt"
	"github.com/mhvis/samil"
)

func main() {
	flag.Parse()
	fmt.Println("searching for inverter")
	s := newConnection()
	fmt.Println("connected to:", s.Samil.RemoteAddr())
	fmt.Printf("model: %+v\n", *s.Model())
	fmt.Printf("data: %+v\n", *s.Data())
	fmt.Println("history")
	c := make(chan samil.HistoryDay, 1)
	s.History(17, 17, c)
	for d := range c {
		fmt.Printf("%v/%v/%v\n", d.Month, d.Day, d.Year)
		fmt.Println(d.Value)
	}
}
