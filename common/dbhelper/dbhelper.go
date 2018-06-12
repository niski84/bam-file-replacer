// Helper functions related to execution of sql scripts using psql
package dbhelper

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"regexp"
	"runtime"
	"deployment/cl/common/confighelper"
	"database/sql"
	_ "github.com/lib/pq"
)

// struct for db connection details
type DbConfig struct {
	HostName    	string 
	Port 			string 
	Username 		string 
	Password 		string 
	Database	    string 
}

func GetDbConfig(node string, jsonFile string) (dbConfigStruct *DbConfig, err error) {
	
	// jsonFile := filepath.Join(dir,"testdata","cldeploy-test.json")
	confighelper.SetConfigFilePath(jsonFile)

	// get json config node
	jsonCfgMap, err := confighelper.GetConfigNode(node)
	if err != nil {
		return
	}
	
	// convert map to dbConfig struct
	dbConfigStruct = &DbConfig{}	
	confighelper.FillStruct(jsonCfgMap, dbConfigStruct)

	return 
	
}


func CheckPsqlExists() error {
	_, err := exec.LookPath("psql")
	return err
}

func CheckPsqlVer() (bool, string) {

	const ReqMajor = "9"
	const ReqMinor = "6"
	const ReqRevision = "2"
	var ReqVer = ReqMajor+"."+ReqMinor+"."+ReqRevision

	out, err := exec.Command("psql.exe", "--version").Output()
    if err != nil {
        log.Fatal(err)
	}

	versionRegex, _ := regexp.Compile(`\d\.\d\.\d`)
	version := versionRegex.FindAllStringSubmatch(string(out), -1)[0][0]

	splitVer := strings.Split(version,".")
	VerMajor := splitVer[0]
	VerMinor := splitVer[1]
	VerRevision := splitVer[2]

	if VerMajor < ReqMajor {return false, ReqVer}
	if VerMinor < ReqMinor {return false, ReqVer}
	if VerRevision < ReqRevision {return false, ReqVer}
	
	return true, ReqVer
}
 
func ValidateDBConnDetails(user string, password string, hostname string, port string, databaseName string) (err error) {

	os.Setenv("VERBOSITY", "verbose")
	os.Setenv("PGPASSWORD", password)
	os.Setenv("PGCLIENTENCODING", "UTF8")

	cmdArgs := []string{"-h", hostname, "-p", port, "-U", user, "-d", databaseName}
	cmdName := "psql.exe"
	cmd := exec.Command(cmdName, cmdArgs...)

	stdoutStderr, err := cmd.CombinedOutput()

	if strings.Contains(string(stdoutStderr), "FATAL") {
		return errors.New("Error connecting to database. Check your connection details in the config " + string(stdoutStderr))
	}

	return nil
}

// assumes cmdArgs either contains command at least ("-f", "install.sql" | "-c",  ) 
func ExecSqlCmdArgs(dbcfg *DbConfig, cmdArgs []string, packageRoot string, logSuffix string) (err error, stdErrOut string) {

	os.Setenv("VERBOSITY", "verbose")
	os.Setenv("PGPASSWORD", dbcfg.Password)
	os.Setenv("PGCLIENTENCODING", "UTF8")

	args := []string{"-h", dbcfg.HostName, "-p", dbcfg.Port, "-U", dbcfg.Username, "-d", dbcfg.Database}
	args = append(args, cmdArgs...)
	args = append(args, "-v", "ON_ERROR_STOP=1", "--pset", "pager=off", "--single-transaction", "-b", "-L",
		filepath.Join(packageRoot, (logSuffix + "-output.log")))

	cmdName := "psql.exe"

	log.Println("\nExecuting: "+cmdName, strings.Join(args, " ")+"\n")

	cmd := exec.Command(cmdName, args...)

	stderr, err := cmd.StderrPipe()
	if err != nil {
		log.Println(err)
	}

	if err := cmd.Start(); err != nil {
		log.Println(err)
	}

	slurp, _ := ioutil.ReadAll(stderr)
	stdErrOut = string(slurp)

	var errorsOnly string
	if err := cmd.Wait(); err != nil {
		lines:=strings.Split(stdErrOut,`psql:`)
		for _,line := range(lines) {
			if strings.Contains(string(line), "ERROR:") {
				errorsOnly += string(line)
			}
		}
		
		// return the last line; sometimes it doesn't say "ERROR"
		errorsOnly += lines[len(lines)-1]

		// write errors to file:
		ioutil.WriteFile(filepath.Join(packageRoot,(logSuffix + "-" + "error.log")), []byte(stdErrOut), os.FileMode(0644))
		
		return err, errorsOnly
	}

	return err, errorsOnly 
}
func ExecSql(hostname string, password string, port string, user string, databaseName string, installSQL string, packageRoot string, dbtype string, outputlog string) (err error, stdErrOut string) {

	os.Setenv("VERBOSITY", "verbose")
	os.Setenv("PGPASSWORD", password)
	os.Setenv("PGCLIENTENCODING", "UTF8")

	cmdArgs := []string{"-h", hostname, "-p", port, "-U", user, "-d", databaseName, "-f",
		installSQL, "-v", "ON_ERROR_STOP=1", "--pset", "pager=off", "--single-transaction", "-b", "-L",
		filepath.Join(packageRoot, (dbtype + "-" + outputlog))}

	cmdName := "psql.exe"

	log.Println("\nExecuting: "+cmdName, strings.Join(cmdArgs, " ")+"\n")

	cmd := exec.Command(cmdName, cmdArgs...)

	stderr, err := cmd.StderrPipe()
	if err != nil {
		log.Println(err)
	}

	if err := cmd.Start(); err != nil {
		log.Println(err)
	}

	slurp, _ := ioutil.ReadAll(stderr)
	stdErrOut = string(slurp)

	var errorsOnly string
	if err := cmd.Wait(); err != nil {
		lines:=strings.Split(stdErrOut,`psql:`)
		for _,line := range(lines) {
			if strings.Contains(string(line), "ERROR:") {
				errorsOnly += string(line)
			}
		}
		
		// return the last line; sometimes it doesn't say "ERROR"
		errorsOnly += lines[len(lines)-1]

		// write errors to file:
		ioutil.WriteFile(filepath.Join(packageRoot,(dbtype + "-" + "error.log")), []byte(stdErrOut),os.FileMode(0644))
		
		return err, errorsOnly
	}

	return err, errorsOnly 
}

