package main

import (
	"log/slog"
	"os"
	"os/exec"
	"yathahi/source/library"
)

//const socketPath = "/tmp/ekathvam.engine.sock"

func main() {
	slog.SetDefault(slog.New(slog.NewTextHandler(os.Stdout, nil)))
	slog.Info("Initializing The Web Engine....")

	cmd := exec.Command("lxc-attach", "-n", library.ProjectName, "--", "apt-get", "install", "-y", "nginx")
	output, err := cmd.CombinedOutput()
	if err != nil {
		slog.Error("Failed: %s\n Error: %v\n", string(output), err)
	}
	slog.Info("Success!")

	cmd1 := exec.Command("lxc-attach", "-n", library.ProjectName,"--", "systemctl", "enable", "--now", "nginx")
	output, err = cmd1.CombinedOutput()
	if err != nil {
		slog.Error("Failed: %s\n Error: %v\n", string(output), err)
	}
	slog.Info("Success!")

	select {}
}
