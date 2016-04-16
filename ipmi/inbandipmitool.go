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
	"os"
	"sync"
)

// LinuxInBandIpmitool implements communication with ipmitool on linux
type LinuxInBandIpmitool struct {
	Device  string
	Channel string
	Slave   string
	mutex   sync.Mutex
}

// BatchExecRaw performs batch of requests to given device.
// Returns array of responses in order corresponding to requests.
// Error is returned when any of requests failed.
func (al *LinuxInBandIpmitool) BatchExecRaw(requests []IpmiRequest, host string) ([]IpmiResponse, error) {
	al.mutex.Lock()
	defer al.mutex.Unlock()

	results := make([]IpmiResponse, len(requests))

	for i, r := range requests {
		results[i].Data = ExecIpmiToolLocal(r.Data, al)
		results[i].IsValid = 1
	}

	return results, nil

}

// GetPlatformCapabilities returns host NM capabilities
func (al *LinuxInBandIpmitool) GetPlatformCapabilities(requests []RequestDescription, _ []string) map[string][]RequestDescription {
	host, _ := os.Hostname()
	validRequests := make(map[string][]RequestDescription, 0)
	validRequests[host] = make([]RequestDescription, 0)

	for _, request := range requests {
		response := ExecIpmiToolLocal(request.Request.Data, al)
		j := 0

		for i := range response {
			if response[i] == 0 {
				j++
			}
		}
		if j != len(response) {
			validRequests[host] = append(validRequests[host], request)
		}
	}

	return validRequests

}
