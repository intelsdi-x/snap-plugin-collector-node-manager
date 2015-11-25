# Snap Intel Node Manager Collector Plugin

 Plugin to collect data from Intel's Node Manager. Which is presenting low level metrics like power consumption, cpu temperature, etc.
 Currently it is using IPMI device to collect data from NM.

1. [Getting Started](#getting-started)
  * [System Requirements](#system-requirements)
  * [Configuration and Usage](configuration-and-usage)
2. [Documentation](#documentation)
  * [Collected Metrics](#collected-metrics)
  * [Examples](#examples)
  * [Roadmap](#roadmap)
3. [Community Support](#community-support)
4. [Contributing](#contributing)
5. [License](#license-and-authors)
6. [Acknowledgements](#acknowledgements)

## Getting Started

 Plugin collects specified metrics in-band on OS level

### System Requirements

Include:

 - Plugin needs to be run on server platform which supports Intel Node Manager.
 - Currently it works only on Linux Servers


### Configuration and Usage

 On OS level user needs to load modules:
  - ipmi_msghandler
  - ipmi_devintf
  - ipmi_si
 Those modules provides specific IPMI device which can collect data from NM

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

### Roadmap
As we launch this plugin, we have a few items in mind for the next release:
- Out-Of-Band Support
- Ipmitool support

## Community Support
This repository is one of **many** plugins in the **Snap Framework**: a powerful telemetry agent framework. To reach out on other use cases, visit:

* Snap Gitter channel (@TODO Link)
* Our Google Group (@TODO Link)

The full project is at http://github.com:intelsdi-x/snap.

## Contributing
We love contributions! :heart_eyes:

There's more than one way to give back, from examples to blogs to code updates. See our recommended process in [CONTRIBUTING.md](CONTRIBUTING.md).

## License
Snap, along with this plugin, is an Open Source software released under the Apache 2.0 [License](LICENSE).

## Acknowledgements
List authors, co-authors and anyone you'd like to mention

* Author: [Lukasz Mroz](https://github.com/lmroz)
* Author: [Marcin Krolik](https://github.com/marcin-krolik)
* Author: [Patryk Matyjasek](https://github.com/PatrykMatyjasek)

**Thank you!** Your contribution is incredibly important to us.
