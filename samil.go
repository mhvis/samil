// Package samil provides an API for quering networked Samil Power inverters.
// Supported inverters are those that work with SolarPower Browser.
// These are: SolarRiver TD, SolarRiver TL-D and SolarLake TL inverters.
package samil

import (
	"net"
)

type Samil struct {
	conn *net.TCPConn

	//readErr chan error
	//readMsg chan []byte
}

// NewConnection searches for an inverter in the network and returns the
// connection if one is found.
//
// Inverters that are connected already to a client will not initiate a new
// connection. Therefore calling this function multiple times while leaving the
// connections open will connect to different inverters.
func NewConnection() (Samil, error) {
	conn, err := connect()
	if err != nil {
		return Samil{}, err
	}
	return Samil{conn}, nil
}

// ModelInfo.
func (s Samil) ModelInfo() error {
	_, err := s.conn.Write(modelInfo)
	return err
}

// Status.
func (s Samil) Status() error {
	_, err := s.conn.Write(data)
	return err
}

func (s Samil) History() error {
	panic("")
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
