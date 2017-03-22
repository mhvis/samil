package main

import (
	"fmt"
	"github.com/mhvis/samil"
	"os"
)

// samilHandler is a layer over the Samil type that handles errors by exiting.
type samilHandler struct {
	samil.Samil
}

func (h samilHandler) ModelInfo() samil.ModelInfo {
	modelInfo, err := h.Samil.ModelInfo()
	h.handle(err, "model info")
	return modelInfo
}

func (h samilHandler) Data() samil.InverterData {
	data, err := h.Samil.Data()
	h.handle(err, "data")
	return data
}

func (h samilHandler) History(start, end int) {
	err := h.Samil.History(start, end)
	h.handle(err, "history")
	return
}

func newConnection() samilHandler {
	s, err := samil.NewConnection()
	h := samilHandler{s}
	h.handle(err, "search")
	return h
}

// Prints error, closes socket and exits.
func (h samilHandler) handle(err error, action string) {
	if err == nil {
		return
	}
	fmt.Fprintln(os.Stderr, action, "failed:", err)
	os.Exit(1)
}
