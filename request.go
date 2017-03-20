package samil

import (
	"encoding/binary"
)

var advertisement = forgeRequest([3]byte{0, 64, 2}, []byte("I AM SERVER"), [2]byte{4, 58})

var modelInfo = forgeRequest([3]byte{1, 3, 2}, nil, [2]byte{1, 5})

var data = forgeRequest([3]byte{1, 2, 2}, nil, [2]byte{1, 4})

// Start and end are both the last two digits of a year number.
// I.e. 05, 07 means 2005 and 2007.
func historyRequest(start, end int) []byte {
	payload := []byte{byte(start), byte(end)}
	return forgeRequest([3]byte{6, 1, 2}, payload, [2]byte{1, 42})
}

func forgeRequest(header [3]byte, payload []byte, end [2]byte) []byte {
	request := make([]byte, 2+3+2+len(payload)+2)
	copy(request[0:2], []byte{85, 170})
	copy(request[2:5], header[:])
	binary.BigEndian.PutUint16(request[5:7], uint16(len(payload)))
	copy(request[7:7+len(payload)], payload)
	copy(request[7+len(payload):], end[:])
	return request
}
