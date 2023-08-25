/* Copyright (C) 2023 Sondre Jørgensen - All Rights Reserved
 * You may use, distribute and modify this code under the
 * terms of the CC BY 4.0 license
 */
package collector

import (
	"edge_exporter/pkg/config"
	"edge_exporter/pkg/http"
	"edge_exporter/pkg/utils"
	"encoding/xml"
	"log"
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

type sSBCdata struct {
	XMLname    xml.Name   `xml:"root"`
	Status     sStatus    `xml:"status"`
	SystemData systemData `xml:"historicalstatistics"`
}
type sStatus struct {
	HTTPcode string `xml:"http_code"`
}
type systemData struct {
	Href                 string `xml:"href,attr"`
	Rt_CPUUsage          int    `xml:"rt_CPUUsage"`    // Average percent usage of the CPU.
	Rt_MemoryUsage       int    `xml:"rt_MemoryUsage"` // Average percent usage of system memory. int
	Rt_CPUUptime         int    `xml:"rt_CPUUptime"`
	Rt_FDUsage           int    `xml:"rt_FDUsage"`
	Rt_CPULoadAverage1m  int    `xml:"rt_CPULoadAverage1m"`
	Rt_CPULoadAverage5m  int    `xml:"rt_CPULoadAverage5m"`
	Rt_CPULoadAverage15m int    `xml:"rt_CPULoadAverage15m"`
	Rt_TmpPartUsage      int    `xml:"rt_TmpPartUsage"` //Percentage of the temporary partition used. int
	Rt_LoggingPartUsage  int    `xml:"rt_LoggingPartUsage"`
}

func SystemCollector(host *config.HostCompose) (m []prometheus.Metric, successfulScrape bool) {

	var (
		Rt_CPUUsage = prometheus.NewDesc("edge_system_CPUUsage",
			"Average percent usage of the CPU.",
			[]string{"hostip", "hostname", "chassis_type", "serial_number"}, nil,
		)
		Rt_MemoryUsage = prometheus.NewDesc("edge_system_MemoryUsage",
			"Average percent usage of system memory.",
			[]string{"hostip", "hostname", "chassis_type", "serial_number"}, nil,
		)
		Rt_CPUUptime = prometheus.NewDesc("edge_system_CPUUptime",
			"The total duration in seconds, that the system CPU has been UP and running.",
			[]string{"hostip", "hostname", "chassis_type", "serial_number"}, nil,
		)
		Rt_FDUsage = prometheus.NewDesc("edge_disk_FDUsage",
			"Number of file descriptors used by the system.",
			[]string{"hostip", "hostname", "chassis_type", "serial_number"}, nil,
		)
		Rt_CPULoadAverage1m = prometheus.NewDesc("edge_system_CPULoadAverage1m",
			"Average number of processes over the last one minute waiting to run because CPU is busy.",
			[]string{"hostip", "hostname", "chassis_type", "serial_number"}, nil,
		)
		Rt_CPULoadAverage5m = prometheus.NewDesc("edge_system_CPULoadAverage5m",
			"Average number of processes over the last five minutes waiting to run because CPU is busy.",
			[]string{"hostip", "hostname", "chassis_type", "serial_number"}, nil,
		)
		Rt_CPULoadAverage15m = prometheus.NewDesc("edge_system_CPULoadAverage15m",
			"Average number of processes over the last fifteen minutes waiting to run because CPU is busy.",
			[]string{"hostip", "hostname", "chassis_type", "serial_number"}, nil,
		)
		Rt_TmpPartUsage = prometheus.NewDesc("edge_disk_TmpPartUsage",
			"Percentage of the temporary partition used.",
			[]string{"hostip", "hostname", "chassis_type", "serial_number"}, nil,
		)
		Rt_LoggingPartUsage = prometheus.NewDesc("edge_disk_LoggingPartUsage",
			"Percentage of the logging partition used. This is applicable only for the SBC2000.",
			[]string{"hostip", "hostname", "chassis_type", "serial_number"}, nil,
		)
		Error_ip = prometheus.NewDesc("edge_system_status",
			"Returns 1 if the SBC Edge scrape was successful, and 0 if not.",
			[]string{"hostip", "hostname"}, nil,
		)
	)

		dataStr := "https://" + host.Ip + "/rest/system/historicalstatistics/1"
		timeReportedByExternalSystem := time.Now()

		if (!http.SBCIsUp(host.Ip)){
			m = append(m, prometheus.NewMetricWithTimestamp(
				timeReportedByExternalSystem,
				prometheus.MustNewConstMetric(
					Error_ip, prometheus.GaugeValue, 0, host.Ip, host.Hostname),
			))
			return m, false
		}

		chassisType, serialNumber, err := utils.GetChassisLabels(host.Ip, "null")
		if err != nil {
			chassisType, serialNumber = "Error fetching chassisinfo", "Error fetching chassisinfo"
			log.Print(err)
		}
		
		phpsessid, err := http.APISessionAuth(host.Username, host.Password, host.Ip)
		if err != nil {
			log.Println("Error retrieving session cookie (system): ", log.Flags(), err)
			m = append(m, prometheus.NewMetricWithTimestamp(
				timeReportedByExternalSystem,
				prometheus.MustNewConstMetric(
					Error_ip, prometheus.GaugeValue, 0, host.Ip, host.Hostname),
			))

			return m, false//trying next ip address
		}
		//fetching labels from DB or if not exist yet; from router
		if chassisType == "Error fetching chassisinfo" {
			chassisType, serialNumber, err = utils.GetChassisLabels(host.Ip, phpsessid)
			if err != nil {
				chassisType, serialNumber = "Error fetching chassisinfo", "Error fetching chassisinfo"
				log.Print(err)
			}
		}
		//Fetching systemdata
		_, data, err := http.GetAPIData(dataStr, phpsessid)
		if err != nil {
			log.Print("Error collecting from host: ", log.Flags(), err, "\n")
			m = append(m, prometheus.NewMetricWithTimestamp(
				timeReportedByExternalSystem,
				prometheus.MustNewConstMetric(
					Error_ip, prometheus.GaugeValue, 0, host.Ip, host.Hostname),
			))
			return m, false
		}
		ssbc := &sSBCdata{}
		err = xml.Unmarshal(data, &ssbc) //Converting XML data to variables
		if err != nil {
			log.Print("XML error system", err)
		}

		metricValue1 := float64(ssbc.SystemData.Rt_CPULoadAverage15m)
		metricValue2 := float64(ssbc.SystemData.Rt_CPULoadAverage1m)
		metricValue3 := float64(ssbc.SystemData.Rt_CPULoadAverage5m)
		metricValue4 := float64(ssbc.SystemData.Rt_CPUUptime)
		metricValue5 := float64(ssbc.SystemData.Rt_CPUUsage)
		metricValue6 := float64(ssbc.SystemData.Rt_FDUsage)
		metricValue7 := float64(ssbc.SystemData.Rt_LoggingPartUsage)
		metricValue8 := float64(ssbc.SystemData.Rt_MemoryUsage)
		metricValue9 := float64(ssbc.SystemData.Rt_TmpPartUsage)

		m = append(m, prometheus.MustNewConstMetric(
			Error_ip, prometheus.GaugeValue, 1, host.Ip, host.Hostname))

		m = append(m, prometheus.MustNewConstMetric(Rt_CPULoadAverage15m, prometheus.GaugeValue, metricValue1, host.Ip, host.Hostname, chassisType, serialNumber))
		m = append(m, prometheus.MustNewConstMetric(Rt_CPULoadAverage1m, prometheus.GaugeValue, metricValue2, host.Ip, host.Hostname, chassisType, serialNumber))
		m = append(m, prometheus.MustNewConstMetric(Rt_CPULoadAverage5m, prometheus.GaugeValue, metricValue3, host.Ip, host.Hostname, chassisType, serialNumber))
		m = append(m, prometheus.MustNewConstMetric(Rt_CPUUptime, prometheus.GaugeValue, metricValue4, host.Ip, host.Hostname, chassisType, serialNumber))
		m = append(m, prometheus.MustNewConstMetric(Rt_CPUUsage, prometheus.GaugeValue, metricValue5, host.Ip, host.Hostname, chassisType, serialNumber))
		m = append(m, prometheus.MustNewConstMetric(Rt_FDUsage, prometheus.GaugeValue, metricValue6, host.Ip, host.Hostname, chassisType, serialNumber))
		m = append(m, prometheus.MustNewConstMetric(Rt_LoggingPartUsage, prometheus.GaugeValue, metricValue7, host.Ip, host.Hostname, chassisType, serialNumber))
		m = append(m, prometheus.MustNewConstMetric(Rt_MemoryUsage, prometheus.GaugeValue, metricValue8, host.Ip, host.Hostname, chassisType, serialNumber))
		m = append(m, prometheus.MustNewConstMetric(Rt_TmpPartUsage, prometheus.GaugeValue, metricValue9, host.Ip, host.Hostname, chassisType, serialNumber))
	
	return m, true
}
