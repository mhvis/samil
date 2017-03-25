package samil

import (
	"strconv"
	"strings"
)

// Device type values.
const (
	SinglePhaseInverter = 1
	ThreePhaseInverter  = 2
	SolarEnviMonitor    = 3
	RPhaseInverter      = 4
	SPhaseInverter      = 5
	TPhaseInverter      = 6
)

// Model stores model and version information.
type Model struct {
	// Device type, see constants for possible values
	DeviceType int
	// Volt-ampere rating, e.g. "4500"
	VARating string
	// Firmware version, e.g. "V1.30"
	FirmwareVersion string
	// Model name, e.g. "River 4500TL-D"
	ModelName string
	// Manufacturer, e.g. "SamilPower"
	Manufacturer string
	// Serial number, e.g. "DW413B8080"
	SerialNumber string
	// Communication version, e.g. "V1.30"
	CommunicationVersion string
	// Other version, I don't know what it means, for me it is "V1.30"
	OtherVersion string
	// General, I don't know what it means (maybe a version code), for me it is
	// 2. When your inverter returns something different than 2, there is a
	// chance that generation data is not correctly interpreted, notably the PV2
	// voltage and current.
	General int
}

// Model requests and returns model and version information.
func (s *Samil) Model() (*Model, error) {
	err := s.write(model)
	if err != nil {
		return nil, err
	}
	payload, err := s.readFor(func(header [3]byte) bool {
		return header[0] == 1 && header[1] == 0x83
	})
	if err != nil {
		return nil, err
	}
	return modelFrom(payload), nil
}

// Converts payload to Model struct.
func modelFrom(payload []byte) *Model {
	if len(payload) != 71 {
		panic("Invalid model payload length")
	}
	return &Model{
		DeviceType:           atoi(string(payload[0:1])),
		VARating:             stringFrom(payload[1:7]),
		FirmwareVersion:      stringFrom(payload[7:12]),
		ModelName:            stringFrom(payload[12:28]),
		Manufacturer:         stringFrom(payload[28:44]),
		SerialNumber:         stringFrom(payload[44:60]),
		CommunicationVersion: stringFrom(payload[60:65]),
		OtherVersion:         stringFrom(payload[65:70]),
		General:              atoi(string(payload[70:71])),
	}
}

func atoi(s string) int {
	i, err := strconv.Atoi(s)
	if err != nil {
		panic("Invalid byte in model payload: " + err.Error())
	}
	return i
}

// Converts and trims cstring
func stringFrom(b []byte) string {
	return strings.TrimSpace(string(b[:clen(b)]))
}

// Length of a cstring
func clen(n []byte) int {
	for i := 0; i < len(n); i++ {
		if n[i] == 0 {
			return i
		}
	}
	return len(n)
}
