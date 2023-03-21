package main

import (
	"crypto/tls"
	"database/sql"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
	_ "github.com/mattn/go-sqlite3"
)
type Cookie struct {
	Ipaddress string    `json:"ipaddress"`
	Phpsessid string    `json:"phpsessid"`
	Time      time.Time `json:"time"`
}
type Cookies struct {
	Name string `json:"Name"`
	Cookies []Cookie `json:"cookies"`
}
// The functions APISessionAuth(...) and getAPIData(...) utilizes curl-to-go translator but is modified for cookie management.
// Generated by curl-to-Go: https://mholt.github.io/curl-to-go

// curl -k --data "Username=student&Password=PanneKake23" -i -v https://10.233.230.11/rest/login

// TODO: This is insecure; use only in dev environments.
func APISessionAuth(username string, password string, ipaddress string) (string,error) {
	//var read []byte
	var phpsessid string

	cfg := getConf(&Config{})
	timeout := cfg.Authtimeout
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	client := &http.Client{Transport: tr,Timeout: time.Duration(timeout) * time.Second}

	params := url.Values{}
	params.Add("Username", username)
	params.Add("Password", password)
	body := strings.NewReader(params.Encode())

	req, err := http.NewRequest("POST", "https://"+ipaddress+"/rest/login", body)
	if err != nil {
		log.Flags()
			fmt.Println("error in auth:", err)
			return "Error fetching data", err
		//	fmt.Println("error in systemExporter:", error)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := client.Do(req)
	if err != nil {
		log.Flags()
		fmt.Println("error in auth:", err)
		return "Error fetching data", err
		//fmt.Println("error in systemExporter:", err)
	}

	  m := make(map[string]string)
	  for _, c := range resp.Cookies() {
		 m[c.Name] = c.Value
	  }
	phpsessid = m["PHPSESSID"]


	defer resp.Body.Close()



	return phpsessid,err
	}


func getAPIData(url string, phpsessid string) (string,error){

tr2 := &http.Transport{
	TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
}
client2 := &http.Client{Transport: tr2}
cookie1 := &http.Cookie{
	Name:   "PHPSESSID",
	Value:  phpsessid,
	//Path:     "/",
	MaxAge:   3600,
	HttpOnly: false,
	Secure:   true,
}
req2, err := http.NewRequest("GET", url, nil)
if err != nil {
	log.Flags()
		fmt.Println("error in getapidata():", err)
		return "Error fetching data", err
	//	fmt.Println("error in systemExporter:", error)
}
req2.AddCookie(cookie1)
	resp2, err := client2.Do(req2)
	if err != nil {
		log.Flags()
			fmt.Println("error in getapidata():", err)
			return "Error fetching data", err
	}

	b, err := ioutil.ReadAll(resp2.Body)
	defer resp2.Body.Close()

	return string(b), err
}


func main() {

	php, err  := APISessionAuth("student", "PanneKake23", "10.233.234.11")
	//php2, err  := APISessionAuth("student", "PanneKake23", "10.233.230.11")
	os.Remove("sqlite-database.db") // I delete the file to avoid duplicated records.
	// SQLite is a file based database.

	//log.Println("Creating sqlite-database.db...")
	file, err := os.Create("sqlite-database.db") // Create SQLite file
	if err != nil {
		log.Fatal(err.Error())
	}
	fmt.Println("sqlite-database.db created")
//sql.Open()
	var sqliteDatabase *sql.DB

	sqliteDatabase, _ = sql.Open("sqlite3", "./sqlite-database.db") // Open the created SQLite File
	 // Defer Closing the database
	createTable(sqliteDatabase) // Create Database Tables
	// INSERT RECORDS
	insertAuth(sqliteDatabase, "ipaddress", "phpsessid", time.Now().String())

	// DISPLAY INSERTED RECORDS
	//fmt.Println(ip, sess, time)
	defer sqliteDatabase.Close()
	defer file.Close()
	fmt.Println(php,err)


	//fmt.Println(php,err)
}
