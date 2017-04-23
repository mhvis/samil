// A very simple command-line application for Samil Power inverters.
package main

import (
	"fmt"
	"github.com/mhvis/samil"
	"net"
	"os"
	"time"
)

func main() {
	firstRound := true
	for {
		if firstRound {
			fmt.Println("searching for inverters")
		} else {
			fmt.Println()
			fmt.Println("searching for another inverter")
		}

		inverter, err := samil.NewConnection()
		if e, ok := err.(net.Error); ok && e.Timeout() {
			// Stop application at I/O timeout
			if firstRound {
				fmt.Println("no inverters found")
			} else {
				fmt.Println("no new inverters found")
			}
			return
		}
		checkError(err, "connect")

		// Inverter found, print info
		fmt.Println("found inverter on address", inverter.RemoteAddr())
		model, err := inverter.Model()
		checkError(err, "model")
		fmt.Printf("model info: %+v\n", model)
		data, err := inverter.Data()
		checkError(err, "data")
		fmt.Printf("data info: %+v\n", data)

		// Keep inverter connected by sending keepalive packets
		// (to prevent reconnection to this inverter later)
		go func(inverter *samil.Samil) {
			for {
				time.Sleep(10 * time.Second)
				_, err := inverter.Data()
				checkError(err, "keepalive")
			}
		}(inverter)
		firstRound = false
	}
}

// Prints error and exits (sockets are automatically closed on exit).
func checkError(err error, action string) {
	if err == nil {
		return
	}
	fmt.Fprintln(os.Stderr, action, "failed:", err)
	os.Exit(1)
}
