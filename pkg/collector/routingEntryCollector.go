// routingentry
package collector

import (
	"edge_exporter/pkg/config"
	"edge_exporter/pkg/database"
	"edge_exporter/pkg/http"
	"encoding/xml"
	"fmt"
	//"sync"
	//"log"
	"regexp"
	"github.com/prometheus/client_golang/prometheus"
	//"strconv"
	//"time"
	//"exporter/sqlite"
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
)

// rest/routingtable/2/routingentry
// first request
// rest/routingtable/
type routingTables struct {
	// Value  float32 `xml:",chardata"`
	XMLName        xml.Name       `xml:"root"`
	RoutingTables2 routingTables2 `xml:"routingtable_list"`
}
type routingTables2 struct {
	RoutingTables3 routingTables3 `xml:"routingtable_pk"`
}
type routingTables3 struct {
	Attr  []string `xml:"id,attr"`
	Value string   `xml:",chardata"`
}

// Second request
// rest/routingtable/4/routingentry
type routingEntries struct {
	XMLName       xml.Name      `xml:"root"`
	RoutingEntry2 routingEntry2 `xml:"routingentry_list"`
}
type routingEntry2 struct {
	RoutingEntry3 routingEntry3 `xml:"routingentry_pk"`
}
type routingEntry3 struct {
	Attr  []string `xml:"id,attr"`
	Value string   `xml:",chardata"`
}

// Third request
// rest/routingtable/2/routingentry/1/historicalstatistics/1
type rSBCdata struct {
	XMLname     xml.Name    `xml:"root"`
	Status      rStatus     `xml:"status"`
	RoutingData routingData `xml:"historicalstatistics"`
}
type rStatus struct {
	HTTPcode string `xml:"http_code"`
}
type routingData struct {
	Href              string `xml:"href,attr"`
	Rt_RuleUsage      int    `xml:"rt_RuleUsage"`
	Rt_ASR            int    `xml:"rt_ASR"`
	Rt_RoundTripDelay int    `xml:"rt_RoundTripDelay"`
	Rt_Jitter         int    `xml:"rt_Jitter"`
	Rt_MOS            int    `xml:"rt_MOS"`
	Rt_QualityFailed  int    `xml:"rt_QualityFailed"`
}

// Metrics for each routingentry
type rMetrics struct {
	Rt_RuleUsage      *prometheus.Desc
	Rt_ASR            *prometheus.Desc
	Rt_RoundTripDelay *prometheus.Desc
	Rt_Jitter         *prometheus.Desc
	Rt_MOS            *prometheus.Desc
	Rt_QualityFailed  *prometheus.Desc
	Error_ip          *prometheus.Desc
}

func routingCollector() *rMetrics {

	return &rMetrics{
		Rt_RuleUsage: prometheus.NewDesc("rt_RuleUsage",
			"NoDescriptionYet",
			[]string{"Instance", "hostname", "job", "routing_table", "routing_entry", "HTTP_status"}, nil,
		),
		Rt_ASR: prometheus.NewDesc("rt_ASR",
			"NoDescriptionYet",
			[]string{"Instance", "hostname", "job", "routing_table", "routing_entry", "HTTP_status"}, nil,
		),
		Rt_RoundTripDelay: prometheus.NewDesc("rt_RoundTripDelay",
			"NoDescriptionYet",
			[]string{"Instance", "hostname", "job", "routing_table", "routing_entry", "HTTP_status"}, nil,
		),
		Rt_Jitter: prometheus.NewDesc("rt_Jitter",
			"NoDescriptionYet",
			[]string{"Instance", "hostname", "job", "routing_table", "routing_entry", "HTTP_status"}, nil,
		),
		Rt_MOS: prometheus.NewDesc("rt_MOS",
			"NoDescriptionYet.",
			[]string{"Instance", "hostname", "job", "routing_table", "routing_entry", "HTTP_status"}, nil,
		),
		Rt_QualityFailed: prometheus.NewDesc("rt_QualityFailed",
			"NoDescriptionYet",
			[]string{"Instance", "hostname", "job", "routing_table", "routing_entry", "HTTP_status"}, nil,
		),
		Error_ip: prometheus.NewDesc("error_edge_routing",
			"NoDescriptionYet",
			[]string{"Instance", "hostname"}, nil,
		),
	}
}

// Each and every collector must implement the Describe function.
// It essentially writes all descriptors to the prometheus desc channel.
func (collector *rMetrics) Describe(ch chan<- *prometheus.Desc) {
	//Update this section with the each metric you create for a given collector
	ch <- collector.Rt_RuleUsage
	ch <- collector.Rt_ASR
	ch <- collector.Rt_RoundTripDelay
	ch <- collector.Rt_Jitter
	ch <- collector.Rt_MOS
	ch <- collector.Rt_QualityFailed
	ch <- collector.Error_ip
}

