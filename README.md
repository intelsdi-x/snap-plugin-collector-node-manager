DISCONTINUATION OF PROJECT. 

This project will no longer be maintained by Intel.

This project has been identified as having known security escapes.

Intel has ceased development and contributions including, but not limited to, maintenance, bug fixes, new releases, or updates, to this project.  

Intel no longer accepts patches to this project.
# DISCONTINUATION OF PROJECT 

**This project will no longer be maintained by Intel.  Intel will not provide or guarantee development of or support for this project, including but not limited to, maintenance, bug fixes, new releases or updates.  Patches to this project are no longer accepted by Intel. If you have an ongoing need to use this project, are interested in independently developing it, or would like to maintain patches for the community, please create your own fork of the project.**

# Plugin status

This plugin is no longer being actively maintained by intel. Work has instead shifted to the [intel-dcm-platform](https://github.com/intelsdi-x/snap-plugin-collector-intel-dcm-platform) plugin.

# snap collector plugin - Intel Node Manager

 Plugin to collect data from Intel's Node Manager. Which is presenting low level metrics like power consumption, cpu temperature, etc.
 Currently it is using Ipmitool to collect data from NM.

1. [Getting Started](#getting-started)
  * [System Requirements](#system-requirements)
  * [Installation](#installation)
  * [Configuration and Usage](#configuration-and-usage)
2. [Documentation](#documentation)
  * [Collected Metrics](#collected-metrics)
  * [Examples](#examples)
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
You can get the pre-built binaries for your OS and architecture at plugin's [Github Releases](releases) page.

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
This builds the plugin in `./build/`

### Configuration and Usage

 On OS level user needs to load modules:
  - ipmi_msghandler
  - ipmi_devintf
  - ipmi_si
 
Those modules provides specific IPMI device which can collect data from NM

There are currently 6 configuration options:
 - mode - defines mode of plugin work, possible values: legacy_inband, legacy_inband_openipmi, oob
 - channel - defines communication channel address (default: "0x00")
 - slave - defines target address (default: "0x00")
 - user - for OOB mode only, user for authentication to remote host
 - password - for OOB mode only, password for authentication to remote host
 - host - for OOB mode only, IP of host which will be monitored OOB

Sample configuration of node manager plugin:
```
{
    "control": {
        "plugins": {
            "collector": {
                "node-manager": {
                    "all": {
                        "mode": "legacy_inband",
                        "channel": "6",
                        "slave": "0x2c"
                    }
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
/intel/node_manager/airflow | uint16 | Current Volumetric Airflow
/intel/node_manager/airflow/avg | uint16 | Average Volumetric Airflow 
/intel/node_manager/airflow/max | uint16 | Maximal Volumetric Airflow 
/intel/node_manager/airflow/min | uint16 | Minimal Volumetric Airflow 
/intel/node_manager/cups/cpu_cstate | uint16 | CUPS CPU Bandwidth
/intel/node_manager/cups/io_bandwith | uint16 | CUPS I/O Bandwidth
/intel/node_manager/cups/memory_bandwith | uint16 | CUPS Memory Bandwidth
/intel/node_manager/power/cpu | uint16 | Current CPU power consumption
/intel/node_manager/power/cpu/avg | uint16 | Average CPU power consumption
/intel/node_manager/power/cpu/max | uint16 | Maximal CPU power consumption
/intel/node_manager/power/cpu/min | uint16 | Minimal CPU power consumption
/intel/node_manager/power/policy/power_limit | uint16 | Power policy
/intel/node_manager/margin/cpu/tj  | uint16 | Margin-to-throttle functional  (CPU)
/intel/node_manager/margin/cpu/tj/margin_offset | uint16 | Margin-to-spec reliability (CPU)
/intel/node_manager/power/memory | uint16 | Current Memory power consumption
/intel/node_manager/power/memory/avg | uint16 | Average Memory power consumption
/intel/node_manager/power/memory/max | uint16 | Maximal Memory power consumption
/intel/node_manager/power/memory/min | uint16 | Minimal Memory power consumption
/intel/node_manager/power/system | uint16 | Current Platform power consumption
/intel/node_manager/power/system/avg | uint16 | Average Platform power consumption
/intel/node_manager/power/system/max | uint16 | Maximal Platform power consumption
/intel/node_manager/power/system/min | uint16 | Minimal Platform power consumption
/intel/node_manager/temperature/cpu/cpu/<cpu_id> | uint16 | Current CPU temperature
/intel/node_manager/temperature/pmbus/VR/<VR_id> | uint16 | Current VR's temperature
/intel/node_manager/temperature/memory/dimm/<dimm_id> | uint16 | Current Memory dimms temperature
/intel/node_manager/temperature/outlet | uint16 | Current Outlet (exhaust air) temperature
/intel/node_manager/temperature/outlet/avg | uint16 | Average Outlet (exhaust air) temperature
/intel/node_manager/temperature/outlet/max | uint16 | Maximal Outlet (exhaust air) temperature
/intel/node_manager/temperature/outlet/min | uint16 | Minimal Outlet (exhaust air) temperature
/intel/node_manager/temperature/inlet | uint16 | Current Inlet Temperature
/intel/node_manager/temperature/inlet/avg | uint16 | Average Inlet Temperature
/intel/node_manager/temperature/inlet/max | uint16 | Maximal Inlet Temperature
/intel/node_manager/temperature/inlet/min | uint16 | Minimal Inlet Temperature

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
                "/intel/node_manager/airflow": {},
                "/intel/node_manager/airflow/avg": {},
                "/intel/node_manager/airflow/max": {},
                "/intel/node_manager/airflow/min": {},
                "/intel/node_manager/cups/cpu_cstate": {},
                "/intel/node_manager/cups/io_bandwith": {},
                "/intel/node_manager/cups/memory_bandwith": {},
                "/intel/node_manager/margin/cpu/tj": {},
                "/intel/node_manager/margin/cpu/tj/margin_offset": {},
                "/intel/node_manager/power/cpu": {},
                "/intel/node_manager/power/cpu/avg": {},
                "/intel/node_manager/power/cpu/max": {},
                "/intel/node_manager/power/cpu/min": {},
                "/intel/node_manager/power/memory": {},
                "/intel/node_manager/power/memory/avg": {},
                "/intel/node_manager/power/memory/max": {},
                "/intel/node_manager/power/memory/min": {},
                "/intel/node_manager/power/system": {},
                "/intel/node_manager/power/system/avg": {},
                "/intel/node_manager/power/system/max": {},
                "/intel/node_manager/power/system/min": {},
                "/intel/node_manager/temperature/cpu/cpu/0": {},
                "/intel/node_manager/temperature/cpu/cpu/1": {},
                "/intel/node_manager/temperature/memory/dimm/0": {},
                "/intel/node_manager/temperature/memory/dimm/1": {},
                "/intel/node_manager/temperature/memory/dimm/10": {},
                "/intel/node_manager/temperature/memory/dimm/11": {},
                "/intel/node_manager/temperature/memory/dimm/12": {},
                "/intel/node_manager/temperature/memory/dimm/13": {},
                "/intel/node_manager/temperature/memory/dimm/14": {},
                "/intel/node_manager/temperature/memory/dimm/15": {},
                "/intel/node_manager/temperature/memory/dimm/16": {},
                "/intel/node_manager/temperature/memory/dimm/17": {},
                "/intel/node_manager/temperature/memory/dimm/18": {},
                "/intel/node_manager/temperature/memory/dimm/19": {},
                "/intel/node_manager/temperature/memory/dimm/2": {},
                "/intel/node_manager/temperature/memory/dimm/20": {},
                "/intel/node_manager/temperature/memory/dimm/21": {},
                "/intel/node_manager/temperature/memory/dimm/22": {},
                "/intel/node_manager/temperature/memory/dimm/23": {},
                "/intel/node_manager/temperature/memory/dimm/24": {},
                "/intel/node_manager/temperature/memory/dimm/25": {},
                "/intel/node_manager/temperature/memory/dimm/26": {},
                "/intel/node_manager/temperature/memory/dimm/27": {},
                "/intel/node_manager/temperature/memory/dimm/28": {},
                "/intel/node_manager/temperature/memory/dimm/29": {},
                "/intel/node_manager/temperature/memory/dimm/3": {},
                "/intel/node_manager/temperature/memory/dimm/30": {},
                "/intel/node_manager/temperature/memory/dimm/31": {},
                "/intel/node_manager/temperature/memory/dimm/4": {},
                "/intel/node_manager/temperature/memory/dimm/5": {},
                "/intel/node_manager/temperature/memory/dimm/6": {},
                "/intel/node_manager/temperature/memory/dimm/7": {},
                "/intel/node_manager/temperature/memory/dimm/8": {},
                "/intel/node_manager/temperature/memory/dimm/9": {},
                "/intel/node_manager/temperature/outlet": {},
                "/intel/node_manager/temperature/outlet/avg": {},
                "/intel/node_manager/temperature/outlet/max": {},
                "/intel/node_manager/temperature/outlet/min": {}
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


## Community Support
This repository is one of **many** plugins in **Snap**, a powerful telemetry framework. See the full project at http://github.com/intelsdi-x/snap To reach out to other users, head to the [main framework](https://github.com/intelsdi-x/snap#community-support)

## Contributing
We love contributions!

There's more than one way to give back, from examples to blogs to code updates. See our recommended process in [CONTRIBUTING.md](CONTRIBUTING.md).

## License
[Snap](http://github.com:intelsdi-x/snap), along with this plugin, is an Open Source software released under the Apache 2.0 [License](LICENSE).

## Acknowledgements

* Author: [Lukasz Mroz](https://github.com/lmroz)
* Author: [Marcin Krolik](https://github.com/marcin-krolik)
* Author: [Patryk Matyjasek](https://github.com/PatrykMatyjasek)

And **thank you!** Your contribution, through code and participation, is incredibly important to us.
