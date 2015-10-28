// +build linux

/*
http://www.apache.org/licenses/LICENSE-2.0.txt


Copyright 2015 Intel Corporation

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
	"errors"
	"fmt"
)

// GenericValidator performs basic response validation. Checks response code ensures response
// has non-zero lenght.
type GenericValidator struct {
}

// Validate method verifies responses from IPMI device before running parsers
func (gv *GenericValidator) Validate(response []byte) error {
	if len(response) > 0 {
		if response[0] == 0 {
			return nil
		}
		return fmt.Errorf("Unexpected error code : %d", response[0])
	}
	return errors.New("Zero length response")

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
	return []string{"cpu_cstate", "memory_bandwith", "io_bandwith"}
}

// Parse method returns data in human readable format
func (p *ParserCUPS) Parse(response []byte) map[string]uint16 {
	m := map[string]uint16{}
	// Parsing is based on command Get CUPS Data (65h). Bytes 5:6 contains CPU CUPS dynamic load factor
	m["cpu_cstate"] = uint16(response[4]) + uint16(response[5])*256
	// Bytes 7:8 contains memory CUPS dynamic load factor
	m["memory_bandwith"] = uint16(response[6]) + uint16(response[7])*256
	// Bytes 9:10 contains IO CUPS dynamic load factor
	m["io_bandwith"] = uint16(response[8]) + uint16(response[9])*256
	return m
}

// ParserNodeManager extracs data from Node manager response format.
// Data contains current, min, max and average value.
type ParserNodeManager struct {
	*GenericValidator
}

// Instance of ParserNodeManager
var FormatNodeManager = &ParserNodeManager{}

// GetMetrics method returns metric for CUPS parser: "current_value", "min", "max", "avg"
func (p *ParserNodeManager) GetMetrics() []string {
	return []string{"", "min", "max", "avg"}
}

// Parse method returns data in human readable format
func (p *ParserNodeManager) Parse(response []byte) map[string]uint16 {
	m := map[string]uint16{}
	// Parsing is based on command Get Node Manager Statistics (C8h). Bytes 5:6 contains current value
	m[""] = uint16(response[4]) + uint16(response[5])*256
	// Bytes 7:8 contains minimum value
	m["min"] = uint16(response[6]) + uint16(response[7])*256
	// Bytes 9:10 contains maximum value
	m["max"] = uint16(response[8]) + uint16(response[9])*256
	// Bytes 11:12 contains average value
	m["avg"] = uint16(response[10]) + uint16(response[11])*256
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
	a := []string{"cpu/cpu0", "cpu/cpu1", "cpu/cpu2", "cpu/cpu3"}
	for i := 0; i < 64; i++ {
		c := fmt.Sprintf("memory/dimm%d", i)
		a = append(a, c)
	}
	return a
}

// Parse method returns data in human readable format
func (p *ParserTemp) Parse(response []byte) map[string]uint16 {
	m := map[string]uint16{}
	// Parsing is based on Get CPU and Memory Temperature (4Bh). Bytes 5:8 contains temperatures of each socket (up to 4)
	m["cpu/cpu0"] = uint16(response[4])
	m["cpu/cpu1"] = uint16(response[5])
	m["cpu/cpu2"] = uint16(response[6])
	m["cpu/cpu3"] = uint16(response[7])
	for i := 8; i < len(response); i++ {
		// Bytes 9:72 contains temperatures of each dimm (up to 64)
		a := fmt.Sprintf("memory/dimm%d", i-8)
		m[a] = uint16(response[i])
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
	return []string{"", "margin_offset"}
}

// Parse method returns data in human readable format
func (p *ParserPECI) Parse(response []byte) map[string]uint16 {
	m := map[string]uint16{}
	// Based on Send raw PECI command (40h). Byte 7 returns margin offset
	m["margin_offset"] = uint16(response[6])
	// Bytes 8:9 returns TJmax
	m[""] = uint16(response[7]) + uint16(response[8])*256
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
	return []string{"VR0", "VR1", "VR2", "VR3", "VR4", "VR5"}
}

// Parse method returns data in human readable format
func (p *ParserPMBus) Parse(response []byte) map[string]uint16 {
	m := map[string]uint16{}
	// Based on Send Raw PMBus Command (D9h). Bytes 9:N contains data received from PSU
	m["VR0"] = uint16(response[4]) + uint16(response[5])*256
	m["VR1"] = uint16(response[6]) + uint16(response[7])*256
	m["VR2"] = uint16(response[8]) + uint16(response[9])*256
	m["VR3"] = uint16(response[10]) + uint16(response[11])*256
	m["VR4"] = uint16(response[12]) + uint16(response[13])*256
	m["VR5"] = uint16(response[14]) + uint16(response[15])*256
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
	return []string{""}
}

// Parse method returns data in human readable format
func (p *ParserSR) Parse(response []byte) map[string]uint16 {
	m := map[string]uint16{}
	// Based on Get Sensor Reading (2Dh)
	m[""] = uint16(response[1])
	return m
}
