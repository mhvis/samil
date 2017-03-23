package samil

import (
	"bytes"
	"net"
	"testing"
)

func TestReadNext(t *testing.T) {
	s, c := samilClient()
	go func() {
		c.Write([]byte{85, 170, 1})
		c.Write([]byte{128, 0, 0})
		c.Write([]byte{27})
		c.Write([]byte{0, 1, 2, 4, 5, 9, 10, 12, 17})
		c.Write([]byte{23, 24, 27, 28, 29, 30, 31, 32, 33, 34, 39, 40, 49, 50, 51, 52})
		c.Write([]byte{53, 54, 4})
		c.Write([]byte{126})
	}()
	msg, err := s.readNext()
	if err != nil {
		t.Error(err)
	}
	headerExpect := [3]byte{1, 128, 0}
	if headerExpect != msg.header {
		t.Errorf("Incorrect header, expected %v, got %v", headerExpect, msg.header)
	}
	payloadExpect := []byte{0, 1, 2, 4, 5, 9, 10, 12, 17, 23, 24, 27, 28, 29,
		30, 31, 32, 33, 34, 39, 40, 49, 50, 51, 52, 53, 54}
	if !bytes.Equal(payloadExpect, msg.payload) {
		t.Errorf("Incorrect payload, expected %v, got %v", payloadExpect, msg.payload)
	}
}

func samilClient() (s Samil, c net.Conn) {
	c, c2 := net.Pipe()
	s.conn = c2
	return
}
