package cmd

import (
	"InstanceServerCX/logs"
	"fmt"
	"github.com/AlecAivazis/survey/v2"
	"github.com/spf13/cobra"
	"io"
	"os"
	"os/exec"
	"path/filepath"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "InstanceServerCX",
	Short: "Instance Server CX",
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

const (
	sitesFolder = "./.exports/sites"
	port        = "5555"
)

var selectedFolder string

func init() {
	folders := getFolders(sitesFolder)
	if len(folders) > 0 {
		selectFolder(folders)
		copyConfiguration()
	} else {
		fmt.Println("No folders were found in ", sitesFolder)
	}
}

func getFolders(pathFolder string) (folders []string) {
	files, err := os.ReadDir(pathFolder)
	if err != nil {
		logs.ErrorLog(fmt.Sprintln("Error reading directory: ", err.Error()))
	}

	for _, file := range files {
		if file.IsDir() {
			folders = append(folders, file.Name())
		}
	}

	return folders
}

func selectFolder(folders []string) {
	prompt := &survey.Select{
		Message: "Select an option:",
		Options: folders,
	}

	err := survey.AskOne(prompt, &selectedFolder, nil)
	if err != nil {
		logs.FatalLog(fmt.Sprintln(err))
	}
}

func copyFile(origen, destination string) error {
	source, err := os.Open(origen)
	if err != nil {
		return err
	}
	defer func(source *os.File) {
		err := source.Close()
		if err != nil {
			logs.FatalLog(fmt.Sprintln(err))
		}
	}(source)

	destinationFile, err := os.Create(destination)
	if err != nil {
		return err
	}
	defer func(destinoFile *os.File) {
		err := destinoFile.Close()
		if err != nil {
			logs.FatalLog(fmt.Sprintln(err))
		}
	}(destinationFile)

	_, err = io.Copy(destinationFile, source)
	if err != nil {
		return err
	}

	return nil
}

func copyConfiguration() {
	source := "./serve.json"
	destination := filepath.Join(sitesFolder, selectedFolder, "dist", "serve.json")

	err := copyFile(source, destination)
	if err != nil {
		logs.ErrorLog(fmt.Sprintln("Error when copying the configuration: ", err.Error()))
	}

	logs.SuccessLog("Copied server configuration")
	logs.SuccessLog(fmt.Sprintln("Serving from the folder: ", filepath.Join(sitesFolder, selectedFolder, "dist")))
	logs.SuccessLog(fmt.Sprintln("http://localhost: ", port, "/"))

	launchServer()
}

func launchServer() {
	logs.SuccessLog(fmt.Sprintln("Starting the server in the folder: ", selectedFolder))
	cmd := exec.Command("serve", "-p", port, filepath.Join(sitesFolder, selectedFolder, "dist"))
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err := cmd.Run()
	if err != nil {
		logs.ErrorLog(fmt.Sprintln("Error starting the server: ", err.Error()))
	}
}