// Build install.sql. Can include sql statements and paths to .sql files
// outputFile location: install-<dbtype>.sql 
// sqlList: string array containg sql paths and commands
// cwd: true/false: meta change working directory into folder containing .sql
func GenerateInstallSql(outputFile string, sqlList []string, cwd bool) (installsql string, err error) {
	var contents string

	// forcing utf-8.
	contents +=`\encoding 'UTF8'`+"\r\n"

	// write .sql paths and commands to file.
	// dlmap is expected to be mix of commands and file paths
	for i := range(sqlList){
		var lineSuffix = ""
		if sqlList[i][len(sqlList[i])-4:] == ".sql"{
			if cwd == true {
				contents += `\cd `+ `'`+ strings.Replace(filepath.Dir(sqlList[i]),`\`,"/",-1)+`'`+"\r\n"
			}

			lineSuffix = `\i `
			contents += lineSuffix + `'`+ strings.Replace(sqlList[i],`\`,"/",-1)+`'`+"\r\n"
		} else {
			// this is an sql statement
			contents += sqlList[i]+"\r\n"
		}
	}

	// ioutil; no need to defer.close
	err = ioutil.WriteFile(outputFile, []byte(contents),os.FileMode(0644))
    if err != nil {
        return "", err
    }
	return outputFile, err
} 

// Execute sql statement using postgres driver
func PqSql(dbcfg *DbConfig, statement string)(output []string, err error){

	dbinfo := fmt.Sprintf("host=%s port=%s user=%s "+
		"password=%s dbname=%s sslmode=disable",
		dbcfg.HostName, dbcfg.Port, dbcfg.Username, dbcfg.Password, dbcfg.Database)


	db, err := sql.Open("postgres", dbinfo)
	if err != nil {
		log.Println(err.Error()) 
	}
	defer db.Close()

	rows, err := db.Query(statement)
    if err != nil {
        log.Printf("hrm %c\n", err)
    }

	output = make([]string, 1)

	// TODO: needs to support returning []interface{} for complext results
    for rows.Next() {
        var name string
        var val string

        err = rows.Scan(&name, &val)
        if err != nil {
         log.Printf("rows.Scan error: %v\n", err)
         log.Printf("hrm %c\n", err)
        }

		output = append(output, val)

	}
	
	return output, err 

}

// Return path of package 
func PkgPath() (path string) {
	_, path, _, _ = runtime.Caller(0)
	path = filepath.Dir(path)
	return
}

// PG_DUMP to a dump file
func PgDump(dbcfg *DbConfig, dumpFile string) (err error) {

	os.Setenv("PGCLIENTENCODING", "UTF8")
	os.Setenv("PGDATABASE", dbcfg.Database)
	

	args := []string{"-h", dbcfg.HostName, "-p", dbcfg.Port, "-U", dbcfg.Username, "-w" ,
	 	"-Fp", "-f", dumpFile}
	
	cmdName := "pg_dump.exe"

	log.Println("\nExecuting: "+cmdName, strings.Join(args, " ")+"\n")

	cmd := exec.Command(cmdName, args...)

	_, err = cmd.StderrPipe()
	if err != nil {
		return err
	}

	if err := cmd.Start(); err != nil {
		return err
	}

	if err := cmd.Wait(); err != nil {
		return err
	}

	return  
}

func PgRestore(dbcfg *DbConfig, dumpFile string) (err error) {

	os.Setenv("PGCLIENTENCODING", "UTF8")

	args := []string{"-h", dbcfg.HostName, "-p", dbcfg.Port, "-U", dbcfg.Username, "-w" ,
	 	"-Fp", "-f", dumpFile}
	
	cmdName := "pg_restore.exe"

	log.Println("\nExecuting: "+cmdName, strings.Join(args, " ")+"\n")

	cmd := exec.Command(cmdName, args...)

	_, err = cmd.StderrPipe()
	if err != nil {
		return err
	}

	if err := cmd.Start(); err != nil {
		return err
	}

	if err := cmd.Wait(); err != nil {
		return err
	}

	return  
}
