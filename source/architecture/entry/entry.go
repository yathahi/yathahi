package main

import (
	"fmt"
	"io"
	"log/slog"
	"os"
	"path/filepath"
	"regexp"
	"sync"
	"yathahi/source/library"
	packagemanager "yathahi/source/library/package-manager"
)

//const containerName = ""+library.ProjectName+"" // later user custom container name can be added for multi servers with incrementing value like ekathvam-web-1 ekathvam-web-2 and so in
//const containerPath = "/"+library.ProjectName+"/containers"

var option string
var wg sync.WaitGroup
var coreWG sync.WaitGroup

// ---------------------------------------------------------------
// Main - Official Entry To The Ecosystem
// ---------------------------------------------------------------

func main() {
	fmt.Println("Verifying Superuser Privileage...")
	if os.Geteuid() == 0 {
		fmt.Println("Successfull. Running with root privileges!")
	} else {
		fmt.Println("Running as a regular user.\nRun with sudo or as root user.\nExiting....")
		return
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

	coreWG.Go(func() {

		packagemanager.PackageManager("zstd")
		packagemanager.PackageManager(library.CorePackageName)
	})

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
			fmt.Println("Waiting for background tasks to finish before exiting....")
			coreWG.Wait()
			wg.Wait() // Blocks until the goroutine calls wg.Done()
			fmt.Println("Exiting....")
			return // Triggers defer f.Close() cleanly

		default:
			fmt.Println("❌ Invalid selection. Please choose a valid option.")
		}
	}
}

//--------------------------------------------------------------------------
//	Web Entry - asks additional datas and process them to be sent.
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
	wg.Go(func() {
		coreWG.Wait()
		library.UdsClient(data, library.EntrySocketPath)
	})
}
