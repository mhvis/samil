package samil

import (
	"strings"
)

// Device type values.
const (
	SinglePhaseInverter = 1
	
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
	str := string(payload)
	return &Model{
		DeviceType:           int(payload[0]),
		VARating:             strings.TrimSpace(str[1:7]),
		FirmwareVersion:      str[7:12],
		ModelName:            str[12:26],
		Manufacturer:         str[28:38],
		SerialNumber:         str[44:54],
		CommunicationVersion: str[60:65],
		OtherVersion:         str[65:70],
		General:              int(payload[70]),
	}
}
