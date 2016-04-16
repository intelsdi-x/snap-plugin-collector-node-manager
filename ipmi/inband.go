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
	"fmt"
	"sync"
	"unsafe"
)

// #include "linux_inband.h"
import "C"

// LinuxInband Implements communication with openipmi driver on linux
type LinuxInband struct {
	Device string
	mutex  sync.Mutex
}

// BatchExecRaw Performs batch of requests to given device.
// Returns array of responses in order corresponding to requests.
// Error is returned when any of requests failed.
func (al *LinuxInband) BatchExecRaw(requests []IpmiRequest, _ string) ([]IpmiResponse, error) {
	al.mutex.Lock()
	defer al.mutex.Unlock()

	n := len(requests)
	info := C.struct_IpmiStatusInfo{}
	inputs := make([]C.struct_IpmiCommandInput, n)
	outputs := make([]C.struct_IpmiCommandOutput, n)

	for i, r := range requests {
		for j, b := range r.Data {
			inputs[i].data[j] = C.char(b)
		}
		inputs[i].data_len = C.int(len(r.Data))
		inputs[i].channel = C.short(r.Channel)
		inputs[i].slave = C.uchar(r.Slave)
	}

	errcode := C.IPMI_BatchCommands(C.CString(al.Device), &inputs[0], &outputs[0],
		C.int(n), C.int(3), &info)

	switch {
	case errcode < 0:
		return nil, fmt.Errorf("%d : Invalid call", errcode)
	case errcode > 0:
		return nil, fmt.Errorf("%d : System error [%d : %s]", errcode,
			info.system_error, C.GoString(&info.error_str[0]))
	}

	results := make([]IpmiResponse, n)

	for i, r := range outputs {
		results[i].Data = C.GoBytes(unsafe.Pointer(&r.data[0]), r.data_len)
		results[i].IsValid = uint(r.is_valid)
	}

	return results, nil
}

// GetPlatformCapabilities returns host NM capabilities
func (al *LinuxInband) GetPlatformCapabilities(requests []RequestDescription, host []string) map[string][]RequestDescription {
	return nil
}
