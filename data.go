package samil

import (
	"encoding/binary"
)

const (
	// Operating modes.
	ModeWait       = 0
	ModeNormal     = 1
	ModePVPowerOff = 5
)

// InverterData stores generation data from the inverter.
type InverterData struct {
	// Internal temperature in decicelcius (375 = 37.5 degrees Celcius)
	InternalTemperature int
	// PV1 voltage in decivolts (2975 = 297.5 V)
	PV1Voltage int
	// PV2 voltage in decivolts
	PV2Voltage int
	// PV1 current in deciampère
	PV1Current int
	// PV2 current in deciampère
	PV2Current int
	// Total operation time in hours
	OperationTime int
	// Operating mode, 0=wait, 1=normal, 5=PV power off (?)
	OperatingMode int
	// Energy produced today in decawatt hour (474 = 4.74 kWh)
	EnergyToday int
	// PV1 input power in watt
	PV1Power int
	// PV2 input power in watt
	PV2Power int
	// Single phase grid current in deciampère
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
// InverterData struct.
func (s Samil) Data() (InverterData, error) {
	if *s.closed != nil {
		return InverterData{}, *s.closed
	}
	_, err := s.conn.Write(data)
	if err != nil {
		return InverterData{}, err
	}
	payload, err := s.readData()
	if err != nil {
		return InverterData{}, err
	}
	return dataFrom(payload), nil
}

// Returns the next data message from the socket.
func (s Samil) readData() ([]byte, error) {
	return s.readFor(func(header [3]byte, end [2]byte) bool {
		return header[0] == 1 && header[1] == 130
	})
}

// Payload to InverterData struct.
func dataFrom(payload []byte) InverterData {
	if len(payload) != 54 {
		panic("Unexpected data length")
	}
	return InverterData{
		InternalTemperature: intFrom(payload[0:2]),
		PV1Voltage:          intFrom(payload[2:4]),
		PV2Voltage:          intFrom(payload[4:6]),
		PV1Current:          intFrom(payload[6:8]),
		PV2Current:          intFrom(payload[8:10]),
		OperationTime:       intFrom(payload[12:14]),
		OperatingMode:       intFrom(payload[14:16]),
		EnergyToday:         intFrom(payload[16:18]),
		PV1Power:            intFrom(payload[38:40]),
		PV2Power:            intFrom(payload[40:42]),
		GridCurrent:         intFrom(payload[42:44]),
		GridVoltage:         intFrom(payload[44:46]),
		GridFrequency:       intFrom(payload[46:48]),
		OutputPower:         intFrom(payload[48:50]),
		EnergyTotal:         intFrom(payload[50:54]),
	}
}

func intFrom(b []byte) int {
	switch len(b) {
	case 2:
		return int(binary.BigEndian.Uint16(b))
	case 4:
		return int(binary.BigEndian.Uint32(b))
	default:
		panic("Invalid integer byte sequence encoding, incorrect length")
	}
}
