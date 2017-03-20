package samil

import (
	"testing"
)

func TestDataFromPayload(t *testing.T) {
	payload := []byte{1, 119, 11, 163, 11, 243, 0, 21, 0, 20, 0, 0, 40, 64, 0,
		1, 1, 218, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
		2, 138, 2, 114, 0, 55, 9, 20, 19, 134, 4, 245, 0, 1, 177, 204}
	expect := InverterData{
		InternalTemperature: 375,
		PV1Voltage:          2979,
		PV2Voltage:          3059,
		PV1Current:          21,
		PV2Current:          20,
		OperationTime:       10304,
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
	data := dataFrom(payload)
	if expect != data {
		t.Errorf("Incorrect data from payload, expected %v, got %v", expect, data)
	}
}
