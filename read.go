package samil

import (
	"encoding/binary"
	"io"
)

type message struct {
	header  [3]byte
	payload []byte
}

// Reads next messages until the condition holds, returns that message.
// The header of message is passed to the condition function.
func (s Samil) readFor(hold func([3]byte) bool) ([]byte, error) {
	for {
		msg, ok := <-s.read
		if !ok {
			return nil, *s.closed
		}
		if hold(msg.header) {
			return msg.payload, nil
		}
	}
}

// Reads continuously, closes the connection at EOF and sets closed error flag.
func (s *Samil) readRoutine() {
	defer s.conn.Close()
	for {
		msg, err := s.readNext()
		if err != nil {
			*s.closed = err
			close(s.read)
			return
		}
		s.read <- msg
	}
}

// Reads next incoming message.
func (s Samil) readNext() (msg message, err error) {
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
