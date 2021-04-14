package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"password-encoder/server"
	"password-encoder/service"
	"sync"
	"syscall"
)

func main() {
	var wg sync.WaitGroup
	s := service.InitializeService(&wg)
	h := server.InitializeHandler(s, &wg)
	r := server.NewRouter(h)
	r.InitializeRoutes()

	server := http.Server{
		Addr:    ":8080",
		Handler: r,
	}

	termChan := make(chan os.Signal)
	signal.Notify(termChan, syscall.SIGINT)

	go func() {
		<-termChan
		log.Print("SIGTERM received. Shutting down server.")
		server.Shutdown(context.Background())
	}()

	if err := server.ListenAndServe(); err != nil {
		if err.Error() != "http: Server closed" {
			log.Printf("HTTP server closed")
		}
		log.Printf("HTTP server shut down")
	}
}
