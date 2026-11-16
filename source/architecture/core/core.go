package main

import (
	"log/slog"
	"os"
	"yathahi/source/library"
)

func main() {
	slog.SetDefault(slog.New(slog.NewTextHandler(os.Stdout, nil)))
	requestsChannel := make(chan library.EntryData)

	go library.UdsServer(library.EntrySocketPath, requestsChannel)

	go Containersetup()

	go func() {

		for {
			receivedData := <-requestsChannel
			EngineInstaller(receivedData.Service)
		}
	}()

	select {}
}
