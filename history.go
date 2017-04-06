package samil

import (
	"encoding/binary"
	"fmt"
	"strings"
)

// History requests and returns history data in the time period provided.
// The time period is provided as two integers for the last digits of the start
// and end year. E.g. start=7 and end=10 are for 2007 and 2010.
//
// The history data is asynchronously returned via the channel parameter. This
// channel must be provided by the caller. The method will block while the data
// is being returned on the channel, and will only unblock after all data is
// received or an error occurred. In the case of an error, the error is
// returned, else the return value is nil. The provided channel will be closed
// after the data is received or an error occurred. The caller should ensure
// that the channel is not full for too long periods, otherwise incoming
// messages may get discarded.
func (s *Samil) History(start, end int, c chan HistoryDay) error {
	err := s.write(historyPacket(start, end))
	if err != nil {
		close(c)
		return err
	}
	store := make(historyStore)
	for {
		msg, err := s.read()
		if err != nil {
			close(c)
			return err
		}
		// Check for end packet
		if msg.header[0] == 6 && msg.header[1] == 0x81 {
			break
		}
		// Check for unknown packet
		if msg.header[0] != 6 || msg.header[1] != 0x61 {
			continue
		}
		d, cnt, seq, value := historyPayload(msg.payload)
		if store.has(d, seq) {
			close(c)
			return fmt.Errorf("received duplicate data part: %v %v %v %v %v",
				s, d, cnt, seq, value)
		}
		store.add(d, cnt, seq, value)
		if store.isFull(d) {
			c <- HistoryDay{d.Year, d.Month, d.Day, store.get(d)}
		}
	}
	// Optionally could check if incompleted days exist (in that case return err)
	close(c)
	return nil
}

func historyPayload(p []byte) (d date, cnt, seq int, value *string) {
	d.Year = int(p[0])
	d.Month = int(p[1])
	d.Day = int(p[2])
	cnt = int(binary.BigEndian.Uint32(p[7:11]))
	seq = int(binary.BigEndian.Uint32(p[11:15]))
	v := string(p[15:])
	value = &v
	return
}

type date struct {
	Year, Month, Day int
}

// Hash map to combine different day part packets into complete days.
type historyStore map[date][]*string

func (s historyStore) has(d date, seq int) bool {
	slice, present := s[d]
	return present && slice[seq] != nil
}

// Adds part of a day to the store
func (s historyStore) add(d date, cnt, seq int, value *string) {
	slice, present := s[d]
	if !present {
		slice = make([]*string, cnt)
		s[d] = slice
	}
	slice[seq] = value
}

func (s historyStore) isFull(d date) bool {
	for _, v := range s[d] {
		if v == nil {
			return false
		}
	}
	return true
}

// Gets the full (combined) day value
func (s historyStore) get(d date) string {
	ptrs := s[d]
	vals := make([]string, len(ptrs))
	for i, ptr := range ptrs {
		vals[i] = *ptr
	}
	return strings.Join(vals, "")
}

// HistoryDay stores the history data for a single day.
type HistoryDay struct {
	// Year number in 2 digits (e.g. 99 for 2099)
	Year int
	// Month number
	Month int
	// Day numer
	Day int
	// Value is a csv-encoded string of generation per hour values.
	// Please check the format yourself as it could be different.
	Value string
}