// Collect implements required collect function for all promehteus collectors
func (collector *rMetrics) Collect(c chan<- prometheus.Metric) {
	hosts := config.GetIncludedHosts("routingentry") //retrieving targets for this exporter
	if len(hosts) <= 0 {
		fmt.Println("no hosts")
		return
	}
	var metricValue1 float64
	var metricValue2 float64
	var metricValue3 float64
	var metricValue4 float64
	var metricValue5 float64
	var metricValue6 float64

	var sqliteDatabase *sql.DB
	sqliteDatabase, err := sql.Open("sqlite3", "./sqlite-database.db")
	if err != nil {
		fmt.Println(err)
	}
	var routingtables []string
	var match []string //variable to hold routingentries cleaned with regex

	for i := range hosts {

		phpsessid, err := http.APISessionAuth(hosts[i].Username, hosts[i].Password, hosts[i].Ip)
			if err != nil {
				fmt.Println("Error auth", hosts[i].Ip, err)
				//continue
				return
			}
			//var match []string

			if (database.RoutingTablesExists(sqliteDatabase,hosts[i].Ip)) { //fetching from database
					routingtables, err = database.GetRoutingTables(sqliteDatabase,hosts[i].Ip)
						if err != nil {
							fmt.Println(err)
						}
			} else { //fetching from router
				_, data, err := http.GetAPIData("https://"+hosts[i].Ip+"/rest/routingtable", phpsessid)
				if err != nil {
					fmt.Println("Error routingtable data", hosts[i].Ip, err)
					//continue
					return
				}
				rt := &routingTables{}
				xml.Unmarshal(data, &rt) //Converting XML data to variables
				routingtables = rt.RoutingTables2.RoutingTables3.Attr //ssbc.Rt2.Rt3.Attr
				if len(routingtables) <= 0 {
					fmt.Println("Routingtables empty")
					return
				}
				err = database.CreateRoutingSqlite(sqliteDatabase)
					if err != nil {
						fmt.Println(err)
					}
				database.StoreRoutingTables(sqliteDatabase, hosts[i].Ip, "test", routingtables)
				}
			for j := range routingtables {

				//Trying to fetch routingentries from database, if not exist yet, fetch new ones
				if (database.RoutingEntriesExists(sqliteDatabase,hosts[i].Ip)				) {
					match, err = database.GetRoutingEntries(sqliteDatabase,hosts[i].Ip,routingtables[j])
						if err != nil {
							fmt.Println(err)
						}
				} else {
					url := "https://" + hosts[i].Ip + "/rest/routingtable/" + routingtables[j] + "/routingentry"
					_, data2, err := http.GetAPIData(url, phpsessid)
					if err != nil {
					}
					//b2 := []byte(data2) //Converting string of data to bytestream
					re := &routingEntries{}
					xml.Unmarshal(data2, &re) //Converting XML data to variables
					routingEntries := re.RoutingEntry2.RoutingEntry3.Attr
					if len(routingEntries) <= 0 {
						fmt.Println("No routingEntry for this routingtable")
						continue
					}
					//
					entries := regexp.MustCompile(`\d+$`)


					for k := range routingEntries {
						tmp := entries.FindStringSubmatch(routingEntries[k])
						for l := range tmp {
							match = append(match, tmp[l])
							//fmt.Println(tmp[l])
						}
					}
					//Storing fetched routingentries

					err = database.StoreRoutingEntries(sqliteDatabase, hosts[i].Ip, routingtables[j], match)
					if err != nil {
						fmt.Println(err)
					}
				}

				for k := range match {

					url := "https://" + hosts[i].Ip + "/rest/routingtable/" + routingtables[j] + "/routingentry/" + match[k] + "/historicalstatistics/1"
					_, data3, err := http.GetAPIData(url, phpsessid)
					if err != nil {
						fmt.Println(err)

						continue
					}

					rData := &rSBCdata{}
					xml.Unmarshal(data3, &rData) //Converting XML data to variables
					//fmt.Println("Successful API call data: ",rData.RoutingData,"\n")

					metricValue1 = float64(rData.RoutingData.Rt_RuleUsage)
					metricValue2 = float64(rData.RoutingData.Rt_ASR)
					metricValue3 = float64(rData.RoutingData.Rt_RoundTripDelay)
					metricValue4 = float64(rData.RoutingData.Rt_Jitter)
					metricValue5 = float64(rData.RoutingData.Rt_MOS)
					metricValue6 = float64(rData.RoutingData.Rt_QualityFailed)

					c <- prometheus.MustNewConstMetric(collector.Rt_RuleUsage, prometheus.GaugeValue, metricValue1, hosts[i].Ip, hosts[i].Hostname, "routingentry", routingtables[j], match[k], "test")
					c <- prometheus.MustNewConstMetric(collector.Rt_ASR, prometheus.GaugeValue, metricValue2, hosts[i].Ip, hosts[i].Hostname, "routingentry", routingtables[j], match[k], "test")
					c <- prometheus.MustNewConstMetric(collector.Rt_RoundTripDelay, prometheus.GaugeValue, metricValue3, hosts[i].Ip, hosts[i].Hostname, "routingentry", routingtables[j], match[k], "test")
					c <- prometheus.MustNewConstMetric(collector.Rt_Jitter, prometheus.GaugeValue, metricValue4, hosts[i].Ip, hosts[i].Hostname, "routingentry", routingtables[j], match[k], "test")
					c <- prometheus.MustNewConstMetric(collector.Rt_MOS, prometheus.GaugeValue, metricValue5, hosts[i].Ip, hosts[i].Hostname, "routingentry", routingtables[j], match[k], "test")
					c <- prometheus.MustNewConstMetric(collector.Rt_QualityFailed, prometheus.GaugeValue, metricValue6, hosts[i].Ip, hosts[i].Hostname, "routingentry", routingtables[j], match[k], "test")
				}
			}
	}
}

func RoutingTablesExists(sqliteDatabase *sql.DB, s string) {
	panic("unimplemented")
}

func RoutingEntryCollector() {
	c := routingCollector()
	prometheus.MustRegister(c)
}