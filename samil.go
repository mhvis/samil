// Package samil provides an API for quering networked Samil Power inverters.
// Supported inverters are those that work with SolarPower Browser V3 software.
// These are: SolarRiver TD, SolarRiver TL-D and SolarLake TL inverters.
// It is only tested and confirmed for SolarRiver 4500TL-D.
package samil

import (
	"net"
)

// Samil maintains a connection to an inverter.
//
// Connections are usually closed by the inverter after 20 seconds of
// inactivity. When a connection is closed, subsequent API calls will return the
// error EOF.
type Samil struct {
	conn   net.Conn
	in     chan message // Buffer for incoming messages
	closed error
}

// NewConnection searches for an inverter in the network and returns the
// connection if one is found.
//
// Inverters that are already connected to a client will not initiate a new
// connection. Therefore calling this function multiple times while leaving the
// connections open will connect to different inverters.
//
// The search will return with an i/o timeout error when no inverter is found
// after a minute.
func NewConnection() (*Samil, error) {
	conn, err := connect()
	if err != nil {
		return nil, err
	}
	s := &Samil{conn, make(chan message, 5), nil}
	go s.readRoutine()
	return s, nil
}

// Writes only after checking if we are still connected.
func (s *Samil) write(b []byte) error {
	if s.closed != nil {
		return s.closed
	}
	_, err := s.conn.Write(b)
	return err
}

// RemoteAddr returns the remote network address. The Addr returned is shared by
// all invocations of RemoteAddr, so do not modify it.
func (s *Samil) RemoteAddr() net.Addr {
	return s.conn.RemoteAddr()
}

// LocalAddr returns the local network address. The Addr returned is shared by
// all invocations of LocalAddr, so do not modify it.
func (s *Samil) LocalAddr() net.Addr {
	return s.conn.LocalAddr()
}

// Close closes the connection.
func (s *Samil) Close() error {
	return s.conn.Close()
}
