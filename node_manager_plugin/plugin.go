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
	"sync"
	"time"

	"github.com/intelsdi-x/snap-plugin-collector-node-manager/ipmi"
	"github.com/intelsdi-x/snap/control/plugin"
	"github.com/intelsdi-x/snap/control/plugin/cpolicy"
)

const (
	Name    = "snap-collector-intel-node-manager"
	Version = 2
	Type    = plugin.CollectorPluginType
)

var namespace_prefix = []string{"intel", "node_manager"}

func makeName(metric string) []string {
	return append(namespace_prefix, strings.Split(metric, "/")...)
}

func parseName(namespace []string) string {
	return strings.Join(namespace[len(namespace_prefix):], "/")
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
	IpmiLayer ipmi.IpmiAL
	NSim      int
	Vendor    []ipmi.RequestDescription

	requestIndexCache     map[string]int
	requestIndexCacheOnce sync.Once
}

func (ic *IpmiCollector) ensureRequestIdxCache() {
	ic.requestIndexCache = make(map[string]int)

	for i, req := range ic.Vendor {
		for _, metric := range req.Format.GetMetrics() {
			path := extendPath(req.MetricsRoot, metric)
			ic.requestIndexCache[path] = i
		}
	}
}

func (ic *IpmiCollector) validateName(namespace []string) error {
	for i, e := range namespace_prefix {
		if namespace[i] != e {
			return fmt.Errorf("Wrong namespace prefix in namespace %v", namespace)
		}
	}

	name := parseName(namespace)
	_, ok := ic.requestIndexCache[name]

	if !ok {
		return fmt.Errorf("Key %s not in index cache", name)
	}

	return nil
}

// Performs metric collection.
// Ipmi request are never duplicated in order to read multiple metrics.
// Timestamp is set to time when batch processing is complete.
// Source is hostname returned by operating system.
func (ic *IpmiCollector) CollectMetrics(mts []plugin.PluginMetricType) ([]plugin.PluginMetricType, error) {

	ic.requestIndexCacheOnce.Do(func() {
		ic.ensureRequestIdxCache()
	})

	requestSet := map[int]bool{}
	for _, mt := range mts {
		ns := mt.Namespace()
		if err := ic.validateName(ns); err != nil {
			return nil, err
		}
		rid := ic.requestIndexCache[parseName(ns)]
		requestSet[rid] = true
	}

	requestList := make([]ipmi.IpmiRequest, 0)
	requestDescList := make([]*ipmi.RequestDescription, 0)
	responseCache := map[string]uint16{}
	for k, _ := range requestSet {
		desc := &ic.Vendor[k]
		requestList = append(requestList, desc.Request)
		requestDescList = append(requestDescList, desc)
	}

	//TODO: nSim from config
	resp, err := ic.IpmiLayer.BatchExecRaw(requestList, ic.NSim)

	if err != nil {
		return nil, err
	}

	valid_metrics := len(mts)
	for i, r := range resp {
		format := requestDescList[i].Format
		if err := format.Validate(r); err != nil {
			valid_metrics--
			submetrics := format.GetMetrics()
			for _, submetric := range submetrics {
				path := extendPath(requestDescList[i].MetricsRoot, submetric)
				responseCache[path] = 0xFFFF
			}
		} else {
			submetrics := format.Parse(r)
			for k, v := range submetrics {
				path := extendPath(requestDescList[i].MetricsRoot, k)
				responseCache[path] = v
			}
		}
	}

	results := make([]plugin.PluginMetricType, valid_metrics)
	t := time.Now()
	host, _ := os.Hostname()

	for i, mt := range mts {
		ns := mt.Namespace()
		key := parseName(ns)
		// to return incomplete metrics remove condition
		if responseCache[key] != 0xFFFF {
			data := responseCache[key]
			metric := plugin.PluginMetricType{Namespace_: ns, Source_: host, Timestamp_: t, Data_: data}
			results[i] = metric
		}
	}

	return results, nil
}

// Returns list of metrics available for current vendor.
func (ic *IpmiCollector) GetMetricTypes(_ plugin.PluginConfigType) ([]plugin.PluginMetricType, error) {
	mts := make([]plugin.PluginMetricType, 0)
	for _, req := range ic.Vendor {
		for _, metric := range req.Format.GetMetrics() {
			path := extendPath(req.MetricsRoot, metric)
			mts = append(mts, plugin.PluginMetricType{Namespace_: makeName(path)})
		}
	}

	return mts, nil
}

func (ic *IpmiCollector) GetConfigPolicy() (*cpolicy.ConfigPolicy, error) {
	c := cpolicy.New()
	return c, nil
}
