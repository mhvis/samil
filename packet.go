package samil

import (
	"encoding/binary"
)

var advertisement = forgePacket([3]byte{0, 64, 2}, []byte("I AM SERVER"))

var model = forgePacket([3]byte{1, 3, 2}, nil)

var data = forgePacket([3]byte{1, 2, 2}, nil)

// Start and end are both the last two digits of a year number.
// I.e. 05, 07 means 2005 and 2007.
func historyPacket(start, end int) []byte {
	payload := []byte{byte(start), byte(end)}
	return forgePacket([3]byte{6, 1, 2}, payload)
}

func forgePacket(header [3]byte, payload []byte) []byte {
	request := make([]byte, 2+3+2+len(payload)+2)
	copy(request[0:2], []byte{0x55, 0xaa})
	copy(request[2:5], header[:])
	binary.BigEndian.PutUint16(request[5:7], uint16(len(payload)))
	copy(request[7:7+len(payload)], payload)
	checksum := uint16(checksum(request[:]))
	binary.BigEndian.PutUint16(request[7+len(payload):], checksum)
	return request
}

func checksum(packet []byte) (sum int) {
	for _, b := range packet {
		sum += int(b)
	}
	return
}
