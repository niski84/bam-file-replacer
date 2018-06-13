package confighelper

import (
	"encoding/json"
	"os"
	"log"
	"io/ioutil"
	"errors"
)

// ConfigFile path to json config file
var ConfigFile string




// Set path to config file
func SetConfigFilePath(cfg string) (err error) {	

	if _, err := os.Stat(cfg); os.IsNotExist(err) {
		if err != nil {
			log.Printf(err.Error())
			return err
		}
	}
	
	ConfigFile = cfg

	return
}

// Get Config node as a map
func GetConfigNode(node string) (result map[string]interface{}, err error) {
	if _, err := os.Stat(ConfigFile); os.IsNotExist(err) {
		if err != nil {
			log.Println(err.Error())
			return nil, err
		}
	}

	// read json config file
	raw, err := ioutil.ReadFile(ConfigFile)
	if err != nil {
		log.Println(err.Error())
		return 
	}

	// unmarshal json into map
	json.Unmarshal([]byte(string(raw)), &result)
	
	if _, ok := result[node]; ok {
		return result[node].(map[string]interface{}), nil
	} else {
		return nil, errors.New("Node does not exist in config file.  bad json formatting? ")
	}
}

