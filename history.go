package samil

import (
	"time"
)

// History requests and returns history data in the time period provided.
// The time period is provided as two integers for the last digits of the start
// and end year. E.g. start=7 and end=10 are for 2007 and 2010.
/// IMPLEMENTATION NOT FINISHED.
func (s Samil) History(start, end int) error {
	if *s.closed != nil {
		return *s.closed
	}
	_, err := s.conn.Write(historyRequest(start, end))
	if err != nil {
		return err
	}
	/*
		payload, err := s.readData()
		if err != nil {
			return InverterData{}, err
		}
		return dataFrom(payload), nil
	*/
	time.Sleep(1 * time.Hour)
	panic("Not yet implemented")
}
