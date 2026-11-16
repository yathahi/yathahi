package main

import (
	"log/slog"
	"os/exec"
	"strings"
	"yathahi/source/library"
)

//--------------------------------------------------------------------------
//	Engine Installer
//--------------------------------------------------------------------------

func EngineInstaller(engine string) {

	err := library.PackageManager(library.BinaryURL + engine + ".tar.zst")
	if err != nil {
		slog.Error("Unable to Download",
			slog.String("engine", engine),
			slog.Any("error", err),
		)

	}

	cmd1 := exec.Command("tar", "--zstd", "-xf", "/tmp/"+engine+".tar.zst", "-C", "/tmp/")
	output, err := cmd1.CombinedOutput()
	if err != nil {
		slog.Error("Unable to extract "+engine+".tar.zst", "output", output)
	}
	cmd2 := exec.Command("mv", "/tmp/"+engine, library.BinaryFolder)
	output, err = cmd2.CombinedOutput()
	if err != nil {
		slog.Error("Unable to mv ", engine, "from /tmp",
			"binary_folder", library.BinaryFolder,
			"error", err,
			"output", string(output))
	} else {
		slog.Info("binary is successfully moved",
			slog.String("engine", engine),
			slog.String("output", string(output)),
		)
	}
	cmd3 := exec.Command("mv", "/tmp/"+engine, library.BinaryFolder)
	output, err = cmd3.CombinedOutput()
	if err != nil {
		slog.Error("Unable to move",
			slog.String("engine", engine),
			slog.String("output", string(output)),
		)
	}
	// 1. Reload the systemd daemon
	reloadCmd := exec.Command("systemctl", "daemon-reload")
	if output, err := reloadCmd.CombinedOutput(); err != nil {
		slog.Error("Unable to reload systemd daemon", "error", err, "output", string(output))
		return
	}

	// 2. Enable and start the service
	serviceName := engine + ".service"
	enableCmd := exec.Command("systemctl", "enable", "--now", serviceName)
	if output, err := enableCmd.CombinedOutput(); err != nil {
		slog.Error("Unable to enable service", "service", serviceName, "error", err, "output", string(output))
		return
	}

	slog.Info("Enabling ", engine, "Service")
	cmd5 := exec.Command("systemctl", "is-active", engine+".service")
	output, err = cmd5.CombinedOutput()
	status := strings.TrimSpace(string(output))
	if err != nil {
		slog.Error("Something is wrong!", engine, "is not running...", "error", err)
	}
	if status == "active" {
		slog.Info("engine", engine, "is running actively.")
	}

}
