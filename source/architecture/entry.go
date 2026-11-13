package main

import (
	"yathahi/source/library"
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"strings"
)

//const containerName = ""+library.ProjectName+"" // later user custom container name can be added for multi servers with incrementing value like ekathvam-web-1 ekathvam-web-2 and so in
//const containerPath = "/"+library.ProjectName+"/containers"

var option string

//--------------------------------------------------------------------------
// Official Entry To The Ecosystem
//--------------------------------------------------------------------------

func entry() {
	fmt.Println("Verifying Superuser Privileage...")
	if os.Geteuid() == 0 {
		fmt.Println("Successfull. Running with root privileges!")
	} else {
		fmt.Println("Running as a regular user.\nRun with sudo or as root user.\nExiting....")
		os.Exit(0)
	}

	determineDeviceRole()

}

//--------------------------------------------------------------------------
//	Deploying the User Selected Server Bootstrap
//--------------------------------------------------------------------------

func determineDeviceRole() error {
	fmt.Println("\nWelcome to", library.ProjectName)
	fmt.Println("\nWhat new server you would like to make?\n\n1. Administrator\n2. Client\n3. Web Server\n4. DNS Server\n5. Exit\n\nSelect an option:")
	for {
		fmt.Print("Enter option (1-4): ")
		fmt.Scanln(&option)

		switch option {
		case "1":

		case "2":

		case "3":
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
	data.Service = "web"
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
	coreInstall()
	library.UdsClient(data, library.EntrySocketPath)
}

//--------------------------------------------------------------------------
//	Install Yathahi-Core
//--------------------------------------------------------------------------

func coreInstall() {
	err := library.PackageManager("zstd")
	if err != nil {
		fmt.Println("Unable to install zstd", "error", err)
	}

	err = library.PackageManager(library.CoreBinaryURL)
	if err != nil {
		fmt.Println("Unable to Download"+library.ProjectName+"-Core", "error", err)
	}
	cmd1 := exec.Command("tar", "--zstd", "-xf", "/tmp/"+library.ProjectName+"-core.tar.zst", "-C", "/tmp/")
	output, err := cmd1.CombinedOutput()
	if err != nil {
		fmt.Printf("Unable to extract "+library.ProjectName+"-core.tar.zst %s", output)
	}
	err = os.MkdirAll("/opt/"+library.ProjectName+"/bin", 0700)
	if err != nil {
		fmt.Println(err)
	}
	cmd2 := exec.Command("mv", "/tmp/"+library.ProjectName+"-core", library.BinaryFolder)
	output, err = cmd2.CombinedOutput()
	if err != nil {
		fmt.Printf("Unable to move "+library.ProjectName+"-core %s", output)
	}
	fmt.Printf("Installing "+library.ProjectName+"-core %s", output)
	cmd3 := exec.Command("mv", "/tmp/"+library.ProjectName+"-core.service", "/etc/systemd/system/")
	output, err = cmd3.CombinedOutput()
	if err != nil {
		fmt.Println("Unable to move "+library.ProjectName+"-core.service", err)
	}
	fmt.Printf("Installing "+library.ProjectName+".service %s", output)

	cmd4 := exec.Command("systemctl", "daemon-reload")
	output, err = cmd4.CombinedOutput()
	if err != nil {
		fmt.Println("Unable to reload systemd", err)
	}
	fmt.Println("Reloading systemd services")

	cmd5 := exec.Command("systemctl", "enable", "--now", library.ProjectName+"-core.service")
	output, err = cmd5.CombinedOutput()
	if err != nil {
		fmt.Println("Unable to enable "+library.ProjectName+"-core.service", err)
	}
	fmt.Println("Enabling " + library.ProjectName + "-core Service")

	cmd6 := exec.Command("systemctl", "is-active", library.ProjectName+"-core.service")

	output, err = cmd6.CombinedOutput()

	status := strings.TrimSpace(string(output))
	if err != nil {
		fmt.Println("Something is wrong! "+library.ProjectName+"-Core is not running...", err)
	}
	if status == "active" {
		fmt.Println("Core is running actively.")
	}

}
