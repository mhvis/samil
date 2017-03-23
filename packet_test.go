package samil

import (
	"testing"
	"bytes"
)

func TestChecksum(t *testing.T) {
	expected := 9
	actual := checksum([]byte{3, 0, 4, 2})
	if expected != actual {
		t.Errorf("checksum failed, expected %v, got %v", expected, actual)
	}
}

func TestForgePacket(t *testing.T) {
	expected := []byte{0x55, 0xaa, 1, 0, 1, 0, 1, 1, 1, 3}
	actual := forgePacket([3]byte{1, 0, 1}, []byte{1})
	if !bytes.Equal(expected, actual) {
		t.Errorf("expected %v, got %v", expected, actual)
	}
}
