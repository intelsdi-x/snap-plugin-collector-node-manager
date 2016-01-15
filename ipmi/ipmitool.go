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
	"os/exec"
	"strconv"
	"strings"
)

// ExecIpmiToolLocal method runs ipmitool command on a local system
func ExecIpmiToolLocal(request []byte, strct *LinuxInBandIpmitool) []byte {

	c, err := exec.LookPath("ipmitool")
	if err != nil {
		fmt.Printf("Unable to find %s", "ipmitool")
		return []byte{0x01, 0x02}
	}

	stringRequest := []string{"-b", strct.Channel, "-t", strct.Slave, "raw"}
	for i := range request {
		stringRequest = append(stringRequest, fmt.Sprintf("0x%02x", request[i]))
	}
	ret, err := exec.Command(c, stringRequest...).CombinedOutput()
	returnStrings := strings.Split(string(ret), " ")
	rets := make([]byte, len(returnStrings))
	for i, element := range returnStrings {
		value, _ := strconv.ParseInt(element, 16, 0)
		rets[i] = byte(value)
	}
	return rets

}
