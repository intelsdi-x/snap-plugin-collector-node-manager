// +build unit

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

// Tests for ipmi commands parser

package ipmi

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestCUPSParsing(t *testing.T) {
	Convey("Check CUPS parser", t, func() {
		validResponse := IpmiResponse{[]byte{0x00, 0x57, 0x01, 0x00, 0x64, 0x00, 0x50, 0x00, 0x00, 0x01}, "host", 1}
		a := &ParserCUPS{}
		metrics := a.GetMetrics()
		parserOut := a.Parse(validResponse)
		expects := []string{"cpu_bandwith", "memory_bandwith", "io_bandwith"}
		So(len(metrics), ShouldEqual, len(expects))
		for i := 0; i < len(expects); i++ {
			So(metrics[i], ShouldEqual, expects[i])
		}
		So(parserOut["cpu_bandwith"], ShouldEqual, 100)
		So(parserOut["memory_bandwith"], ShouldEqual, 80)
		So(parserOut["io_bandwith"], ShouldEqual, 256)
	})
}

func TestNodeManagerParsing(t *testing.T) {
	Convey("Check NodeManager parser", t, func() {

		validResponse := IpmiResponse{[]byte{0x00, 0x57, 0x01, 0x00, 0x69, 0x00, 0x03, 0x00, 0x7d, 0x01, 0x6E, 0x00, 0xC7, 0x3F, 0x05, 0x56, 0xB9, 0xAD, 0x0C, 0x00, 0x50}, "host", 1}
		a := &ParserNodeManager{}
		metrics := a.GetMetrics()
		parserOut := a.Parse(validResponse)
		expects := []string{"value", "min", "max", "avg"}
		So(len(metrics), ShouldEqual, len(expects))
		for i := 0; i < len(expects); i++ {
			So(metrics[i], ShouldEqual, expects[i])
		}
		So(parserOut["value"], ShouldEqual, 105)
		So(parserOut["min"], ShouldEqual, 3)
		So(parserOut["max"], ShouldEqual, 381)
		So(parserOut["avg"], ShouldEqual, 110)
	})
}

func TestPECIParsing(t *testing.T) {
	Convey("Check PECI parser", t, func() {
		validResponse := IpmiResponse{[]byte{0x00, 0x57, 0x01, 0x00, 0x40, 0x00, 0x0A, 0x59, 0x00}, "host", 1}
		a := &ParserPECI{}
		metrics := a.GetMetrics()
		parserOut := a.Parse(validResponse)
		expects := []string{"value", "margin_offset"}
		So(len(metrics), ShouldEqual, len(expects))
		for i := 0; i < len(expects); i++ {
			So(metrics[i], ShouldEqual, expects[i])
		}
		So(parserOut["value"], ShouldEqual, 89)
		So(parserOut["margin_offset"], ShouldEqual, 10)
	})
}

func TestPMBusParsing(t *testing.T) {
	Convey("Check PMBus parser", t, func() {
		validResponse := IpmiResponse{[]byte{0x00, 0x57, 0x01, 0x00, 0x25, 0x00, 0x2A, 0x00, 0x1F, 0x00, 0x21, 0x00, 0x20, 0x00, 0x1F, 0x00}, "host", 1}
		a := &ParserPMBus{}
		metrics := a.GetMetrics()
		parserOut := a.Parse(validResponse)
		expects := []string{"VR/*/value"}
		So(len(metrics), ShouldEqual, len(expects))
		for i := 0; i < len(expects); i++ {
			So(metrics[i], ShouldEqual, expects[i])
		}
		So(parserOut["VR/0/value"], ShouldEqual, 37)
		So(parserOut["VR/1/value"], ShouldEqual, 42)
		So(parserOut["VR/2/value"], ShouldEqual, 31)
		So(parserOut["VR/3/value"], ShouldEqual, 33)
		So(parserOut["VR/4/value"], ShouldEqual, 32)
		So(parserOut["VR/5/value"], ShouldEqual, 31)
	})
}
