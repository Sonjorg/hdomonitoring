package main

import (
	"fmt"
	//"os"
	"gopkg.in/yaml.v2"
    //"flag"
   // "log"
    "ioutil"
)
// Template used for struct and the functions NewConfig(), ValidateConfigPath() and ParseFlags() are copied from:
// https://dev.to/koddr/let-s-write-config-for-your-golang-web-app-on-right-way-yaml-5ggp
    type Config struct {
        Hosts []Host
    }
        type Host struct {
            HostName       string `yaml:"hostname"`
            Ipaddress      string `yaml:"ipaddress"`
            Username       string `yaml:"username"`
            Password       string `yaml:"password"`
            //exclude        string `yaml:"exclude"`
                Exclude struct {
                    // Server is the general server timeout to use
                    // for graceful shutdowns
                    SystemExporter bool `yaml:"systemstats"`
                    CallStats      bool `yaml:"callstats"`
                }`yaml:"exclude"`
            }


            func (c *Config) getConf() (*Config) {

                yamlFile, err := ioutil.ReadFile("config.yml")
                if err != nil {
                    //log.Printf("yamlFile.Get err   #%v ", err)
                    fmt.Println("yamlFile.Get err   #%v ", err)
                }
                err = yaml.Unmarshal(yamlFile, c)
                if err != nil {
                   // log.Fatalf("Unmarshal: %v", err)
                    fmt.Println("yamlFile.Get err   #%v ", err)
                }
                return c
            }
    // NewConfig returns a new decoded Config struct
 /*func getConfig() (*Config, error) {
        // Create config structure
        config := &Config{}
        // Open config file
        file, err := os.Open("./config.yml")
        if err != nil {
            return nil, err
        }
        defer file.Close()
        // Init new YAML decode
        d := yaml.NewDecoder(file)
        // Start YAML decoding from file
        if err := d.Decode(&config); err != nil {
            return nil, err
        }
        return config, nil
    }*/

   // test := NewConfig(.\config).
   // type hosts []hostConfig
   /*func readConfig() (*Config, error) {
    config := &Config{}
    cfgFile, err := ioutil.ReadFile("./config.yaml")
    if err != nil {
        return nil, err
    }
    err = yaml.Unmarshal(cfgFile, config)
    return &config.Config, err
}
   func ValidateConfigPath(path string) error {
    s, err := os.Stat(path)
    if err != nil {
        return err
    }
    if s.IsDir() {
        return fmt.Errorf("'%s' is a directory, not a normal file", path)
    }
    return nil
}*/
/*
// ParseFlags will create and parse the CLI flags
// and return the path to be used elsewhere
func ParseFlags() (string, error) {
    // String that contains the configured configuration path
    var configPath string

    // Set up a CLI flag called "-config" to allow users
    // to supply the configuration file
    flag.StringVar(&configPath, "config", "./config.yml", "./config.yml")

    // Actually parse the flags
    flag.Parse()

    // Validate the path first
    if err := ValidateConfigPath(configPath); err != nil {
        return "", err
    }

    // Return the configuration path
    return configPath, nil
}*/
//implement return pointer
/*func getConfig() (*Config, error){
cfgPath, err := ParseFlags()
    if err != nil {
        fmt.Println(err)
    }
    cfg, err := NewConfig(cfgPath)
    if err != nil {
    fmt.Println(err)
    }
    return cfg, err
}*//*
func readConfig() (*Config, error) {
    config := &Config{}
    cfgFile, err := ioutil.ReadFile("./config.yaml")
    if err != nil {
        return nil, err
    }
    err = yaml.Unmarshal(cfgFile, config)
    return config, err

}*/
func getIpAdrExp(exporterName string) []string{
    /*cfgPath, err := ParseFlags()
    if err != nil {
        fmt.Println(err)
    }*/
    //var c conf
    //c.getConf()
    cfg := getConfig()

	var list []string
    switch exporterName {
        case "systemStats":
           for i := range cfg.Hosts {
            //for i := 0; i < len(cfg.Hosts); i++ {
                if (cfg.Hosts[i].Exclude.SystemExporter == false) {
                    list = append(list, cfg.Hosts[i].Ipaddress)
                }
            }
        case "callStats":
            for i:= range cfg.Hosts {
                if (cfg.Hosts[i].Exclude.CallStats == false) {
                    list = append(list, cfg.Hosts[i].Ipaddress)
                }
            }
            //INFO: have a switch case on all exporters made, NB!: must remember exact exporternames inside each exporter
        }
return list
}

func getAuth(ipadr string) (username string, password string) {
    var u, p string
    cfg, err := getConfig()
    if err != nil {
       fmt.Println(err)
    }

   // map[adr]cfg.Hosts[i].Username
   // map[adr]cfg.Hosts[i].Username
    //yaml.Unmarshal(file_content, &map)
    for i:= range cfg.Hosts {
        if (cfg.Hosts[i].Ipaddress == ipadr) {
            u, p = cfg.Hosts[i].Username, cfg.Hosts[i].Password
        }
    }
   // return "test", "test"
    return u,p
}



func test() {
    ip := getIpAdrExp("systemStats")
    fmt.Println(ip)
}