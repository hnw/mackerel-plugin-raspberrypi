package mpraspberrypi

import (
	"reflect"
	"testing"
)

func TestParseTemperature(t *testing.T) {
	stub := `temp=48.0'C
`
	temperature, _ := parseTemperature(stub)
	expected := float64(48.0)
	if !reflect.DeepEqual(expected, temperature) {
		t.Errorf("Expected %#v, got %#v", expected, temperature)
	}
}

func TestParseClock(t *testing.T) {
	stub := `frequency(50)=200000000
`
	voltage, _ := parseClock(stub)
	expected := float64(200000000)
	if !reflect.DeepEqual(expected, voltage) {
		t.Errorf("Expected %#v, got %#v", expected, voltage)
	}
}

func TestParseVoltage(t *testing.T) {
	stub := `volt=1.2000V
`
	voltage, _ := parseVoltage(stub)
	expected := float64(1.2)
	if !reflect.DeepEqual(expected, voltage) {
		t.Errorf("Expected %#v, got %#v", expected, voltage)
	}
}

func TestParseThrottled(t *testing.T) {
	stub := `throttled=0x50005
`
	voltage, _ := parseThrottled(stub)
	expected := map[string]float64{
		"under_voltage":              1,
		"frequency_capped":           0,
		"throttled":                  1,
		"temperature_limit":          0,
		"under_voltage_occurred":     1,
		"frequency_capped_occurred":  0,
		"throttled_occurred":         1,
		"temperature_limit_occurred": 0,
	}
	if !reflect.DeepEqual(expected, voltage) {
		t.Errorf("Expected %#v, got %#v", expected, voltage)
	}
}
