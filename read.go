package samil

import (
	"encoding/binary"
)

// Reads next messages until header matches, returns that message.
func (s Samil) readFor(header [3]byte) (payload []byte, err error) {
	var head [3]byte
	for {
		head, payload, _, err = s.Read()
		if err != nil || head == header {
			return
		}
	}
}

// Reads next incoming message.
func (s Samil) Read() (header [3]byte, payload []byte, end [2]byte, err error) {
	start := make([]byte, 2)
	_, err = s.conn.Read(start)
	if err != nil {
		return
	}
	if start[0] != 85 || start[1] != 170 {
		panic("Invalid message, not starting with 85 170 bytes")
	}
	_, err = s.conn.Read(header[:])
	if err != nil {
		return
	}
	sizeBytes := make([]byte, 2)
	_, err = s.conn.Read(sizeBytes)
	if err != nil {
		return
	}
	size := int(binary.BigEndian.Uint16(sizeBytes))
	payload = make([]byte, size)
	_, err = s.conn.Read(payload)
	if err != nil {
		return
	}
	_, err = s.conn.Read(end[:])
	return
}
