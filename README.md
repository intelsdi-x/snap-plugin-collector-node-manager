# snap collector plugin - Intel Node Manager

 Plugin to collect data from Intel's Node Manager. Which is presenting low level metrics like power consumption, cpu temperature, etc.
 Currently it is using Ipmitool to collect data from NM.

1. [Getting Started](#getting-started)
  * [System Requirements](#system-requirements)
  * [Installation](#installation)
  * [Configuration and Usage](configuration-and-usage)
2. [Documentation](#documentation)
  * [Collected Metrics](#collected-metrics)
  * [Examples](#examples)
  * [Roadmap](#roadmap)
3. [Community Support](#community-support)
4. [Contributing](#contributing)
5. [License](#license)
6. [Acknowledgements](#acknowledgements)

## Getting Started

 Plugin collects specified metrics in-band on OS level

### System Requirements

 - Plugin needs to be run on server platform which supports Intel Node Manager.
 - Currently it works only on Linux Servers
 - Ipmitool needs to be installed on platform

### Installation
#### Download Intel Node Manager plugin binary:
You can get the pre-built binaries for your OS and architecture at snap's [Github Releases](https://github.com/intelsdi-x/snap/releases) page.

#### To build the plugin binary:
Fork https://github.com/intelsdi-x/snap-plugin-collector-node-manager
Clone repo into `$GOPATH/src/github/intelsdi-x/`:
```
$ git clone https://github.com/<yourGithubID>/snap-plugin-collector-node-manager
```
Build the plugin by running make in repo:
```
$ make
```
This builds the plugin in `/build/rootfs`

### Configuration and Usage

 On OS level user needs to load modules:
  - ipmi_msghandler
  - ipmi_devintf
  - ipmi_si
 
Those modules provides specific IPMI device which can collect data from NM

There are currently 6 configuration options:
 - mode - defines mode of plugin work, possible values: legacy_inband, legacy_oob
 - channel - defines communication channel address (default: "0x00")
 - slave - defines target address (default: "0x00")
 - user - for OOB mode only, user for authentication to remote host
 - password - for OOB mode only, password for authentication to remote host
 - hosts - for OOB mode only, path to file with IPs of host which will be monitored OOB

Sample configuration of node manager plugin:
```
{
    "plugins": {
        "collector": {
            "node-manager": {
                "all": {
                    "mode": "legacy_inband",
                    "channel": "0x00",
                    "slave": "0x00"
                }
            }
        }
    }
}
```

## Documentation

### Collected Metrics
This plugin has the ability to gather the following metrics:

Namespace | Data Type | Description (optional)
----------|-----------|-----------------------
/intel/node_manager/host_id/airflow | uint64 | Current Volumetric Airflow
/intel/node_manager/host_id/airflow/avg | uint64 | Average Volumetric Airflow 
/intel/node_manager/host_id/airflow/max | uint64 | Maximal Volumetric Airflow 
/intel/node_manager/host_id/airflow/min | uint64 | Minimal Volumetric Airflow 
/intel/node_manager/host_id/cups/cpu_cstate | uint64 | CUPS CPU Bandwidth
/intel/node_manager/host_id/cups/io_bandwith | uint64 | CUPS I/O Bandwidth
/intel/node_manager/host_id/cups/memory_bandwith | uint64 | CUPS Memory Bandwidth
/intel/node_manager/host_id/power/cpu | uint64 | Current CPU power consumption
/intel/node_manager/host_id/power/cpu/avg | uint64 | Average CPU power consumption
/intel/node_manager/host_id/power/cpu/max | uint64 | Maximal CPU power consumption
/intel/node_manager/host_id/power/cpu/min | uint64 | Minimal CPU power consumption
/intel/node_manager/host_id/power/policy/power_limit | uint64 | Power policy
/intel/node_manager/host_id/margin/cpu/tj  | uint64 | Margin-to-throttle functional  (CPU)
/intel/node_manager/host_id/margin/cpu/tj/margin_offset | uint64 | Margin-to-spec reliability (CPU)
/intel/node_manager/host_id/power/memory | uint64 | Current Memory power consumption
/intel/node_manager/host_id/power/memory/avg | uint64 | Average Memory power consumption
/intel/node_manager/host_id/power/memory/max | uint64 | Maximal Memory power consumption
/intel/node_manager/host_id/power/memory/min | uint64 | Minimal Memory power consumption
/intel/node_manager/host_id/power/system | uint64 | Current Platform power consumption
/intel/node_manager/host_id/power/system/avg | uint64 | Average Platform power consumption
/intel/node_manager/host_id/power/system/max | uint64 | Maximal Platform power consumption
/intel/node_manager/host_id/power/system/min | uint64 | Minimal Platform power consumption
/intel/node_manager/host_id/temperature/cpu/<cpu_id> | uint64 | Current CPU temperature
/intel/node_manager/host_id/temperature/pmbus/VR/<VR_id> | uint64 | Current VR's temperature
/intel/node_manager/host_id/temperature/dimm/<dimm_id> | uint64 | Current Memory dimms temperature
/intel/node_manager/host_id/temperature/outlet | uint64 | Current Outlet (exhaust air) temperature
/intel/node_manager/host_id/temperature/outlet/avg | uint64 | Average Outlet (exhaust air) temperature
/intel/node_manager/host_id/temperature/outlet/max | uint64 | Maximal Outlet (exhaust air) temperature
/intel/node_manager/host_id/temperature/outlet/min | uint64 | Minimal Outlet (exhaust air) temperature
/intel/node_manager/host_id/temperature/inlet | uint64 | Current Inlet Temperature
/intel/node_manager/host_id/temperature/inlet/avg | uint64 | Average Inlet Temperature
/intel/node_manager/host_id/temperature/inlet/max | uint64 | Maximal Inlet Temperature
/intel/node_manager/host_id/temperature/inlet/min | uint64 | Minimal Inlet Temperature

### Metric Tags
Namespace | Tag | Description
----------|-----|------------
/intel/node_manager/* | source | Host IP address

### Examples
Example task manifest to use Intel Node Manager plugin:
```
{
    "version": 1,
    "schedule": {
        "type": "simple",
        "interval": "5s"
    },
    "workflow": {
        "collect": {
            "metrics": {
                "/intel/node_manager/*/airflow": {},
                "/intel/node_manager/*/cups/cpu_cstate": {},
                "/intel/node_manager/*/cups/io_bandwith": {},
                "/intel/node_manager/*/cups/memory_bandwith": {},
                "/intel/node_manager/*/margin/cpu/tj": {},
                "/intel/node_manager/*/margin/cpu/tj/margin_offset": {},
                "/intel/node_manager/*/power/cpu": {},
                "/intel/node_manager/*/power/memory": {},
                "/intel/node_manager/*/power/memory/avg": {},
                "/intel/node_manager/*/power/memory/max": {},
                "/intel/node_manager/*/power/memory/min": {},
                "/intel/node_manager/*/power/system": {},
                "/intel/node_manager/*/power/system/avg": {},
                "/intel/node_manager/*/power/system/max": {},
                "/intel/node_manager/*/power/system/min": {},
                "/intel/node_manager/*/temperature/cpu/*": {},
                "/intel/node_manager/*/temperature/outlet": {},
            },
            "config": {
            },
            "process": null,
            "publish": [
                {
                    "plugin_name": "file",
                    "plugin_version": 2,
                    "config": {
                        "file": "/tmp/published"
                    }
                }
            ]
        }
    }
}
```


### Roadmap
As we launch this plugin, we have a few items in mind for the next release:
- SDR readings support

## Community Support
This repository is one of **many** plugins in **snap**, a powerful telemetry framework. See the full project at http://github.com/intelsdi-x/snap To reach out to other users, head to the [main framework](https://github.com/intelsdi-x/snap#community-support)

## Contributing
We love contributions!

There's more than one way to give back, from examples to blogs to code updates. See our recommended process in [CONTRIBUTING.md](CONTRIBUTING.md).

## License
[snap](http://github.com:intelsdi-x/snap), along with this plugin, is an Open Source software released under the Apache 2.0 [License](LICENSE).

## Acknowledgements

* Author: [Lukasz Mroz](https://github.com/lmroz)
* Author: [Marcin Krolik](https://github.com/marcin-krolik)
* Author: [Patryk Matyjasek](https://github.com/PatrykMatyjasek)

And **thank you!** Your contribution, through code and participation, is incredibly important to us.
