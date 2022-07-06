# Telegraf

## Docs

* Plugins used
  * <https://github.com/influxdata/telegraf/tree/master/plugins/inputs/redfish>
  * <https://github.com/influxdata/telegraf/tree/master/plugins/inputs/cpu>
  * <https://github.com/influxdata/telegraf/tree/master/plugins/inputs/mem>
  * <https://github.com/influxdata/telegraf/tree/master/plugins/inputs/net>
  * <https://github.com/influxdata/telegraf/tree/master/plugins/inputs/gnmi>
  * <https://github.com/influxdata/telegraf/tree/master/plugins/inputs/dpdk>
  * <https://github.com/influxdata/telegraf/tree/master/plugins/outputs/file>
  * <https://github.com/influxdata/telegraf/tree/master/plugins/outputs/influxdb_v2>
  * <https://github.com/influxdata/telegraf/tree/master/plugins/outputs/opentelemetry>

* Dockers used
  * <https://hub.docker.com/_/telegraf>
  * <https://hub.docker.com/_/influxdb>
  * <https://hub.docker.com/r/otel/opentelemetry-collector>

## Getting started

Run `docker-compose -f docker-compose.telegraf.yml up`

## Example

<!-- markdownlint-disable -->
```text
$ docker run --rm --net=host -v $(pwd)/telegraf-spdk.conf:/etc/telegraf/telegraf.conf:ro telegraf:1.22
2022-03-29T18:47:11Z I! Using config file: /etc/telegraf/telegraf.conf
2022-03-29T18:47:11Z I! Starting Telegraf 1.22.0
2022-03-29T18:47:11Z I! Loaded inputs: http
2022-03-29T18:47:11Z I! Loaded aggregators:
2022-03-29T18:47:11Z I! Loaded processors:
2022-03-29T18:47:11Z W! Outputs are not used in testing mode!
2022-03-29T18:47:11Z I! Tags enabled: host=localhost

mem,host=52ee5c75df01 commit_limit=69312983040i,committed_as=13494636544i,huge_page_size=2097152i,used_percent=10.100053796757296,high_free=0i,inactive=13541511168i,low_free=0i,shared=3904901120i,sreclaimable=812650496i,swap_cached=0i,free=100370612224i,huge_pages_total=2048i,low_total=0i,page_tables=49500160i,used=13567504384i,huge_pages_free=1357i,mapped=901996544i,slab=2243977216i,swap_total=4294963200i,cached=20385955840i,vmalloc_chunk=0i,vmalloc_used=0i,write_back=0i,swap_free=4294963200i,high_total=0i,available_percent=86.20598148102354,available=115801366528i,sunreclaim=1431326720i,total=134331011072i,buffered=6938624i,dirty=856064i,vmalloc_total=14073748835531776i,write_back_tmp=0i,active=8859537408i 1650954170000000000

net,host=52ee5c75df01,interface=eth0 drop_in=0i,drop_out=0i,bytes_sent=16589i,bytes_recv=13986i,packets_sent=89i,packets_recv=110i,err_in=0i,err_out=0i 1650954170000000000

cpu,cpu=cpu0,host=52ee5c75df01 usage_user=99.6999999973923,usage_system=0.09999999999763531,usage_idle=0,usage_iowait=0,usage_softirq=0,usage_steal=0,usage_nice=0,usage_irq=0.19999999999527063,usage_guest=0,usage_guest_nice=0 1650954170000000000
cpu,cpu=cpu1,host=52ee5c75df01 usage_user=99.70000000204891,usage_system=0,usage_irq=0.2999999999974534,usage_steal=0,usage_idle=0,usage_nice=0,usage_iowait=0,usage_softirq=0,usage_guest=0,usage_guest_nice=0 1650954170000000000
cpu,cpu=cpu2,host=52ee5c75df01 usage_system=0,usage_idle=0,usage_iowait=0,usage_guest_nice=0,usage_user=99.79999999981374,usage_nice=0,usage_irq=0.20000000000436557,usage_softirq=0,usage_steal=0,usage_guest=0 1650954170000000000
cpu,cpu=cpu3,host=52ee5c75df01 usage_guest_nice=0,usage_user=99.79999999981374,usage_idle=0,usage_nice=0,usage_iowait=0,usage_guest=0,usage_system=0,usage_irq=0.20000000000436557,usage_softirq=0,usage_steal=0 1650954170000000000
cpu,cpu=cpu4,host=52ee5c75df01 usage_user=99.70029970233988,usage_guest=0,usage_steal=0,usage_guest_nice=0,usage_system=0.09990009990223975,usage_idle=0,usage_nice=0,usage_iowait=0,usage_irq=0.19980019979993657,usage_softirq=0 1650954170000000000
cpu,cpu=cpu5,host=52ee5c75df01 usage_nice=0,usage_iowait=0,usage_irq=0.2997002997044478,usage_softirq=0,usage_steal=0,usage_guest_nice=0,usage_user=99.70029970233988,usage_idle=0,usage_guest=0,usage_system=0 1650954170000000000

telegraf_1  | redfish_thermal_temperatures,address=xpu-bmc,health=OK,host=220e68143c3d,member_id=0,name=CPU1\ Temp,rack=WEB43,row=North,source=web483,state=Enabled reading_celsius=41,upper_threshold_critical=45,upper_threshold_fatal=48 1656718190000000000
telegraf_1  | redfish_thermal_temperatures,address=xpu-bmc,host=220e68143c3d,member_id=1,name=CPU2\ Temp,rack=WEB43,row=North,source=web483,state=Disabled upper_threshold_critical=45,upper_threshold_fatal=48 1656718190000000000
telegraf_1  | redfish_thermal_temperatures,address=xpu-bmc,health=OK,host=220e68143c3d,member_id=2,name=Chassis\ Intake\ Temp,rack=WEB43,row=North,source=web483,state=Enabled upper_threshold_critical=40,upper_threshold_fatal=50,lower_threshold_critical=5,lower_threshold_fatal=0,reading_celsius=25 1656718190000000000
telegraf_1  | redfish_thermal_fans,address=xpu-bmc,health=OK,host=220e68143c3d,member_id=0,name=BaseBoard\ System\ Fan,rack=WEB43,row=North,source=web483,state=Enabled lower_threshold_fatal=0i,reading_rpm=2100i 1656718190000000000
telegraf_1  | redfish_thermal_fans,address=xpu-bmc,health=OK,host=220e68143c3d,member_id=1,name=BaseBoard\ System\ Fan\ Backup,rack=WEB43,row=North,source=web483,state=Enabled lower_threshold_fatal=0i,reading_rpm=2050i 1656718190000000000
telegraf_1  | redfish_power_powersupplies,address=xpu-bmc,health=Warning,host=220e68143c3d,member_id=0,name=Power\ Supply\ Bay,rack=WEB43,row=North,source=web483,state=Enabled line_input_voltage=120,last_power_output_watts=325,power_capacity_watts=800 1656718190000000000
telegraf_1  | redfish_power_voltages,address=xpu-bmc,health=OK,host=220e68143c3d,member_id=0,name=VRM1\ Voltage,rack=WEB43,row=North,source=web483,state=Enabled reading_volts=12,upper_threshold_critical=13,upper_threshold_fatal=15,lower_threshold_critical=11,lower_threshold_fatal=10 1656718190000000000
telegraf_1  | redfish_power_voltages,address=xpu-bmc,health=OK,host=220e68143c3d,member_id=1,name=VRM2\ Voltage,rack=WEB43,row=North,source=web483,state=Enabled reading_volts=5,upper_threshold_critical=7,lower_threshold_critical=4.5 1656718190000000000

...
```
<!-- markdownlint-restore -->
