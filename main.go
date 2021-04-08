package main

import (
	"net/http"
	"os"
	"os/signal"
	"password-encoder/server"
	"password-encoder/service"
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

	stop := make(chan os.Signal)
	signal.Notify(stop, os.Interrupt)

	server.ListenAndServe()
}
