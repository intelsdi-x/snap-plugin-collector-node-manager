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

package main

import (
	"os"

	"github.com/intelsdi-x/snap-plugin-collector-node-manager/ipmi"
	"github.com/intelsdi-x/snap-plugin-collector-node-manager/node_manager_plugin"
	"github.com/intelsdi-x/snap/control/plugin"
)

func main() {

	ipmilayer := &ipmi.LinuxInband{Device: "/dev/ipmi0"}

	ipmiCollector := &node_manager_plugin.IpmiCollector{IpmiLayer: ipmilayer,
		Vendor: ipmi.GenericVendor, NSim: 3}

	plugin.Start(plugin.NewPluginMeta(node_manager_plugin.Name, node_manager_plugin.Version,
		node_manager_plugin.Type, []string{}, []string{plugin.SnapGOBContentType},
		plugin.ConcurrencyCount(1)), ipmiCollector, os.Args[1])
}
