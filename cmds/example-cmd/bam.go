package main

import (
	"templates/common/confighelper"

	"os"
	"flag"
	"log"
	"fmt"
	"path/filepath"
	"io/ioutil"
)

// Process command line flags (arguments)
func processFlags() (jsonFile string) {
	
	flag.StringVar(&jsonFile, "-c", `./bam.json`, "path to config file")
	flag.Parse()

	return 
}

func main() {
	
	log.Println("BAM! - file replacer!\n\n")
	
	// get json config file
	jsonFile := processFlags()
	
	// set the config file path
	confighelper.SetConfigFilePath(jsonFile)
	log.Println("\nExecuting using Config:" + jsonFile)
	
	
	// read config
	cfg, err := confighelper.GetConfigNode("Template")
	if err != nil {
		log.Fatal(err)
	}

	
	log.Println(cfg["glob"].(string))

	files, _:= FindFiles(cfg["destination"].(string), cfg["glob"].(string))
	if err != nil {
		log.Fatal(err)
	} 

	templateText,_ := getTemplate(cfg["template-file"].(string))
	
	log.Println("using text: \n ",templateText)

	for i,filef := range(files){
		log.Println(i, " replacing text in file: ", filef)
		// write errors to file:
	    ioutil.WriteFile(filef, []byte(templateText),os.FileMode(0644))

	}

	
	log.Println(cfg["destination"],templateText)

	log.Println("Processing complete")

}

func getTemplate(location string) (templateText string, err error){
	raw, err := ioutil.ReadFile(location)

	templateText = string(raw)

	if err != nil {
		log.Printf(err.Error())
		return
	}
return
}


// FindFiles searches for files matching the pattern recursively
func FindFiles(dir string, pattern string) ([]string, error) {
	
	fl := []string{}
	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			fmt.Printf("error accessing path %q:%v\n", dir, err)
			return err
		}
			// Get list of files in path
			filelist, err := filepath.Glob(filepath.Join(path, pattern))
			if err != nil {
				fmt.Printf("Error in search pattern %v\n", err)
				return err
			}
			
			fl = append(fl, filelist...)
				
		return nil
	})
	return fl, err
}

