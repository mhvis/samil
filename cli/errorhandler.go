package main

import (
	"fmt"
	"github.com/mhvis/samil"
	"os"
)

// samilHandler is a layer over the Samil type that handles errors by exiting.
type samilHandler struct {
	*samil.Samil
}

func (h samilHandler) Model() *samil.Model {
	model, err := h.Samil.Model()
	handleError(err, "model")
	return model
}

func (h samilHandler) Data() *samil.Data {
	data, err := h.Samil.Data()
	handleError(err, "data")
	return data
}

func (h samilHandler) History(start, end int) {
	err := h.Samil.History(start, end)
	handleError(err, "history")
	return
}

func newConnection() samilHandler {
	s, err := samil.NewConnection()
	handleError(err, "search")
	return samilHandler{s}
}

// Prints error and exits (sockets are automatically closed on exit).
func handleError(err error, action string) {
	if err == nil {
		return
	}
	fmt.Fprintln(os.Stderr, action, "failed:", err)
	os.Exit(1)
}
