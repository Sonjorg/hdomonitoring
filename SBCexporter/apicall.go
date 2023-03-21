package main

import (
	"crypto/tls"
	"fmt"
	"strings"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"
	"log"
	"encoding/json"
	"os"
)
type Cookie struct {
	Ipaddress string    `json:"ipaddress"`
	Phpsessid string    `json:"phpsessid"`
	Time      time.Time `json:"time"`
}
type Cookies struct {
	Name string `json:"Name"`
	Cookie Cookie `json:"cookies"`
}
// The functions APISessionAuth(...) and getAPIData(...) utilizes curl-to-go translator but is modified for cookie management.
// Generated by curl-to-Go: https://mholt.github.io/curl-to-go

// curl -k --data "Username=student&Password=PanneKake23" -i -v https://10.233.230.11/rest/login

// TODO: This is insecure; use only in dev environments.
func APISessionAuth(username string, password string, ipaddress string) (string,error) {
	//var read []byte
	var phpsessid string
	 /*
	read, err := ioutil.ReadFile("tmp.json")
	if err != nil {
		//struct := &Host{}
		//var str Name
		//doc := make(map[string]Host{})
		Hosts := []Cookie{}
		err := json.Unmarshal(read, &Hosts)
		if err != nil {
			fmt.Println("No data retrieved unmarhalling json phpsessid")
		}
		for i := range Hosts {
			if (Hosts[i].Ipaddress == ipaddress) {
				if (Hosts[i].Time.Add(8 * time.Minute).Before(time.Now())) {

					return Hosts[i].Phpsessid,nil
					//phpsessid = Hosts[i].Phpsessid
						fmt.Println("retrieved from file")
				}
			}
		}
	}*/
	var file []byte
	//var phpsessid string
	file, err := ioutil.ReadFile("data.json")
	if err == nil {
		//struct := &Host{}
		//var str Name
		//doc := make(map[string]Host{})
		//var Hosts = &Cookie{}

   // data := []Cookie{}

    // Here the magic happens!
    //json.Unmarshal(file, &data)
		Hosts := []Cookies{}
		err := json.Unmarshal(file, &Hosts)
		if err != nil {
			fmt.Println("No data retrieved unmarhalling json phpsessid",err)
		}
		fmt.Println(Hosts)
		for i := range Hosts {
			if (Hosts[i].Cookie.Ipaddress == "10.233.234.11") {
				//if (Hosts[i].Time.Add(2 * time.Minute).Before(time.Now())) {

					//Hosts[i].Phpsessid
					//phpsessid = Hosts[i].Phpsessid
						fmt.Println(Hosts[i].Cookie.Ipaddress)
						return Hosts[i].Cookie.Ipaddress, nil
			//	}
			}
		}
	} else {fmt.Println("cant open")}
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
	 // fmt.Println(m["PHPSESSID"])
	phpsessid = m["PHPSESSID"]
//d := Cookies{}
	//data := Cookie{ipaddress, phpsessid, time.Now()}
    //data := Cookies{}
	x := Cookie{
		Ipaddress: ipaddress,
		Phpsessid: phpsessid,
		Time:      time.Now(),
	}
	data := Cookies{
		"test",
		x,
    }
	//data = append(data, *c)
	dataBytes, err := json.Marshal(data)
    if err != nil {
        fmt.Println(err)
    }
	dataBytes, _ = json.MarshalIndent(data, "", "  ")
/*
    err = ioutil.WriteFile("data.json", dataBytes, 0644)
    if err != nil {
        fmt.Println(err)
    }*/
	f, err := os.OpenFile("./data.json", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		fmt.Println(err)
	}
	defer f.Close()

	n, err := f.Write(dataBytes)
	if err != nil {
		fmt.Println(n, err)
	}

	if n, err = f.WriteString("\n"); err != nil {
		fmt.Println(n, err)
	}
	//jsonByte, _ := json.Marshal(data)
	//jsonByte, _ = json.MarshalIndent(data, "", "  ")
	/*
	err := os.OpenFile("./data.json", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		fmt.Println(err)

	}
	defer f.Close()
	n, err := f.Write(jsonByte)
		if err != nil {
			fmt.Println(n, err)
		}

		if n, err = f.WriteString("\n"); err != nil {
			fmt.Println(n, err)
		}*/
	/*if n, err = f.WriteString("\n"); err != nil {
		fmt.Println(n, err)
	}*/
/*

	err = ioutil.WriteFile("data.json", jsonByte, 0644)
	if err != nil {
	  fmt.Println(err)
 	 }
*/
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
	php2, err  := APISessionAuth("student", "PanneKake23", "10.233.230.11")

	fmt.Println(php,php2,err)


	//fmt.Println(php,err)
}
