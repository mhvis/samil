package samil

import (
	"strings"
)

// Device type values (known by me).
const (
	SinglePhaseInverter = 1
)

// ModelInfo stores model and version information.
type ModelInfo struct {
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
	// General, I don't know what it means, for me it is 2
	General int
}

// ModelInfo requests and returns model and version information.
// This is only confirmed to work for SolarRiver 4500TL-D.
func (s Samil) ModelInfo() (ModelInfo, error) {
	if *s.closed != nil {
		return ModelInfo{}, *s.closed
	}
	_, err := s.conn.Write(modelInfo)
	if err != nil {
		return ModelInfo{}, err
	}
	payload, err := s.readFor(func(header [3]byte) bool {
		return header == [3]byte{1, 131, 0}
	})
	if err != nil {
		return ModelInfo{}, err
	}
	info := string(payload)
	return ModelInfo{
		DeviceType:           int(payload[0]),
		VARating:             strings.TrimSpace(info[1:7]),
		FirmwareVersion:      info[7:12],
		ModelName:            info[12:26],
		Manufacturer:         info[28:38],
		SerialNumber:         info[44:54],
		CommunicationVersion: info[60:65],
		OtherVersion:         info[65:70],
		General:              int(payload[70]),
	}, nil
}
