package library

import (
	"encoding/json"
	"log/slog"
	"net"
	"os"
)

func UdsServer(socketPath string, output chan<- EntryData) {
	// Force-remove old file so it doesn't crash on rerun
	os.RemoveAll(socketPath)

	// Create the listening socket
	listener, err := net.Listen("unix", socketPath)
	if err != nil {
		slog.Error("Failed to listen", "error", err)
	}
	defer listener.Close()
	slog.Info("UDS started listening...")

	// Wait for exactly ONE connection and then read from it
	conn, err := listener.Accept()
	if err != nil {
		slog.Error("Failed to accept connection", "error", err)
	}
	slog.Info("Listner accepted the connection!")
	defer conn.Close() // Ensure the connection is closed when this function exits

	var data EntryData
	err = json.NewDecoder(conn).Decode(&data)
	if err != nil {
		slog.Error("Failed to parse JSON", "error", err)
	}

	slog.Info("Entry Data parsed successfully!", "service", data.Service, "additional", data.Additional)

	// Return the populated struct
	output <- data
	slog.Info("The data has been transmitted!")
}
