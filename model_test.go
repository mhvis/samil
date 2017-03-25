package samil

import (
	"testing"
)

func TestModelInfoFromPayload(t *testing.T) {
	payload := []byte("1  4500V1.30River 4500TL-D\x00 SamilPower\x00ANKIEDW413B8080\x00\x00\x00\x00\x00\x00V1.30V1.302")
	expect := Model{
		DeviceType:           1,
		VARating:             "4500",
		FirmwareVersion:      "V1.30",
		ModelName:            "River 4500TL-D",
		Manufacturer:         "SamilPower",
		SerialNumber:         "DW413B8080",
		CommunicationVersion: "V1.30",
		OtherVersion:         "V1.30",
		General:              2,
	}
	model := modelFrom(payload)
	if expect != *model {
		t.Errorf("Incorrect model from payload, expected %v, got %v", expect, *model)
	}

}
