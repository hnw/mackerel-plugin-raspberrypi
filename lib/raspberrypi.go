package mpraspberrypi

import (
	"bytes"
	"flag"
	"fmt"
	"os/exec"
	"regexp"
	"strconv"

	mp "github.com/mackerelio/go-mackerel-plugin"
)

// RaspberrypiPlugin mackerel plugin
type RaspberrypiPlugin struct {
	Prefix string
}

var temperaturePattern = regexp.MustCompile(
	`^temp=(\d+.\d+)'C\n$`,
)

var clockPattern = regexp.MustCompile(
	`^frequency\(\d+\)=(\d+)\n$`,
)

var voltagePattern = regexp.MustCompile(
	`^volt=(\d+.\d+)V\n$`,
)

var throttledPattern = regexp.MustCompile(
	`^throttled=0x([\da-fA-F]+)\n$`,
)

// MetricKeyPrefix interface for PluginWithPrefix
func (r RaspberrypiPlugin) MetricKeyPrefix() string {
	if r.Prefix == "" {
		r.Prefix = "raspberrypi"
	}
	return r.Prefix
}

// GraphDefinition interface for mackerelplugin
func (r RaspberrypiPlugin) GraphDefinition() map[string]mp.Graphs {
	return map[string]mp.Graphs{
		"temperature": {
			Label: "Temperature ['C]",
			Unit:  mp.UnitFloat,
			Metrics: []mp.Metrics{
				{Name: "soc_temp", Label: "SoC Temperature"},
			},
		},
		"clock": {
			Label: "Clock frequency [MHz]",
			Unit:  mp.UnitFloat,
			Metrics: []mp.Metrics{
				{Name: "arm", Label: "ARM cores", Scale: 0.000001},
				{Name: "core", Label: "VC4 scaler cores", Scale: 0.000001},
				{Name: "H264", Label: "H264 block", Scale: 0.000001},
				{Name: "isp", Label: "Image Signal Processor", Scale: 0.000001},
				{Name: "v3d", Label: "3D block", Scale: 0.000001},
				{Name: "uart", Label: "UART", Scale: 0.000001},
				{Name: "pwm", Label: "PWM block", Scale: 0.000001},
				{Name: "emmc", Label: "SD card interface", Scale: 0.000001},
				{Name: "pixel", Label: "Pixel valve", Scale: 0.000001},
				{Name: "vec", Label: "Analogue video encoder", Scale: 0.000001},
				{Name: "hdmi", Label: "HDMI", Scale: 0.000001},
				{Name: "dpi", Label: "Display Peripheral Interface", Scale: 0.000001},
			},
		},
		"voltage": {
			Label: "Voltage [V]",
			Unit:  mp.UnitFloat,
			Metrics: []mp.Metrics{
				{Name: "core_volts", Label: "VC4 core voltage"},
				{Name: "sdram_c", Label: "SDRAM Core Voltage"},
				{Name: "sdram_i", Label: "SDRAM I/O voltage"},
				{Name: "sdram_p", Label: "SDRAM Phy Voltage"},
			},
		},
		"state": {
			Label: "Throttled state",
			Unit:  mp.UnitInteger,
			Metrics: []mp.Metrics{
				{Name: "under_voltage", Label: "Under-voltage detected"},
				{Name: "frequency_capped", Label: "Arm frequency capped"},
				{Name: "throttled", Label: "Throttled"},
				{Name: "temperature_limit", Label: "Soft temperature limit active"},
			},
		},
	}
}

// FetchMetrics interface for mackerelplugin
func (r RaspberrypiPlugin) FetchMetrics() (map[string]float64, error) {
	funcs := []func() (metrics map[string]float64, err error){
		getTemperature,
		getClock,
		getVoltage,
		getThrottled,
	}
	metrics := make(map[string]float64)
	for _, f := range funcs {
		tmpMetrics, err := f()
		if err != nil {
			return nil, err
		}
		for k, v := range tmpMetrics {
			metrics[k] = v
		}
	}
	return metrics, nil
}

