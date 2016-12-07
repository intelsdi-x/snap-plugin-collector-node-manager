// +build unit

/*
http://www.apache.org/licenses/LICENSE-2.0.txt


Copyright 2016 Intel Corporation

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

// Tests for plugin functionality

package nodeManagerPlugin

import (
	"testing"

	"github.com/intelsdi-x/snap-plugin-utilities/str"
	"github.com/intelsdi-x/snap/control/plugin"
	. "github.com/smartystreets/goconvey/convey"
	"strings"
)

func TestGetMetrics(t *testing.T) {
	Convey("Check GetMetricsTypes", t, func() {
		ic := New()
		cfg := plugin.ConfigType{}
		mts, err := ic.GetMetricTypes(cfg)
		Convey("There no error should be reported", func() {
			So(err, ShouldBeNil)
		})
		Convey("Proper metrics should be returned", func() {
			metricNames := []string{}
			for _, m := range mts {
				metricNames = append(metricNames, m.Namespace().String())
			}
			So(len(mts), ShouldEqual, 42)
			So(str.Contains(metricNames, "/intel/node_manager/*/cups/dynamic_load_factor/cpu_bandwith"), ShouldBeTrue)
			So(str.Contains(metricNames, "/intel/node_manager/*/cups/dynamic_load_factor/memory_bandwith"), ShouldBeTrue)
			So(str.Contains(metricNames, "/intel/node_manager/*/cups/dynamic_load_factor/io_bandwith"), ShouldBeTrue)
			So(str.Contains(metricNames, "/intel/node_manager/*/cups/index"), ShouldBeTrue)
			So(str.Contains(metricNames, "/intel/node_manager/*/power/system/value"), ShouldBeTrue)
			So(str.Contains(metricNames, "/intel/node_manager/*/power/system/min"), ShouldBeTrue)
			So(str.Contains(metricNames, "/intel/node_manager/*/power/system/max"), ShouldBeTrue)
			So(str.Contains(metricNames, "/intel/node_manager/*/power/system/avg"), ShouldBeTrue)
			So(str.Contains(metricNames, "/intel/node_manager/*/power/cpu/value"), ShouldBeTrue)
			So(str.Contains(metricNames, "/intel/node_manager/*/power/cpu/min"), ShouldBeTrue)
			So(str.Contains(metricNames, "/intel/node_manager/*/power/cpu/max"), ShouldBeTrue)
			So(str.Contains(metricNames, "/intel/node_manager/*/power/cpu/avg"), ShouldBeTrue)
			So(str.Contains(metricNames, "/intel/node_manager/*/power/memory/value"), ShouldBeTrue)
			So(str.Contains(metricNames, "/intel/node_manager/*/power/memory/min"), ShouldBeTrue)
			So(str.Contains(metricNames, "/intel/node_manager/*/power/memory/max"), ShouldBeTrue)
			So(str.Contains(metricNames, "/intel/node_manager/*/power/memory/avg"), ShouldBeTrue)
			So(str.Contains(metricNames, "/intel/node_manager/*/temperature/inlet/value"), ShouldBeTrue)
			So(str.Contains(metricNames, "/intel/node_manager/*/temperature/inlet/min"), ShouldBeTrue)
			So(str.Contains(metricNames, "/intel/node_manager/*/temperature/inlet/max"), ShouldBeTrue)
			So(str.Contains(metricNames, "/intel/node_manager/*/temperature/inlet/avg"), ShouldBeTrue)
			So(str.Contains(metricNames, "/intel/node_manager/*/temperature/outlet/value"), ShouldBeTrue)
			So(str.Contains(metricNames, "/intel/node_manager/*/temperature/outlet/min"), ShouldBeTrue)
			So(str.Contains(metricNames, "/intel/node_manager/*/temperature/outlet/max"), ShouldBeTrue)
			So(str.Contains(metricNames, "/intel/node_manager/*/temperature/outlet/avg"), ShouldBeTrue)
			So(str.Contains(metricNames, "/intel/node_manager/*/temperature/cpu/*/value"), ShouldBeTrue)
			So(str.Contains(metricNames, "/intel/node_manager/*/temperature/dimm/*/value"), ShouldBeTrue)
			So(str.Contains(metricNames, "/intel/node_manager/*/temperature/psu/hot_spot/0"), ShouldBeTrue)
			So(str.Contains(metricNames, "/intel/node_manager/*/temperature/psu/hot_spot/1"), ShouldBeTrue)
			So(str.Contains(metricNames, "/intel/node_manager/*/temperature/psu/ambient/0"), ShouldBeTrue)
			So(str.Contains(metricNames, "/intel/node_manager/*/temperature/psu/ambient/1"), ShouldBeTrue)
			So(str.Contains(metricNames, "/intel/node_manager/*/airflow/value"), ShouldBeTrue)
			So(str.Contains(metricNames, "/intel/node_manager/*/airflow/min"), ShouldBeTrue)
			So(str.Contains(metricNames, "/intel/node_manager/*/airflow/max"), ShouldBeTrue)
			So(str.Contains(metricNames, "/intel/node_manager/*/airflow/avg"), ShouldBeTrue)
			So(str.Contains(metricNames, "/intel/node_manager/*/margin/cpu/tj/value"), ShouldBeTrue)
			So(str.Contains(metricNames, "/intel/node_manager/*/margin/cpu/tj/margin_offset"), ShouldBeTrue)
			So(str.Contains(metricNames, "/intel/node_manager/*/temperature/chipset/value"), ShouldBeTrue)
			So(str.Contains(metricNames, "/intel/node_manager/*/power/policy/power_limit"), ShouldBeTrue)
		})
	})
}

func TestExtendPath(t *testing.T) {
	Convey("Check extendPath", t, func() {
		a := []string{"intel", "node_manager", "127.0.0.1", "cups", "cpu_bandwith" }
		path := extendPath(a)
		So(strings.Contains(path, "/intel/node_manager/127.0.0.1/cups/cpu_bandwith"), ShouldBeTrue)

	})
}

func TestValidateResponse(t *testing.T) {
	Convey("Check validate response", t, func() {
		validResponse := []byte{0x00, 0x57, 0x01}
		invalidResponse1 := []byte{0x01, 0x57, 0x01}
		invalidResponse2 := []byte{0x00, 0x00, 0x01}
		invalidResponse3 := []byte{0x00, 0x57, 0x00}
		zeroResponse := []byte{}
		So(validateResponse(validResponse), ShouldBeTrue)
		So(validateResponse(invalidResponse1), ShouldBeFalse)
		So(validateResponse(invalidResponse2), ShouldBeFalse)
		So(validateResponse(invalidResponse3), ShouldBeFalse)
		So(validateResponse(zeroResponse), ShouldBeFalse)
	})
}
