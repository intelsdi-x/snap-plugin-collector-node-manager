##Pulse IPMI Collector plugin

#Description
Plugin collects platform's metrics using in-band IPMI device compatible with OpenIPMI driver protocol.

#Assumptions
* Linux kernel version >= 2.6.12-rc2
* IPMI module loaded

#Metrics
 - /intel
	 - /ipmi
		 - /cups
			 - /cpu_cstate - cpu utilization [%]
			 - /memory_bandwidth - mem bandwidth [%]
			 - /io_bandwidth - io bandwidth [%]
		 - /power
			 - /system - total platform power utilization [W]
				 - /min - minimal value
				 - /max - maximal value
				 - /avg - average value
			 - /cpu - cpu power utilization [W]
				 - /min - minimal value
				 - /max - maximal value
				 - /avg - average value
			 - /memory - memory power utilization [W]
				 - /min - minimal value
				 - /max - maximal value
				 - /avg - average value
		 - /temperature
			 - /inlet - inlet air temperature  [℃]
				 - /min - minimal value
				 - /max - maximal value
				 - /avg - average value
			 - /outlet - outlet air temperature  [℃]
				 - /min - minimal value
				 - /max - maximal value
				 - /avg - average value
			 - /pmbus
				 - VR[0..7] - temperature of given voltage regulator
			 - /chipset - chipset temperature  [℃]
			 - /cpu[0..3] - temperature of given CPU [℃]
			 - /memory
				 - /dimm[0..63] - temperature of given memory module [℃]
		 - /airflow - global volumetric airflow statistics [dCFM]
			 - /min - minimal value
			 - /max - maximal value
			 - /avg - average value
		 - /margin
			 - /cpu
				 - /tj - current temperature limit to cpu throttling
					 - /offset_margin - current offset value
