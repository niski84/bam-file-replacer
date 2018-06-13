package main

import (
	"github.com/niski84/bam-file-replacer/common/confighelper"

	"os"
	"flag"
	"log"
	"fmt"
	"path/filepath"
	"io/ioutil"
	"strings"
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
	cfg, err := confighelper.GetConfigNode("bam-config")
	if err != nil {
		log.Fatal(err)
	}

	
	
	log.Println("Searching recursively from: ", cfg["source"].(string))

	files, _:= FindFiles(cfg["source"].(string), "*")
	if err != nil {
		log.Fatal(err)
	} 

	if len(files) < 1 {
		log.Fatal("No files found under: ", cfg["source"].(string))	
	}

	log.Println("reading the keyfile file:",cfg["keyfile"].(string))
	templateText,err := getTemplate(cfg["keyfile"].(string))
	if err != nil {
		log.Fatal("oopsie.. keyfile not found in current directory", cfg["keyfile"].(string) )
	}

	log.Println("using text: \n ",templateText)

	// create destination folder if it doesn't exist
	if _, err := os.Stat(cfg["destination"].(string)); os.IsNotExist(err) {
		os.MkdirAll(cfg["destination"].(string), os.ModePerm)
	}

	// iterate over found files and replace source with destination dir
	for _,filef := range(files){
		log.Println("found: ", filef)
		fileDest := strings.Replace(filef, cfg["source"].(string), cfg["destination"].(string),1)
		
		// create destination folder if it doesn't exist
		if _, err := os.Stat(filepath.Dir(fileDest)); os.IsNotExist(err) {
			log.Println("creating dir:",filepath.Dir(fileDest))
			err = os.MkdirAll(filepath.Dir(fileDest), os.ModePerm)
			if err != nil {
				log.Fatal("error creating dir:",filepath.Dir(fileDest),err)
			}
		}
		
		log.Println("creating file:",fileDest)

		// write out the new file with the contents from the keyfile
	    ioutil.WriteFile(fileDest, []byte(templateText),os.FileMode(0644))
	}

	log.Println("bam-o! Processing complete")

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
			
			// add file as long as it isn't a dir
			for _,file := range filelist {
				if info, err := os.Stat(file); err != nil || !info.IsDir() {
					fl = append(fl, file)	
			}}
		return nil
	})
	return fl, err
}

