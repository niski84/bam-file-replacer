package clihelper


import(
	"bufio"
	"fmt"
	"os"
	"strings"
	"path/filepath"
)

// Prompt user for input. Default text is used if no input given.
func PromptUser(prompttext string, defaulttext string) (input string, err error) {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print(prompttext+" ["+defaulttext+"]: ")
	input, err = reader.ReadString('\n')
	
	if strings.TrimSpace(input) == "" {
		input = defaulttext
	}

	input = strings.TrimSpace(input)
	
	if err != nil {
		return input, err
	}
	return input, nil
}

// Returns path the .exe is executing from
func GetAppDir() (cwd string, err error){
	cwd, err = filepath.Abs(filepath.Dir(os.Args[0]))
    if err != nil {
            return cwd, err
    }
	return cwd, err
}