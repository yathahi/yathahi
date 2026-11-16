package library

import (
	"encoding/json"
	"log/slog"
	"net"
	"time"
)

func UdsClient(data EntryData, socketPath string) {
	var conn net.Conn
	var err error

	// 1. Connect (dial) to your listening Unix domain socket
	start1 := time.Now()
	for {
		conn, err = net.Dial("unix", socketPath)
		if err != nil {
			slog.Error("Unable to listen to UDS", "error", err)
			if time.Since(start1) >= 30*time.Second {
				slog.Error("UDS-Client Timeout....")
				return // Exit immediately so we don't use a nil conn
			}
			time.Sleep(time.Second)
			continue
		}
		slog.Info("Socket is reachable!") // Changed from Error to Info
		break
	}

	defer conn.Close()
	start2 := time.Now()

	for {
		// 2. Initialize and populate your EntryData struct using the function parameters
		message := EntryData{
			Service:    data.Service,
			Additional: data.Additional,
		}

		// 3. Point the JSON Encoder directly at the socket stream and send the struct
		if err := json.NewEncoder(conn).Encode(message); err != nil {
			slog.Error("Unable to Encode or send data to UDS", "error", err)
		} else {
			slog.Info("Data has been sent to the Server")
			break
		}

		if time.Since(start2) >= 30*time.Second {
			slog.Error("UDS-Client Timeout....")
			break
		}
		time.Sleep(time.Second)
	}
}
