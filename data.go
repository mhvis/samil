package samil

import (
	"encoding/binary"
)

// Possible operating modes returned by the inverter.
const (
	Wait           = 0
	Normal         = 1
	Fault          = 2
	PermanentFault = 3
	Check          = 4
	PVPowerOff     = 5
)

// Data stores generation and operational data from the inverter.
type Data struct {
	// Internal temperature in decicelcius (375 = 37.5 degrees Celsius)
	InternalTemperature int
	// PV1 voltage in decivolts (2975 = 297.5 V)
	PV1Voltage int
	// PV2 voltage in decivolts
	PV2Voltage int
	// PV1 current in deciampere
	PV1Current int
	// PV2 current in deciampere
	PV2Current int
	// Total operation time in hours
	OperationTime int
	// Operating mode, see constants for possible values
	OperatingMode int
	// Energy produced today in decawatt hour (474 = 4.74 kWh)
	EnergyToday int
	// PV1 input power in watt
	PV1Power int
	// PV2 input power in watt
	PV2Power int
	// Single phase grid current in deciampere
	GridCurrent int
	// Grid voltage in decivolts
	GridVoltage int
	// Grid frequency in centihertz (4998 = 49.98 Hz)
	GridFrequency int
	// Output power in watt
	OutputPower int
	// Total energy produced in hectowatt hour (114649 = 11464.9 kWh)
	EnergyTotal int
}

// Data requests current data values from the inverter and returns them in the
// Data struct.
func (s *Samil) Data() (*Data, error) {
	err := s.write(data)
	if err != nil {
		return nil, err
	}
	payload, err := s.readData()
	if err != nil {
		return nil, err
	}
	return dataFrom(payload), nil
}

// Returns the next data message from the socket.
func (s *Samil) readData() ([]byte, error) {
	return s.readFor(func(header [3]byte) bool {
		return header[0] == 1 && header[1] == 0x82
	})
}

// Payload to Data struct.
func dataFrom(payload []byte) *Data {
	if len(payload) != 54 {
		panic("Unexpected data length")
	}
	return &Data{
		InternalTemperature: intFrom(payload[0:2], true),
		PV1Voltage:          intFrom(payload[2:4], false),
		PV2Voltage:          intFrom(payload[4:6], false),
		PV1Current:          intFrom(payload[6:8], false),
		PV2Current:          intFrom(payload[8:10], false),
		OperationTime:       intFrom(payload[10:14], false),
		OperatingMode:       intFrom(payload[14:16], false),
		EnergyToday:         intFrom(payload[16:18], false),
		PV1Power:            intFrom(payload[38:40], false),
		PV2Power:            intFrom(payload[40:42], false),
		GridCurrent:         intFrom(payload[42:44], false),
		GridVoltage:         intFrom(payload[44:46], false),
		GridFrequency:       intFrom(payload[46:48], false),
		OutputPower:         intFrom(payload[48:50], false),
		EnergyTotal:         intFrom(payload[50:54], false),
	}
}

func intFrom(b []byte, signed bool) int {
	switch len(b) {
	case 2:
		i := binary.BigEndian.Uint16(b)
		if signed {
			return int(int16(i))
		}
		return int(i)
	case 4:
		i := binary.BigEndian.Uint32(b)
		if signed {
			return int(int32(i))
		}
		return int(i)
	default:
		panic("Invalid integer byte sequence encoding, incorrect length")
	}
}
