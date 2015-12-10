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

package node_manager_plugin

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/intelsdi-x/snap-plugin-collector-node-manager/ipmi"
	"github.com/intelsdi-x/snap/control/plugin"
	"github.com/intelsdi-x/snap/control/plugin/cpolicy"
	"github.com/intelsdi-x/snap/core/ctypes"
)

const (
	Name    = "node-manager"
	Version = 5
	Type    = plugin.CollectorPluginType
)

var namespacePrefix = []string{"intel", "node_manager"}

func makeName(metric string) []string {
	return append(namespacePrefix, strings.Split(metric, "/")...)
}

func parseName(namespace []string) string {
	return strings.Join(namespace[len(namespacePrefix):], "/")
}

func extendPath(path, ext string) string {
	if ext == "" {
		return path
	} else {
		return path + "/" + ext
	}
}

// Ipmi Collector Plugin class.
// IpmiLayer specifies interface to perform ipmi commands.
// NSim is number of requests allowed to be 'in processing' state.
// Vendor is list of request descriptions. Each of them specifies
// RAW request data, root path for metrics
// and format (which also specifies submetrics)
type IpmiCollector struct {
	IpmiLayer   ipmi.IpmiAL
	Vendor      map[string][]ipmi.RequestDescription
	Hosts       []string
	Mode        string
	Initialized bool
	NSim        int
}

func (ic *IpmiCollector) validateName(namespace []string) error {
	for i, e := range namespacePrefix {
		if namespace[i] != e {
			return fmt.Errorf("Wrong namespace prefix in namespace %v", namespace)
		}
	}
	return nil
}

// Performs metric collection.
// Ipmi request are never duplicated in order to read multiple metrics.
// Timestamp is set to time when batch processing is complete.
// Source is hostname returned by operating system.
func (ic *IpmiCollector) CollectMetrics(mts []plugin.PluginMetricType) ([]plugin.PluginMetricType, error) {
	if !ic.Initialized {
		ic.construct(mts[0].Config().Table()) //reinitialize plugin
	}
	requestList := make(map[string][]ipmi.IpmiRequest, 0)
	requestDescList := make(map[string][]ipmi.RequestDescription, 0)
	responseCache := map[string]map[string]uint16{}
	for _, host := range ic.Hosts {
		requestList[host] = make([]ipmi.IpmiRequest, 0)
		requestDescList[host] = make([]ipmi.RequestDescription, 0)
		for _, request := range ic.Vendor[host] {
			requestList[host] = append(requestList[host], request.Request)
			requestDescList[host] = append(requestDescList[host], request)
		}
	}
	response := make(map[string][]ipmi.IpmiResponse, 0)

	for _, host := range ic.Hosts {
		response[host], _ = ic.IpmiLayer.BatchExecRaw(requestList[host], host)
	}

	for nmResponseIdx, hostResponses := range response {
		cached := map[string]uint16{}
		for i, resp := range hostResponses {
			format := requestDescList[nmResponseIdx][i].Format
			if err := format.Validate(resp); err != nil {
				return nil, err
			}
			submetrics := format.Parse(resp)
			for k, v := range submetrics {
				path := extendPath(requestDescList[nmResponseIdx][i].MetricsRoot, k)
				cached[path] = v
			}
			responseCache[nmResponseIdx] = cached
		}
	}

	results := make([]plugin.PluginMetricType, len(mts))
	responseMetrics := make([]plugin.PluginMetricType, 0)
	t := time.Now()

	for _, host := range ic.Hosts {
		for i, mt := range mts {
			ns := mt.Namespace()
			key := parseName(ns)
			data := responseCache[host][key]
			metric := plugin.PluginMetricType{Namespace_: ns, Source_: host,
				Timestamp_: t, Data_: data}
			results[i] = metric
			responseMetrics = append(responseMetrics, metric)
		}
	}

	return responseMetrics, nil
}

func getMode(config map[string]ctypes.ConfigValue) string {
	if mode, ok := config["mode"]; ok {
		return mode.(ctypes.ConfigValueStr).Value
	}
	return "legacy_inband" //Default mode
}

func getChannel(config map[string]ctypes.ConfigValue) string {
	if channel, ok := config["channel"]; ok {
		return channel.(ctypes.ConfigValueStr).Value
	}
	return "0x00" //Default channel addr
}

func getSlave(config map[string]ctypes.ConfigValue) string {
	if slave, ok := config["slave"]; ok {
		return slave.(ctypes.ConfigValueStr).Value
	}
	return "0x00" //Default slave addr
}

func (ic *IpmiCollector) construct(cfg map[string]ctypes.ConfigValue) {
	var hostList []string
	var ipmiLayer ipmi.IpmiAL
	ic.Mode = getMode(cfg)
	channel := getChannel(cfg)
	slave := getSlave(cfg)

	host, _ := os.Hostname()
	fmt.Println(host)
	if ic.Mode == "legacy_inband" {
		ipmiLayer = &ipmi.LinuxInBandIpmitool{Device: "ipmitool", Channel: channel, Slave: slave}
		hostList = []string{host}
	} else {
		return
	}
	ic.IpmiLayer = ipmiLayer
	ic.Hosts = hostList
	ic.Vendor = ipmiLayer.GetPlatformCapabilities(ipmi.GenericVendor, hostList)

}

// Returns list of metrics available for current vendor.
func (ic *IpmiCollector) GetMetricTypes(cfg plugin.PluginConfigType) ([]plugin.PluginMetricType, error) {
	ic.construct(cfg.Table())
	mts := make([]plugin.PluginMetricType, 0)
	if ic.IpmiLayer == nil {
		return mts, fmt.Errorf("Wrong mode configuration")
	}
	for _, host := range ic.Hosts {
		for _, req := range ic.Vendor[host] {
			for _, metric := range req.Format.GetMetrics() {
				path := extendPath(req.MetricsRoot, metric)
				mts = append(mts, plugin.PluginMetricType{Namespace_: makeName(path), Source_: host})
			}
		}
	}
	ic.Initialized = true
	return mts, nil
}

func (ic *IpmiCollector) GetConfigPolicy() (*cpolicy.ConfigPolicy, error) {
	c := cpolicy.New()
	return c, nil
}

func New() *IpmiCollector {
	collector := &IpmiCollector{Initialized: false}
	return collector
}