func getTemperature() (metrics map[string]float64, err error) {
	cmd := exec.Command("vcgencmd", "measure_temp")
	var out bytes.Buffer
	cmd.Stdout = &out
	err = cmd.Run()
	if err != nil {
		return
	}
	metrics = make(map[string]float64)
	metrics["soc_temp"], err = parseTemperature(out.String())
	if err != nil {
		return nil, err
	}
	return
}

func parseTemperature(out string) (float64, error) {
	if matches := temperaturePattern.FindStringSubmatch(out); matches != nil {
		return strconv.ParseFloat(matches[1], 64)
	}
	return 0, fmt.Errorf("Failed to parse temperature: %s", out)
}

func getClock() (metrics map[string]float64, err error) {
	metrics = make(map[string]float64)
	devices := []string{
		"arm", "core", "H264", "isp", "v3d", "uart", "pwm", "emmc", "pixel",
		"vec", "hdmi", "dpi",
	}
	for _, device := range devices {
		cmd := exec.Command("vcgencmd", "measure_clock", device)
		var out bytes.Buffer
		cmd.Stdout = &out
		err := cmd.Run()
		if err != nil {
			return nil, err
		}
		metrics[device], err = parseClock(out.String())
		if err != nil {
			return nil, err
		}
	}
	return
}

func parseClock(out string) (float64, error) {
	if matches := clockPattern.FindStringSubmatch(out); matches != nil {
		return strconv.ParseFloat(matches[1], 64)
	}
	return 0, fmt.Errorf("Failed to parse temperature: %s", out)
}

func getVoltage() (metrics map[string]float64, err error) {
	metrics = make(map[string]float64)
	blocks := []string{"core", "sdram_c", "sdram_i", "sdram_p"}
	for _, block := range blocks {
		cmd := exec.Command("vcgencmd", "measure_volts", block)
		var out bytes.Buffer
		cmd.Stdout = &out
		err = cmd.Run()
		if err != nil {
			return nil, err
		}
		if block == "core" {
			block = "core_volts"
		}
		metrics[block], err = parseVoltage(out.String())
		if err != nil {
			return nil, err
		}
	}
	return
}

func parseVoltage(out string) (float64, error) {
	if matches := voltagePattern.FindStringSubmatch(out); matches != nil {
		return strconv.ParseFloat(matches[1], 64)
	}
	return 0, fmt.Errorf("Failed to parse temperature: %s", out)
}

func getThrottled() (metrics map[string]float64, err error) {
	cmd := exec.Command("vcgencmd", "get_throttled")
	var out bytes.Buffer
	cmd.Stdout = &out
	err = cmd.Run()
	if err != nil {
		return nil, err
	}
	return parseThrottled(out.String())
}

func parseThrottled(out string) (metrics map[string]float64, err error) {
	if matches := throttledPattern.FindStringSubmatch(out); matches != nil {
		var rawValue int64
		rawValue, err = strconv.ParseInt(matches[1], 16, 32)
		if err != nil {
			return
		}
		metrics = make(map[string]float64)
		metrics["under_voltage"] = float64(rawValue & 0x1 >> 0)
		metrics["frequency_capped"] = float64(rawValue & 0x2 >> 1)
		metrics["throttled"] = float64(rawValue & 0x4 >> 2)
		metrics["temperature_limit"] = float64(rawValue & 0x8 >> 3)
		return
	}
	return nil, fmt.Errorf("Failed to parse throttled: %s", out)
}

// Do the plugin
func Do() {
	optPrefix := flag.String("metric-key-prefix", "raspberrypi", "Metric key prefix")
	optTempfile := flag.String("tempfile", "", "Temp file name")
	flag.Parse()

	r := RaspberrypiPlugin{
		Prefix: *optPrefix,
	}
	plugin := mp.NewMackerelPlugin(r)
	plugin.Tempfile = *optTempfile
	plugin.Run()
}
