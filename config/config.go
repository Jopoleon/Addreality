package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
)

//host = "localhost"
//port = 5432
//user = "postgres"
//password = "qwe"
//dbname = "test"
type ConfigType struct {
	ServerPort string
	Postgre    string
	Host       string
	Port       string
	User       string
	Password   string
	DBname     string

	Metric_1_Max int
	Metric_1_Min int
	Metric_2_Max int
	Metric_2_Min int
	Metric_3_Max int
	Metric_3_Min int
	Metric_4_Max int
	Metric_4_Min int
	Metric_5_Max int
	Metric_5_Min int
}

var Config = new(ConfigType)

func GetConfig() *ConfigType {
	log.Printf("GetConfig values: %+v", Config)
	return Config
}

func InitConf(filename string) error {
	//log.Println("....InitConf used")
	var c = &ConfigType{}
	if filename == "" {
		return fmt.Errorf(`Error: Don't dicribe path to a config file.`)
	}
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		if err = genConfig(filename); err != nil {
			return err
		}
	}

	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return fmt.Errorf("Can't read config file: %s", err)
	}
	log.Println("Config file " + filename + " found. Reading...")

	if err = json.Unmarshal(data, c); err != nil {
		fmt.Errorf("Can't read config file: %s", err)
		return err
	}
	Config = c
	//log.Println("...... Config file values: ", Config)
	return nil
}

func genConfig(filename string) error {
	log.Println("NOPE. There is no such config file ", filename)
	log.Println("Configuration file not found. Created new with name " + filename + ". " +
		"\n 		     Please, fill it with values you need and RESTART application")
	f, err := os.Create(filename)
	if err != nil {
		log.Println("ReadConfigFile os.Create error: ", err)
		return err
	}
	f.Close()

	var initjson = ConfigType{
		ServerPort: "3000",
		Postgre:    "ssa",
		DBname:     "addrealty",
		Host:       "localhost",
		Port:       "5432",
		User:       "postgres",
		Password:   "qwe",
	}

	writebytes, err := json.MarshalIndent(initjson, "", "\t")
	if err != nil {
		panic(err)
		return err
	}
	err = ioutil.WriteFile(filename, writebytes, 0644)
	if err != nil {
		panic(err)
		return err
	}
	os.Exit(5)
	return nil
}
