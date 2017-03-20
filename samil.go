// Package samil provides an API for quering networked Samil Power inverters.
// Supported inverters are those that work with SolarPower Browser.
// These are: SolarRiver TD, SolarRiver TL-D and SolarLake TL inverters.
package samil

import (
	"net"
)

type Samil struct {
	conn   net.Conn
	read   chan message
	closed *error
}

// NewConnection searches for an inverter in the network and returns the
// connection if one is found.
//
// Inverters that are already connected to a client will not initiate a new
// connection. Therefore calling this function multiple times while leaving the
// connections open will connect to different inverters.
func NewConnection() (Samil, error) {
	conn, err := connect()
	if err != nil {
		return Samil{}, err
	}
	s := Samil{conn, make(chan message, 5), &err}
	go s.readRoutine()
	return s, nil
}

// ModelInfo requests and returns a string with model information.
// This API will change.
func (s Samil) ModelInfo() (string, error) {
	if *s.closed != nil {
		return "", *s.closed
	}
	_, err := s.conn.Write(modelInfo)
	if err != nil {
		return "", err
	}
	payload, err := s.readFor(func(header [3]byte, end [2]byte) bool {
		return header == [3]byte{1, 131, 0}
	})
	return string(payload), err
}

// History requests and returns history data. Not yet implemented.
func (s Samil) History() error {
	if *s.closed != nil {
		return *s.closed
	}
	panic("Not yet implemented")
}

// RemoteAddr returns the remote network address. The Addr returned is shared by
// all invocations of RemoteAddr, so do not modify it.
func (s Samil) RemoteAddr() net.Addr {
	return s.conn.RemoteAddr()
}

// LocalAddr returns the local network address. The Addr returned is shared by
// all invocations of LocalAddr, so do not modify it.
func (s Samil) LocalAddr() net.Addr {
	return s.conn.LocalAddr()
}

// Close closes the connection.
func (s Samil) Close() error {
	return s.conn.Close()
}
