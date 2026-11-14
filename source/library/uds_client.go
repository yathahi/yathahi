package library

import (
	"encoding/json"
	"log/slog"
	"net"
	"time"
)

func UdsClient(data EntryData, socketPath string) {
	// 1. Connect (dial) to your listening Unix domain socket
	start := time.Now()
	for {
		conn, err := net.Dial("unix", socketPath)
		if err != nil {
			slog.Error("Unable to listen to UDS ", "error", err)
			continue
		}
		defer conn.Close()

		// 2. Initialize and populate your EntryData struct using the function parameters
		message := EntryData{
			Service:    data.Service,    // Use the passed 'service' argument
			Additional: data.Additional, // Use the passed 'additional' argument
		}

		// 3. Point the JSON Encoder directly at the socket stream and send the struct
		err = json.NewEncoder(conn).Encode(message)
		if err != nil {
			slog.Error("Unable to Encode or sent data  to UDS ", "error", err)
		} else {
			slog.Info("Data has sent to the Server")
			break
		}
		if time.Since(start) >= 30*time.Second {
			slog.Error("UDS-CLient Timeout....")
			break
		}
		time.Sleep(time.Second)
	}

}
