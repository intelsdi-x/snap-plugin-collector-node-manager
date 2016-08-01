// +build linux

/*
http://www.apache.org/licenses/LICENSE-2.0.txt


Copyright 2015-2016 Intel Corporation

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package ipmi

import (
	"fmt"
)

// GenericValidator performs basic response validation. Checks response code ensures response
// has non-zero length.
type GenericValidator struct {
}

// ParserCUPS extracts data from CUPS specific response format.
// Data contains info about cpu utilization and memory & io bandwidth.
type ParserCUPS struct {
	*GenericValidator
}

// Instance of ParserCUPS
var FormatCUPS = &ParserCUPS{}

// GetMetrics method returns metric for CUPS parser: "cpu_cstate", "memory_bandwith", "io_bandwith"
func (p *ParserCUPS) GetMetrics() []string {
	return []string{"cpu_bandwith", "memory_bandwith", "io_bandwith"}
}

// Parse method returns data in human readable format
func (p *ParserCUPS) Parse(response IpmiResponse) map[string]uint64 {
	m := map[string]uint64{}
	// Parsing is based on command Get CUPS Data (65h). Bytes 5:6 contains CPU CUPS dynamic load factor
	// Bytes 7:8 contains memory CUPS dynamic load factor
	// Bytes 9:10 contains IO CUPS dynamic load factor
	var names = map[string]uint{
		"cpu_bandwith":    4,
		"memory_bandwith": 6,
		"io_bandwith":     8,
	}
	for metricName, startIndex := range names {
		m[metricName] = uint64(response.Data[startIndex]) + uint64(response.Data[startIndex+1])*256
	}
	return m
}

// ParserCUPS extracts data from CUPS specific response format.
// Data contains info about cpu utilization and memory & io bandwidth.
type ParserCUPSUtilization struct {
	*GenericValidator
}

// ParserCUPSIndex extracts CUPS Index from Node Manager
type ParserCUPSIndex struct {
	*GenericValidator
}

// Instance of ParserCUPS
var FormatCUPSIndex = &ParserCUPSIndex{}

// GetMetrics method returns metric for CUPS parser: "index"
func (p *ParserCUPSIndex) GetMetrics() []string {
	return []string{"index"}
}

// Parse method returns data in human readable format
func (p *ParserCUPSIndex) Parse(response IpmiResponse) map[string]uint64 {
	m := map[string]uint64{}
	// Parsing is based on command Get CUPS Data (65h). Bytes 5:6 contains CPU CUPS Index
	m["index"] = uint64(response.Data[4]) + uint64(response.Data[5])*256
	return m
}

// ParserNodeManager extracts data from Node manager response format.
// Data contains current, min, max and average value.
type ParserNodeManager struct {
	*GenericValidator
}

// Instance of ParserNodeManager
var FormatNodeManager = &ParserNodeManager{}

// GetMetrics method returns metric for CUPS parser: "current_value", "min", "max", "avg"
func (p *ParserNodeManager) GetMetrics() []string {
	return []string{"value", "min", "max", "avg"}
}

// Parse method returns data in human readable format
func (p *ParserNodeManager) Parse(response IpmiResponse) map[string]uint64 {
	m := map[string]uint64{}
	// Parsing is based on command Get Node Manager Statistics (C8h). Bytes 5:6 contains current value
	// Bytes 7:8 contains minimum value
	// Bytes 9:10 contains maximum value
	// Bytes 11:12 contains average value
	var names = map[string]uint{
		"value": 4,
		"min":   6,
		"max":   8,
		"avg":   10,
	}
	for metricName, startIndex := range names {
		m[metricName] = uint64(response.Data[startIndex]) + uint64(response.Data[startIndex+1])*256

	}
	return m
}

// ParserTemp extracts temperature data.
// Data contains info about temperatures for first 4 cpus
// and 64 dimms.
type ParserTemp struct {
	*GenericValidator
}

// Instance of ParserTempMargin.
var FormatTemp = &ParserTemp{}

// GetMetrics method returns metric for temperature parser: temperature of each cpu (up to 4),
// temperature of each dimm (up to 64)
func (p *ParserTemp) GetMetrics() []string {
	a := []string{"cpu/*", "dimm/*"}
	return a
}

// Parse method returns data in human readable format
func (p *ParserTemp) Parse(response IpmiResponse) map[string]uint64 {
	m := map[string]uint64{}
	// Parsing is based on Get CPU and Memory Temperature (4Bh). Bytes 5:8 contains temperatures of each socket (up to 4)
	m["cpu/0/value"] = uint64(response.Data[4])
	m["cpu/1/value"] = uint64(response.Data[5])
	m["cpu/2/value"] = uint64(response.Data[6])
	m["cpu/3/value"] = uint64(response.Data[7])
	// Bytes 9:72 contains temperatures of each dimm (up to 64)
	for i := 8; i < len(response.Data); i++ {
		a := fmt.Sprintf("dimm/%d/value", i-8)
		m[a] = uint64(response.Data[i])
	}

	return m
}

// ParserPECI extracts temperature margin datas from PECI response.
// Main metric value is TJ max.
// margin_offset current value of margin offset, which is value
// of TJ max reduction.
type ParserPECI struct {
	*GenericValidator
}

// Instance of ParserPECI.
var FormatPECI = &ParserPECI{}

// GetMetrics method returns metrics for PECI parser: TJmax, margin_offset
func (p *ParserPECI) GetMetrics() []string {
	return []string{"value", "margin_offset"}
}

// Parse method returns data in human readable format
func (p *ParserPECI) Parse(response IpmiResponse) map[string]uint64 {
	m := map[string]uint64{}
	// Based on Send raw PECI command (40h). Byte 7 returns margin offset
	// Bytes 8:9 returns TJmax
	m["margin_offset"] = uint64(response.Data[6])
	m["value"] = uint64(response.Data[7]) + uint64(response.Data[8])*256

	return m
}

// ParserPMBus extracts temperatures of voltage regulators.
type ParserPMBus struct {
	*GenericValidator
}

// Instance of ParserPMBus.
var FormatPMBus = &ParserPMBus{}

// GetMetrics returns metrics for PMBus parser: VR[0:5]
func (p *ParserPMBus) GetMetrics() []string {
	return []string{"VR/*/value"}
}

