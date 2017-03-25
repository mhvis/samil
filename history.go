package samil

// History requests and returns history data in the time period provided.
// The time period is provided as two integers for the last digits of the start
// and end year. E.g. start=7 and end=10 are for 2007 and 2010.
// IMPLEMENTATION NOT FINISHED.
func (s *Samil) History(start, end int) error {
	err := s.write(historyPacket(start, end))
	if err != nil {
		return err
	}
	// Read until EOF
	_, err = s.readFor(func(header [3]byte) bool {
		return false
	})
	return err
}
