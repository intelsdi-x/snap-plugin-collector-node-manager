// +build linux

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

package ipmi

import (
	"sync"
	"time"

	log "github.com/Sirupsen/logrus"
)

// LinuxOutOfBand implements communication with openipmi driver on linux
type LinuxOutOfBand struct {
	Device  string
	Channel string
	Slave   string
	Addr    []string
	User    string
	Pass    string
	mutex   sync.Mutex
}

// BatchExecRaw Performs batch of requests to given device.
// Returns array of responses in order corresponding to requests.
// Error is returned when any of requests failed.
func (al *LinuxOutOfBand) BatchExecRaw(requests []IpmiRequest, host string) ([]IpmiResponse, error) {
	var wg sync.WaitGroup
	wg.Add(len(requests))
	results := make([]IpmiResponse, len(requests))

	a := time.Now()
	for i, r := range requests {
		go func(i int, r IpmiRequest) {
			defer wg.Done()
			al.mutex.Lock()
			results[i] = fillStruct(r.Data, al, host)
			al.mutex.Unlock()
		}(i, r)
	}
	wg.Wait()
	b := time.Now()
	c := (b.Second() - a.Second())
	log.Debug("[COLLECTION] Collection took: ", c)

	return results, nil

}

func fillStruct(request []byte, strct *LinuxOutOfBand, addr string) IpmiResponse {
	var res IpmiResponse
	res.Data = ExecIpmiToolRemote(request, strct, addr)
	res.IsValid = 1
	return res
}

// GetPlatformCapabilities returns host NM capabilities
func (al *LinuxOutOfBand) GetPlatformCapabilities(requests []RequestDescription, host []string) map[string][]RequestDescription {
	validRequests := make(map[string][]RequestDescription, 0)
	var wg sync.WaitGroup

	a := time.Now()
	for _, addr := range host {
		validRequests[addr] = make([]RequestDescription, 0)
		wg.Add(len(requests))

		for _, req := range requests {
			go func(req RequestDescription, addr string) {
				al.mutex.Lock()
				a := ExecIpmiToolRemote(req.Request.Data, al, addr)
				al.mutex.Unlock()
				j := 0

				for i := range a {
					if a[i] == 0 {
						j++
					}
				}
				if j != len(a) {
					validRequests[addr] = append(validRequests[addr], req)
				}

				wg.Done()
			}(req, addr)
		}
		wg.Wait()
		b := time.Now()
		c := (b.Second() - a.Second())
		log.Debug("[INIT] Initialization took: ", c)
	}

	return validRequests

}