// Parse method returns data in human readable format
func (p *ParserPMBus) Parse(response IpmiResponse) map[string]uint64 {
	m := map[string]uint64{}
	var names = map[string]uint{"VR/0/value": 4, "VR/1/value": 6, "VR/2/value": 8, "VR/3/value": 10, "VR/4/value": 12, "VR/5/value": 14}
	//if not all VRs are present on platform return 0xFFFE
	if len(response.Data) < 14 {
		for metricName := range names {
			m[metricName] = 0xFFFE
		}
		return m
	}
	// Based on Send Raw PMBus Command (D9h). Bytes 9:N contains data received from PSU
	for metricName, startIndex := range names {
		m[metricName] = uint64(response.Data[startIndex]) + uint64(response.Data[startIndex+1])*256

	}
	return m
}

// ParserPSU extracts temperatures of PSU.
type ParserPSU struct {
	*GenericValidator
}

// Instance of ParserPSU.
var FormatPSU = &ParserPSU{}

// GetMetrics returns metrics for PSU Parser
func (p *ParserPSU) GetMetrics() []string {
	return []string{"0", "1"}
}

// Parse method returns data in human readable format
func (p *ParserPSU) Parse(response IpmiResponse) map[string]uint64 {
	m := map[string]uint64{}
	m["0"] = uint64(response.Data[4]) + uint64(response.Data[5])*256
	m["1"] = uint64(response.Data[6]) + uint64(response.Data[7])*256
	return m
}

// ParserSR extracts sensor value from response to Get Sensor Record.
type ParserSR struct {
	*GenericValidator
}

// Instance of ParserSR.
var FormatSR = &ParserSR{}

// GetMetrics returns metrics for sensor reading parser: current value
func (p *ParserSR) GetMetrics() []string {
	return []string{"value"}
}

// Parse method returns data in human readable format
func (p *ParserSR) Parse(response IpmiResponse) map[string]uint64 {
	m := map[string]uint64{}
	// Based on Get Sensor Reading (2Dh)
	m["value"] = uint64(response.Data[1])
	return m
}

// ParserPolicy extracts sensor value from response to Get Power Policy.
type ParserPolicy struct {
	*GenericValidator
}

// Instance of Power Policy parser
var FormatPolicy = &ParserPolicy{}

// GetMetrics returns metrics for power limit
func (p *ParserPolicy) GetMetrics() []string {
	return []string{"power_limit"}
}

// Parse method returns data in human readable format
func (p *ParserPolicy) Parse(response IpmiResponse) map[string]uint64 {
	m := map[string]uint64{}
	m["power_limit"] = uint64(response.Data[13]) + uint64(response.Data[14])*256
	return m
}

// ParserPolicy extracts sensor value from response to Get CUPS.
type ParserGetCups struct {
	*GenericValidator
}

// Instance of Cups parser
var FormatGetCups = &ParserGetCups{}

// GetMetrics returns metrics for cups
func (p *ParserGetCups) GetMetrics() []string {
	return []string{"value"}
}

// Parse method returns data in human readable format
func (p *ParserGetCups) Parse(response IpmiResponse) map[string]uint64 {
	m := map[string]uint64{}
	m["value"] = uint64(response.Data[1]) * 100 / 255
	return m
}
