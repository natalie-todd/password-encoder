package server

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
	"password-encoder/service"
	"strconv"
	"time"
)

const password = "password"

type Handler struct {
	service *service.Service
}

func InitializeHandler(service *service.Service) *Handler {
	return &Handler{
		service: service,
	}
}

func (h *Handler) CreateHash(w http.ResponseWriter, r *http.Request) {
	startTime := time.Now()
	fiveSecTimer := time.NewTimer(5 * time.Second)

	if err := r.ParseForm(); err != nil {
		fmt.Errorf("500 failed to parse form - this could probably be a better error")
	}
	password := r.PostForm.Get(password)

	h.service.CalculateHashAndDuration(startTime, fiveSecTimer, password)

	json.NewEncoder(w).Encode(len(h.service.GetHashedPasswords()) + 1)
}

func (h *Handler) GetHash(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	intID, err := strconv.Atoi(id)
	if err != nil {
		fmt.Errorf("failed to convert path param to int")
	}

	json.NewEncoder(w).Encode(h.service.GetHashedPasswords()[intID-1])
}

func (h *Handler) GetStats(w http.ResponseWriter, r *http.Request) {
	stats := h.service.CalculateStats()
	json.NewEncoder(w).Encode(stats)
}

func (h *Handler) Shutdown(w http.ResponseWriter, r *http.Request) {
	fmt.Println("shutdown")
	ctx := (*r).Context()
	ctx.Done()
}
