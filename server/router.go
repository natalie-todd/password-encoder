package server

import (
	"github.com/gorilla/mux"
	"net/http"
)

type Router struct {
	*mux.Router
	handler *Handler
}

func NewRouter(handler *Handler) *Router {
	return &Router{
		Router:  mux.NewRouter(),
		handler: handler,
	}
}

func (r *Router) InitializeRoutes() {
	r.HandleFunc("/hash", r.handler.CreateHash).Methods(http.MethodPost)
	r.HandleFunc("/hash/{id}", r.handler.GetHash).Methods(http.MethodGet)
	r.HandleFunc("/stats", r.handler.GetStats).Methods(http.MethodGet)
	r.HandleFunc("/shutdown", r.handler.Shutdown).Methods(http.MethodPost)
}
