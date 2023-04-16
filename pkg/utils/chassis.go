package utils

//system status exporter
//rest/system/historicalstatistics/1

import (
	"encoding/xml"
	"log"
	_ "github.com/mattn/go-sqlite3" // Import go-sqlite3 library
	"database/sql"
	"edge_exporter/pkg/database"
	"edge_exporter/pkg/http"
)

type ChassisData struct {
	XMLname    xml.Name   `xml:"root"`
	Chassis chassis       `xml:"chassis"`
}

type chassis struct {
	Rt_Chassis_Type   string `xml:"rt_Chassis_Type"`
	SerialNumber      string `xml:"SerialNumber"`    // Average percent usage of the CPU.
}


//Collect implements required collect function for all promehteus collectors
func GetChassisLabels(ipaddress string, phpsessid string) (chassisType string, serialNumber string, err error){
	//hosts := config.GetAllHosts()//retrieving targets for this exporter

	//log.Print(hosts)
	var sqliteDatabase *sql.DB
	//var labels *database.Chassis
	sqliteDatabase, err = sql.Open("sqlite3", "./sqlite-database.db")
	if err != nil {
		log.Print(err)
		return "","", err
	} // Open the created SQLite File
	// Defer Closing the database
	defer sqliteDatabase.Close()
	if (database.RowExists(sqliteDatabase, ipaddress)) {
		chassisType, serialNumber, err = database.GetChassis(sqliteDatabase, ipaddress)
		if (chassisType == "" || serialNumber == "" || err != nil) {

				dataStr := "https://"+ipaddress+"/rest/chassis"
				_, data,err := http.GetAPIData(dataStr, phpsessid)
				if err != nil {
						//log.Print("Error collecting from : ", err)
					return "http error","http error",err
				}
				b := []byte(data) //Converting string of data to bytestream
				ssbc := &ChassisData{}
				err = xml.Unmarshal(b, &ssbc) //Converting XML data to variables
				if err != nil {
				return "http error","http error",err
				}
				//log.Print("Successful API call data: ",ssbc.SystemData,"\n")

				chassisType := ssbc.Chassis.Rt_Chassis_Type
				serialNumber := ssbc.Chassis.SerialNumber

				err = database.InsertChassis(sqliteDatabase, ipaddress, chassisType, serialNumber)
					if err != nil {
						log.Print("insert chassis error", err)
					}
				return string(chassisType), string(serialNumber), err
		}
	}
return chassisType, serialNumber, err
}
/*
func test() {
	var sqliteDatabase *sql.DB
	//var labels *database.Chassis
	sqliteDatabase, err = sql.Open("sqlite3", "./sqlite-database.db")
	if err != nil {
		log.Print(err)
	}
	phpsessid,err := http.APISessionAuth("student", "PanneKake", "10.233.234.11")
	dataStr := "https://10.233.234.11/rest/chassis"
				_, data,err := http.GetAPIData(dataStr, phpsessid)
				if err != nil {
						log.Print("Error collecting from : ", err)

				}
				b := []byte(data) //Converting string of data to bytestream
				ssbc := &ChassisData{}
				xml.Unmarshal(b, &ssbc) //Converting XML data to variables
				//log.Print("Successful API call data: ",ssbc.SystemData,"\n")

				chassisType := ssbc.Chassis.Rt_Chassis_Type
				serialNumber := ssbc.Chassis.SerialNumber
				err = database.InsertChassis(sqliteDatabase, ipaddress, chassisType, serialNumber)
					if err != nil {
						log.Print("insert chassis error", err)
					}
				return chassisType, serialNumber, err
}
*/