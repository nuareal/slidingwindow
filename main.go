package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/nuareal/slidingwindow/server"
)

func main() {

	port := ":8080"
	filePath := "cache.json"
	server := server.NewServer(filePath, port)

	done := make(chan bool)
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-quit

		// To guarantee that last state was save it would need to lock acess until shutdown
		if err := server.GetData().WriteToFile(filePath); err != nil {
			log.Fatalln("Could not save state to disk:", err)
		}

		if err := server.Server.Shutdown(context.TODO()); err != nil {
			log.Fatalln("Error on shutdown", err)
		}
		close(done)
	}()

	log.Println("Server is ready to handle requests at", server.Server.Addr)
	if err := server.Server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatal(err)
	}

	<-done
}
