package samil

import (
	"encoding/binary"
	"io"
)

type message struct {
	header  [3]byte
	payload []byte
}

// Returns payload of the next message for which the condition function holds.
// The header of a message is passed to the condition function.
func (s *Samil) readFor(hold func([3]byte) bool) ([]byte, error) {
	var msg message
	var err error
	for {
		msg, err = s.read()
		if err != nil || hold(msg.header) {
			return msg.payload, err
		}
	}
}

// Returns next (buffered) message.
func (s *Samil) read() (message, error) {
	msg, ok := <-s.in
	if !ok {
		return message{}, s.closed
	}
	return msg, nil
}

// Read routine to be run as separate goroutine for processing and buffering messages.
// Use read() or readFor() to read buffered messages.
// Closes the connection at EOF and sets closed error flag.
func (s *Samil) readRoutine() {
	defer s.conn.Close()
	for {
		msg, err := s.readNext()
		if err != nil {
			s.closed = err
			close(s.in)
			return
		}
		s.in <- msg
	}
}

// Reads next incoming message, used in the read routine.
func (s *Samil) readNext() (msg message, err error) {
	start := make([]byte, 2)
	_, err = io.ReadFull(s.conn, start)
	if err != nil {
		return
	}
	if start[0] != 0x55 || start[1] != 0xaa {
		panic("Invalid message, not starting with 55 aa bytes")
	}
	_, err = io.ReadFull(s.conn, msg.header[:])
	if err != nil {
		return
	}
	sizeBytes := make([]byte, 2)
	_, err = io.ReadFull(s.conn, sizeBytes)
	if err != nil {
		return
	}
	size := int(binary.BigEndian.Uint16(sizeBytes))
	msg.payload = make([]byte, size)
	_, err = io.ReadFull(s.conn, msg.payload)
	if err != nil {
		return
	}
	chksumBytes := make([]byte, 2)
	_, err = io.ReadFull(s.conn, chksumBytes)
	if err != nil {
		return
	}
	chksum := int(binary.BigEndian.Uint16(chksumBytes))
	if chksum != checksum(start) + checksum(msg.header[:]) + checksum(sizeBytes) + checksum(msg.payload) {
		panic("Invalid message, incorrect checksum")
	}
	return
}
