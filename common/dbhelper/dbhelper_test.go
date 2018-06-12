package dbhelper

import (
	"path/filepath"
	"testing"
)

//Test ValidateDBconnDetails
func TestValidateDBConnDetails(t *testing.T) {
	dir, _ := filepath.Abs(".")

	jsonFile := filepath.Join(dir,"testdata","cldeploy-test.json")

	// read dbconfig into dbConfig struct
	dbMoodyscfg, err := GetDbConfig("MoodysDB", jsonFile)
	if err != nil {
		t.Fatal(err)
	}

	// should not return an error	
	if err := ValidateDBConnDetails(dbMoodyscfg.Username, dbMoodyscfg.Password, 
		dbMoodyscfg.HostName, dbMoodyscfg.Port, dbMoodyscfg.Database); (err != nil) {
		t.Error("ValidateDBConnDetails() error = " + err.Error())
	}

	// should return an error (BAD USER NAME)
	err = ValidateDBConnDetails("BAD USER NAME", dbMoodyscfg.Password, 
		dbMoodyscfg.HostName, dbMoodyscfg.Port, dbMoodyscfg.Database)
	if err == nil {
		t.Error("ValidateDBConnDetails() should have returned error")
	}
}