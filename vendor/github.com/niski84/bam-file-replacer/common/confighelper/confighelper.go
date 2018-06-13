package confighelper

import (
	"encoding/json"
	"os"
	"log"
	"io/ioutil"
	"errors"
	"reflect"
)

// ConfigFile path to json config file
var ConfigFile string

// PackageRoot is the folder of the cl package
var PackageRoot string

// SetPackageRoot and validate structure
func SetPackageRoot(path string) (err error) {	
	if _, err := os.Stat(path); os.IsNotExist(err) {
		if err != nil {
			log.Printf(err.Error())
			return err
		}
	}

	// TODO: VALIDATE PACKAGE STRUCTURE
	PackageRoot = path

	return 
}

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
		return nil, errors.New("Node does not exist in config file: ")
	}
}

// FillStruct reflect a map into struct
// map and struct names MUST match
func FillStruct(data map[string]interface{}, result interface{}) {
	t := reflect.ValueOf(result).Elem()
    for k, v := range data {
		val := t.FieldByName(k)
        val.Set(reflect.ValueOf(v))
    }
}