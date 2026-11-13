package main

import (
	"yathahi/source/library"
	"errors"
	"fmt"
	"log/slog"
	"os"
	"os/exec"
	"strings"
)

const stateFilePath = "/var/lib/yathahi/state.json"

func main() {
	slog.SetDefault(slog.New(slog.NewTextHandler(os.Stdout, nil)))
	slog.Error("The 1socket path is ", "output", library.EntrySocketPath) // Temp-Debugger
	requestsChannel := make(chan library.EntryData)
	slog.Error("The channel has just Created!")				// Temp-Debugger - suspect this is not running something wrong with channel creation , function is just stuck here

	go library.UdsServer(library.EntrySocketPath, requestsChannel) 
	slog.Error("The 2socket path is ", "output", library.EntrySocketPath) // Temp-Debugger

	go func() {
		for {
			receivedData := <-requestsChannel
			slog.Info("Recived data is ", "output", receivedData.Service) // Temp-Debugger
			switch receivedData.Service {
			case "administrator":

			case "client":

			case "web":
				Containersetup()
				webEngineInstaller()
			case "dns":

			default:
				slog.Error("Critical Error : Unknown Service type received from entry data.")

			}
		}
	}()

	_, err := os.Stat(stateFilePath)

	if err == nil {
		slog.Info("State file", "exists on", stateFilePath)

	}

	if errors.Is(err, os.ErrNotExist) {
		slog.Warn("State file", "not_found on", stateFilePath)
	}

	// Handles scenarios where the file might exist but cannot be accessed (e.g., permission issues)
	slog.Error("failed to determine  state file status", stateFilePath, "Error :")

	select {}
}

//--------------------------------------------------------------------------
//	Web-Engine Installer
//--------------------------------------------------------------------------

func webEngineInstaller() {

	err := library.PackageManager("https://commons.wikimedia.org/wiki/File:Ko%C5%9Bcielec_view_2.jpg")
	if err != nil {
		slog.Error("Unable to Download web-engine", "error", err)
	}

	cmd1 := exec.Command("tar", "--zstd", "-xf", "/tmp/web-engine.tar.zst", "-C", "/tmp/")
	output, err := cmd1.CombinedOutput()
	if err != nil {
		slog.Error("Unable to extract web-engine.tar.zst", "output", output)
	}
	cmd2 := exec.Command("mv", "/tmp/web-engine", library.BinaryFolder)
	output, err = cmd2.CombinedOutput()
	if err != nil {
		slog.Error("Unable to mv web-engine from /tmp",
			"binary_folder", library.BinaryFolder,
			"error", err,
			"output", string(output))
	} else {
		slog.Info("web-engine binary is successfully moved",
			"output", string(output))
	}
	cmd3 := exec.Command("mv", "/tmp/web-engine", library.BinaryFolder)
	output, err = cmd3.CombinedOutput()
	if err != nil {
		slog.Error("Unable to move Web Engine!", "output", output)
	}
	cmd4 := exec.Command("systemctl", "daemon-reload",
		"&&", "systemctl", "enable", "--now", "web-engine.service")
	output, err = cmd4.CombinedOutput()
	if err != nil {
		slog.Error("Unable to enable web-engine.service", "error", err)
	}
	slog.Info("Enabling "+library.ProjectName+"-core Service")
	cmd5 := exec.Command("systemctl", "is-active", "web-engine.service")
	output, err = cmd5.CombinedOutput()
	status := strings.TrimSpace(string(output))
	if err != nil {
		slog.Error("Something is wrong! web-engine is not running...", "error", err)
	}
	if status == "active" {
		slog.Info("Web-Engine is running actively.")
	}


}

//--------------------------------------------------------------------------
//	Initiating The Container
//--------------------------------------------------------------------------


func Containersetup() {
	slog.Info("Initializing LXC Deployment...", "container", library.ProjectName)

	// 1. Enable lxc-net Service
	slog.Info("Enabling lxc-net service...")
	out, err := exec.Command("systemctl", "enable", "--now", "lxc-net").CombinedOutput()
	if err != nil {
		slog.Error("Failed to enable lxc-net",
			"error", err,
			"output", string(out),
		)
		os.Exit(1)
	}

	// 2. Create the LXC Container
	slog.Info("Creating Container...", "container", library.ProjectName)
	out, err = exec.Command(
		"lxc-create",
		"-n", library.ProjectName,
		"-t", "download",
		"--",
		"--dist", "debian",
		"--release", "trixie",
		"--arch", "amd64",
	).CombinedOutput()
	if err != nil {
		slog.Error("Failed to create container",
			"container", library.ProjectName,
			"error", err,
			"output", string(out),
		)
		os.Exit(1)
	}
	if len(out) > 0 {
		slog.Info("Container creation output", "output", string(out))
	}

	// 3. Start the Container
	slog.Info("Starting LXC Container...", "container", library.ProjectName)
	out, err = exec.Command("lxc-start", "-n", library.ProjectName).CombinedOutput()
	if err != nil {
		slog.Error("Failed to start container",
			"container", library.ProjectName,
			"error", err,
			"output", string(out),
		)
		return
	}

	slog.Info("LXC Deployment Complete!", "container", library.ProjectName)
	slog.Info("Configuring Firewall To make internet access inside the lxc Container")
	err = provisionFirewall()
	if err != nil {
		slog.Error("Something went wrong with Firewall", "error", err)
	}
}

//-----------------------------------------------------------------------------
//	Provision Firewall
//-----------------------------------------------------------------------------

func provisionFirewall() error {

	// 1. Check for firewalld
	fwActiveCmd := exec.Command("systemctl", "is-active", "firewalld")
	fwOut, _ := fwActiveCmd.Output()

	ufwActiveCmd := exec.Command("systemctl", "is-active", "ufw")
	ufwOut, _ := ufwActiveCmd.Output()

	if strings.TrimSpace(string(fwOut)) == "active" {
		err := configureFirewalld()
		if err != nil {
			return err
		}

	} else if strings.TrimSpace(string(ufwOut)) == "active" {
		err := configureUfw()
		if err != nil {
			return err
		}
	}

	slog.Warn("Neither firewalld nor ufw is currently active..\nAn active firewall is Recommended.")
	return nil

}

//--------------------------------------------------------------------------
//	Configure Firewalld
//--------------------------------------------------------------------------

func configureFirewalld() error {

	slog.Info("Detected active firewalld. Applying rules...")

	rules := [][]string{
		{"--zone=trusted", "--add-interface=lxcbr0", "--permanent"},
		{"--zone=public", "--add-masquerade", "--permanent"},
		{"--reload"},
	}

	for _, args := range rules {
		cmd := exec.Command("firewall-cmd", args...)
		if out, err := cmd.CombinedOutput(); err != nil {
			return fmt.Errorf("firewall-cmd %s failed: %v (output: %s)\n", strings.Join(args, " "), err, string(out))
		}
	}
	slog.Info("firewalld configuration complete.")
	return nil

}

//--------------------------------------------------------------------------
//	Configure UFW
//--------------------------------------------------------------------------

func configureUfw() error {
	// 2. Check for UFW
	slog.Info("Detected active UFW. Applying rules...")

	rules := [][]string{
		{"allow", "in", "on", "lxcbr0"},
		{"route", "allow", "in", "on", "lxcbr0"},
	}

	for _, args := range rules {
		cmd := exec.Command("ufw", args...)
		if out, err := cmd.CombinedOutput(); err != nil {
			return fmt.Errorf("ufw %s failed: %v (output: %s)\n", strings.Join(args, " "), err, string(out))
		}
	}
	slog.Info("UFW configuration complete.")
	return nil
}


