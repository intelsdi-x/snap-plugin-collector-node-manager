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

package nodeManagerPlugin

import (
	"os"
	"strings"
	"time"

	"bufio"
	"github.com/intelsdi-x/snap-plugin-collector-node-manager/ipmi"
	"github.com/intelsdi-x/snap-plugin-utilities/config"
	"github.com/intelsdi-x/snap/control/plugin"
	"github.com/intelsdi-x/snap/control/plugin/cpolicy"
	"github.com/intelsdi-x/snap/core"
	"sync"
	"fmt"
	log "github.com/Sirupsen/logrus"
)

const (
	//Name is name of plugin
	Name = "node-manager"
	//Version of plugin
	Version = 8
	//Type of plugin
	Type = plugin.CollectorPluginType

	pluginVendor = "intel"
	pluginName   = "node_manager"
)

func extendPath(path []string) string {
	fullPath := string("")
	if len(path) > 0 {
		for _, ext := range path {
			fullPath += "/" + ext
		}
	}
	return fullPath
}

// IpmiCollector Plugin class.
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
}

func validateResponse(response []byte) bool {
	if len(response) > 0 {
		if response[0] == 0x00 && response[1] == 0x57 && response[2] == 0x01 {
			return true
		} else if len(response) == 4 && response[0] == 0x00 && response[2] == 0xC0 { //validate Sensor reading
			return true
		}
	}
	log.Debug("Response Invalid")
	return false
}

func (ic *IpmiCollector) CollectMetrics(mts []plugin.MetricType) ([]plugin.MetricType, error) {
	if !ic.Initialized {
		log.Debug("Plugin not initialized!")
		ic.construct(mts[0]) //reinitialize plugin
	}
	var wg sync.WaitGroup
	c := make(chan ipmi.IpmiResponse, len(ic.Hosts)*len(ipmi.GenericVendor))
	wg.Add(len(ic.Hosts) * len(ipmi.GenericVendor))

	log.Debug("Collection started")
	for _, host := range ic.Hosts {
		for i, req := range ipmi.GenericVendor {
			go func(host string, i int, req ipmi.RequestDescription) {
				defer wg.Done()
				c <- ic.IpmiLayer.RunParallelRequests(req.Request, host, i)
			}(host, i, req)
		}
	}
	wg.Wait()
	close(c)
	log.Debug("Collection done")
	values := make([]ipmi.IpmiResponse, 0)
	for a := range c {
		values = append(values, a)
	}
	t := time.Now()
	responseMetrics := make([]plugin.MetricType, 0)
	responseCache := make(map[string]uint64)
	log.Debug("Processing started")
	for _, value := range values {
		if len(value.Data) > 0 {
			if validateResponse(value.Data) {
				format := ipmi.GenericVendor[value.Index].Format
				metricRoot := ipmi.GenericVendor[value.Index].MetricsRoot
				submetrics := format.Parse(value)
				for k, v := range submetrics {
					path := extendPath([]string{value.Source, metricRoot, k})
					responseCache[path] = v
				}
			}
		}
	}
	var metric plugin.MetricType

	for key, resp := range responseCache {
		args := []string{pluginVendor, pluginName}
		args = append(args, strings.Split(key, "/")[1:]...)
		source := args[2]
		ns := core.NewNamespace(args[0:2]...).
			AddDynamicElement("host_id", "Host ID").
			AddStaticElements(args[3:]...)
		ns[2].Value = fmt.Sprintf("%s", source)
		tags := map[string]string{"source": source}
		metric = plugin.MetricType{
			Namespace_: ns,
			Data_:      resp,
			Timestamp_: t,
			Tags_:      tags,
		}
		responseMetrics = append(responseMetrics, metric)
	}
	log.Debug("Processing done")

	return responseMetrics, nil
}

