package main

import (
	"fmt"
	"io"
	"log"
	"log/slog"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
	"yathahi/source/library"
)

//const containerName = ""+library.ProjectName+"" // later user custom container name can be added for multi servers with incrementing value like ekathvam-web-1 ekathvam-web-2 and so in
//const containerPath = "/"+library.ProjectName+"/containers"

var option string

// ---------------------------------------------------------------
// Main - Official Entry To The Ecosystem
// ---------------------------------------------------------------

func main() {
	fmt.Println("Verifying Superuser Privileage...")
	if os.Geteuid() == 0 {
		fmt.Println("Successfull. Running with root privileges!")
	} else {
		fmt.Println("Running as a regular user.\nRun with sudo or as root user.\nExiting....")
		os.Exit(0)
	}

	folderpath := filepath.Dir(library.EntryLogFile)
	var out io.Writer = io.Discard

	if os.MkdirAll(folderpath, 0755) == nil {
		if f, err := os.OpenFile(library.EntryLogFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644); err == nil {
			out = f
			defer f.Close()
		}
	}

	slog.SetDefault(slog.New(slog.NewTextHandler(out, nil)))

	go checkCoreAvailablility()

	fmt.Println("\nWelcome to", library.ProjectName)
	fmt.Println("\nWhat new server you would like to make?\n\n1. Administrator\n2. Client\n3. Web Server\n4. DNS Server\n5. Exit\n\nSelect an option:")
	for {
		fmt.Print("Enter option (1-5): ")
		fmt.Scanln(&option)

		switch option {
		case "1":

		case "2":

		case "3":
			fmt.Println("The Service is getting established.....")
			webEntry()

		case "4":

		case "5":
			fmt.Println("Exiting....")
			os.Exit(0)
		default:
			fmt.Println("❌ Invalid selection. Please choose a valid option.")
		}
	}
}

//--------------------------------------------------------------------------
//	Web Entry - asks additional dtas and process them to be sent.
//--------------------------------------------------------------------------

func webEntry() {

	var data library.EntryData
	data.Service = library.WebPackageName
	for {

		fmt.Printf("What is the domain that you wanted to resolved to this page? (like example.com,google.com,github.org: )")
		fmt.Scanln(&data.Additional)

		matched, _ := regexp.MatchString(`^[a-zA-Z0-9-]+\.[a-zA-Z]{2,}$`, data.Additional)

		if !matched {
			fmt.Printf("Invalid format: %s   (enter vaild domain like "+library.ProjectName+".eka or example.com)", data.Additional)
		} else {
			break
		}
	}
	library.UdsClient(data, library.EntrySocketPath)
}

//--------------------------------------------------------------------------
//	Install Yathahi-Core
//--------------------------------------------------------------------------

func checkCoreAvailablility() {

	_, err := os.Stat(library.BinaryFolder + library.CorePackageName)
	if err != nil {
		slog.Info("Core is not installed!")
		coreInstall()
	}
	slog.Info("Core is installed")

	cmd6 := exec.Command("systemctl", "is-active", library.ProjectName+"-core.service")

	output, err := cmd6.CombinedOutput()

	status := strings.TrimSpace(string(output))
	if err != nil {
		slog.Info("Core is not enabled!")
		enableCore()
	}
	if status == "active" {
		slog.Info("Core Already exists")
	}

}

//--------------------------------------------------------------------------
//	Install the core
//--------------------------------------------------------------------------

func coreInstall() {
	err := library.PackageManager("zstd")
	if err != nil {
		log.Println("Unable to install zstd", "error", err)
	}

	err = library.PackageManager(library.BinaryURL + library.CorePackageName)
	if err != nil {
		log.Println("Unable to Download "+library.ProjectName+"-Core", "error", err)
	}
	cmd1 := exec.Command("tar", "--zstd", "-xf", "/tmp/"+library.ProjectName+"-core.tar.zst", "-C", "/tmp/")
	output, err := cmd1.CombinedOutput()
	if err != nil {
		log.Printf("Unable to extract "+library.ProjectName+"-core.tar.zst %s", output)
	}
	err = os.MkdirAll(library.BinaryFolder, 0700)
	if err != nil {
		log.Println(err)
	}
	cmd2 := exec.Command("mv", "/tmp/"+library.ProjectName+"-core", library.BinaryFolder)
	output, err = cmd2.CombinedOutput()
	if err != nil {
		log.Printf("Unable to move "+library.ProjectName+"-core %s", output)
	}
	log.Printf("Installing "+library.ProjectName+"-core %s", output)
	cmd3 := exec.Command("mv", "/tmp/"+library.CoreServiceName, "/etc/systemd/system/")
	output, err = cmd3.CombinedOutput()
	if err != nil {
		log.Println("Unable to move "+library.ProjectName+"-core.service", err)
	}
	log.Printf("Installing "+library.ProjectName+"-core.service %s", output)

	enableCore()
}

//--------------------------------------------------------------------------
//	Enable The Yathahi-Core
//--------------------------------------------------------------------------

func enableCore() {

	cmd4 := exec.Command("systemctl", "daemon-reload")
	output, err := cmd4.CombinedOutput()
	if err != nil {
		log.Println("Unable to reload systemd", err)
	}
	log.Println("Reloading systemd services")

	cmd5 := exec.Command("systemctl", "enable", "--now", library.CoreServiceName)
	output, err = cmd5.CombinedOutput()
	if err != nil {
		log.Println("Unable to enable "+library.ProjectName+"-core.service", err)
	}
	log.Println("Enabling " + library.ProjectName + "-core Service")

	cmd6 := exec.Command("systemctl", "is-active", library.CoreServiceName)

	output, err = cmd6.CombinedOutput()

	status := strings.TrimSpace(string(output))
	if err != nil {
		log.Println("Something is wrong! "+library.ProjectName+"-Core is not running...", err)
	}
	if status == "active" {
		log.Println("Core is running actively.")
	}
}
