package samil

import (
	"net"
	"testing"
)

func TestHistoryStoreAdd(t *testing.T) {
	s := make(historyStore)
	d := date{12, 3, 4}
	val := "test"
	s.add(d, 2, 1, &val)
	slice := s[d]
	if len(slice) != 2 {
		t.Errorf("wrong slice length, expected 2, got %v", len(slice))
	}
	if slice[0] != nil {
		t.Errorf("wrong slice element value, expected nil, got %v", *slice[0])
	}
	if slice[1] != &val {
		t.Errorf("wrong slice element value, expected \"test\", got %v", *slice[1])
	}
}

func TestHistoryStoreHas(t *testing.T) {
	s := make(historyStore)
	d := date{12, 3, 4}
	val := "test"
	if s.has(d, 0) {
		t.Errorf("expected initially false, got true")
	}
	s.add(d, 1, 0, &val)
	if !s.has(d, 0) {
		t.Errorf("expected true, got false")
	}
}

func TestHistoryStoreIsFull(t *testing.T) {
	s := make(historyStore)
	d := date{12, 3, 4}
	val := "test"
	s.add(d, 3, 0, &val)
	if s.isFull(d) {
		t.Errorf("expected initially false, got true")
	}
	s.add(d, 3, 1, &val)
	if s.isFull(d) {
		t.Errorf("expected false after adding 2 of 3, but got true")
	}
	s.add(d, 3, 2, &val)
	if !s.isFull(d) {
		t.Errorf("expected true after adding all elements, got false")
	}
}

func TestHistoryStoreGet(t *testing.T) {
	s := make(historyStore)
	d := date{12, 3, 4}
	val := "test"
	s.add(d, 2, 0, &val)
	s.add(d, 2, 1, &val)
	expect := "testtest"
	actual := s.get(d)
	if expect != actual {
		t.Errorf("expected %v, got %v", expect, actual)
	}
}

func TestHistoryPayload(t *testing.T) {
	payload := []byte{11, 3, 12,
		0, 0, 1, 0x38,
		1, 1, 1, 1,
		1, 1, 1, 2,
		'a', 'n', 'k', 'i', 'e',
	}
	d, cnt, seq, val := historyPayload(payload)
	exD := date{11, 3, 12}
	exCnt := 16843009
	exSeq := 16843010
	exVal := "ankie"
	if exD != d || exCnt != cnt || exSeq != seq || exVal != *val {
		t.Errorf("expected %v %v %v %v, got %v %v %v %v",
			exD, exCnt, exSeq, exVal, d, cnt, seq, val)
	}
	// Empty data
	payload = []byte{11, 3, 12,
		0, 0, 1, 0x38,
		0, 0, 0, 1,
		0, 0, 0, 0,
	}
	d, cnt, seq, val = historyPayload(payload)
	exD = date{11, 3, 12}
	exCnt = 1
	exSeq = 0
	exVal = ""
	if exD != d || exCnt != cnt || exSeq != seq || exVal != *val {
		t.Errorf("expected %v %v %v %v, got %v %v %v %v",
			exD, exCnt, exSeq, exVal, d, cnt, seq, val)
	}
}

func readingSamilClient() (s *Samil, c net.Conn) {
	c, c2 := net.Pipe()
	s = &Samil{c2, make(chan message, 5), nil}
	go s.readRoutine()
	return
}

func TestHistory(t *testing.T) {
	samil, conn := readingSamilClient()
	// Discard incoming messages
	go func() {
		for {
			b := make([]byte, 1)
			_, err := conn.Read(b)
			t.Logf("received byte: %v/%v", b, string(b))
			if err != nil {
				return
			}
		}
	}()

	// Test normal
	ch := make(chan HistoryDay, 2)
	go func() {
		_, err := conn.Write(forgePacket([3]byte{6, 0x61, 0}, []byte{
			10, 1, 1, // Date
			0, 0, 1, 0x38, // ?
			0, 0, 0, 2, // Count
			0, 0, 0, 1, // Sequence number
			'k', 'i', 'e', // Value text
		}))
		if err != nil {
			panic(err)
		}
		conn.Write(forgePacket([3]byte{6, 0x61, 0}, []byte{
			10, 1, 1, // Date
			0, 0, 1, 0x38, // ?
			0, 0, 0, 2, // Count
			0, 0, 0, 0, // Sequence number
			'a', 'n', // Value text
		}))
		conn.Write(forgePacket([3]byte{6, 0x81, 0}, []byte{}))
	}()
	err := samil.History(1, 1, ch)
	if err != nil {
		t.Errorf("got unexpected error: %v", err)
	}
	expect := HistoryDay{10, 1, 1, "ankie"}
	actual := <-ch
	if expect != actual {
		t.Errorf("wrong history day, expected %v, got %v", expect, actual)
	}
	if _, ok := <-ch; ok {
		t.Errorf("channel is not closed after function returned")
	}

	// Test duplicate packet
	ch = make(chan HistoryDay, 2)
	go func() {
		conn.Write(forgePacket([3]byte{6, 0x61, 0}, []byte{
			10, 1, 1, // Date
			0, 0, 1, 0x38, // ?
			0, 0, 0, 2, // Count
			0, 0, 0, 1, // Sequence number
			'k', 'i', 'e', // Value text
		}))
		conn.Write(forgePacket([3]byte{6, 0x61, 0}, []byte{
			10, 1, 1, // Date
			0, 0, 1, 0x38, // ?
			0, 0, 0, 2, // Count
			0, 0, 0, 1, // Sequence number
			'a', 'n', // Value text
		}))
		conn.Write(forgePacket([3]byte{6, 0x81, 0}, []byte{}))
	}()
	err = samil.History(1, 1, ch)
	if err == nil {
		t.Errorf("got no error while should have one")
	}
	if _, ok := <-ch; ok {
		t.Errorf("channel is not closed after function returned")
	}
	samil.Close()
	conn.Close()
}
