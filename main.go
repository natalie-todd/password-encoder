package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"password-encoder/server"
	"password-encoder/service"
	"syscall"
)

func main() {
	s := service.InitializeService()
	h := server.InitializeHandler(s)
	r := server.NewRouter(h)
	r.InitializeRoutes()

	server := http.Server{
		Addr:    ":8080",
		Handler: r,
	}

	termChan := make(chan os.Signal)
	signal.Notify(termChan, syscall.SIGINT)

	go func() {
		<-termChan // Blocks here until interrupted
		log.Print("SIGTERM received. Shutdown process initiated\n")
		server.Shutdown(context.Background())
	}()

	if err := server.ListenAndServe(); err != nil {
		if err.Error() != "http: Server closed" {
			log.Printf("HTTP server closed with: %v\n", err)
		}
		log.Printf("HTTP server shut down")
	}

}
