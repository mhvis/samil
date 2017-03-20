package main

import (
	"fmt"
	"github.com/mhvis/samil"
)

func main() {
	s, err := samil.NewConnection()
	if err != nil {
		fmt.Println(err)
		return
	}
	defer s.Close()
	// ModelInfo
	err = s.ModelInfo()
	if err != nil {
		fmt.Println(err)
		return
	}
	header, payload, end, err := s.Read()
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("Header: %v, payload: %v, end: %v\n", header, string(payload), end)
	// Status
	err = s.Status()
	if err != nil {
		fmt.Println(err)
		return
	}
	header, payload, end, err = s.Read()
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("Header: %v, payload: %v, end: %v\n", header, string(payload), end)
}