// GetMetricTypes Returns list of metrics available for current vendor.
func (ic *IpmiCollector) GetMetricTypes(cfg plugin.ConfigType) ([]plugin.MetricType, error) {
	var mts []plugin.MetricType
	mts = make([]plugin.MetricType, 0)
	var namespace []core.NamespaceElement

	for _, req := range ipmi.GenericVendor {
		for _, metric := range req.Format.GetMetrics() {
			if strings.Contains(metric, "*") {
				if strings.Contains(metric, "cpu") {
					ns := strings.Split(req.MetricsRoot, "/")
					ns = append(ns, strings.Split(metric, "/")[0])
					namespace = core.NewNamespace(pluginVendor, pluginName).
						AddDynamicElement("host_id", "Host ID").
						AddStaticElements(ns...).
						AddDynamicElement("cpu_id", "CPU ID").
						AddStaticElement("value")
				} else if strings.Contains(metric, "dimm") {
					ns := strings.Split(req.MetricsRoot, "/")
					ns = append(ns, strings.Split(metric, "/")[0])
					namespace = core.NewNamespace(pluginVendor, pluginName).
						AddDynamicElement("host_id", "Host ID").
						AddStaticElements(ns...).
						AddDynamicElement("dimm_id", "DIMM ID").
						AddStaticElement("value")
				}
			} else {
				namespace = core.NewNamespace(pluginVendor, pluginName).
					AddDynamicElement("host_id", "Host ID").
					AddStaticElements(strings.Split(req.MetricsRoot, "/")...).
					AddStaticElement(metric)
			}
			mts = append(mts, plugin.MetricType{Namespace_: namespace})
		}
	}
	return mts, nil
}

// GetConfigPolicy creates policy based on global config
func (ic *IpmiCollector) GetConfigPolicy() (*cpolicy.ConfigPolicy, error) {
	c := cpolicy.New()
	return c, nil
}

// New is simple collector constuctor
func New() *IpmiCollector {
	collector := &IpmiCollector{Initialized: false}
	return collector
}

func getHosts(config string) []string {
	if config == "" {
		return nil
	}
	file, err := os.Open(config)
	if err != nil {
		log.Debug("Unable to open file with list of hosts")
		return nil
	}
	defer file.Close()
	var hostList []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		hostList = append(hostList, scanner.Text())
	}
	return hostList
}

func (ic *IpmiCollector) construct(mts plugin.MetricType) {
	var hostList []string
	var ipmiLayer ipmi.IpmiAL
	ic.Initialized = false
	mode, err := config.GetConfigItem(mts, "mode")
	if err != nil {
		log.Debug("Invalid configuration mode")
		log.Debug(mode)
		log.Debug("Unable to initialize plugin")
		return
	}
	ic.Mode = mode.(string)
	host, _ := os.Hostname()
	if ic.Mode == "legacy_inband" {
		configuration, err := config.GetConfigItems(mts, "channel", "slave")
		if err != nil {
			log.Debug("Invalid configuration for legacy_inband")
			log.Debug(configuration)
			log.Debug("Unable to initialize plugin")
			return
		}
		ipmiLayer = &ipmi.LinuxInBandIpmitool{
			Device:  "ipmitool",
			Channel: configuration["channel"].(string),
			Slave:   configuration["slave"].(string),
		}
		hostList = []string{host}
	} else if ic.Mode == "legacy_oob" {
		configuration, err := config.GetConfigItems(mts, "channel", "slave", "user", "password")
		if err != nil {
			log.Debug("Invalid configuration for legacy_oob")
			log.Debug(configuration)
			log.Debug("Unable to initialize plugin")
			return
		}
		ipmiLayer = &ipmi.LinuxOutOfBand{
			Device:  "ipmitool",
			Channel: configuration["channel"].(string),
			Slave:   configuration["slave"].(string),
			User:    configuration["user"].(string),
			Pass:    configuration["password"].(string),
		}
		path, err := config.GetConfigItem(mts, "hosts")
		if err != nil {
			log.Debug("Unable to get hosts list")
			log.Debug(path)
			log.Debug("Unable to initialize plugin")
			return
		}
		hostList = getHosts(path.(string))
		if hostList == nil {
			log.Debug("Host list empty")
			return 
		}
	} else {
		log.Debug("Invalid mode configuration")
		log.Debug("Unable to initialize plugin")
		return
	}

	ic.IpmiLayer = ipmiLayer
	ic.Hosts = hostList
	ic.Initialized = true

}
