package samil

import (
	"testing"
)

func TestDataFromPayload(t *testing.T) {
	payload := []byte{0xff, 0xf7,
		11, 163,
		0x80, 243,
		0, 21,
		0, 20,
		0, 1, 1, 64,
		0, 1,
		1, 218,
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
		2, 138,
		2, 114,
		0, 55,
		9, 20,
		19, 134,
		4, 245,
		0, 1, 177, 204}
	expect := Data{
		InternalTemperature: -9,
		PV1Voltage:          2979,
		PV2Voltage:          33011,
		PV1Current:          21,
		PV2Current:          20,
		OperationTime:       65856,
		OperatingMode:       1,
		EnergyToday:         474,
		PV1Power:            650,
		PV2Power:            626,
		GridCurrent:         55,
		GridVoltage:         2324,
		GridFrequency:       4998,
		OutputPower:         1269,
		EnergyTotal:         111052,
	}
	data, err := dataFrom(payload)
	if err != nil {
		t.Errorf("unexpected error response: %v", err)
	}
	if expect != data {
		t.Errorf("incorrect data from payload, expected %v, got %v", expect, data)
	}
}

func TestIntFrom(t *testing.T) {
	// 16 bit unsigned
	assertInt(t, 65535, intFrom([]byte{0xff, 0xff}, false))
	// 16 bit signed
	assertInt(t, -1, intFrom([]byte{0xff, 0xff}, true))
	// 32 bit unsigned
	assertInt(t, 4294967295, intFrom([]byte{0xff, 0xff, 0xff, 0xff}, false))
	// 32 bit signed
	assertInt(t, -1, intFrom([]byte{0xff, 0xff, 0xff, 0xff}, true))
}

func assertInt(t *testing.T, expected, actual int) {
	if expected != actual {
		t.Errorf("incorrect int, expected %v, got %v", expected, actual)
	}
}
