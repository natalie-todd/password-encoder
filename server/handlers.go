package server

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"net/http"
	"password-encoder/service"
	"strconv"
	"syscall"
	"time"
)

const password = "password"

type Handler struct {
	service service.Servicer
}

func InitializeHandler(service service.Servicer) *Handler {
	return &Handler{
		service: service,
	}
}

func (h *Handler) CreateHash(w http.ResponseWriter, r *http.Request) {
	startTime := time.Now()
	fiveSecTimer := time.NewTimer(5 * time.Second)

	if err := r.ParseForm(); err != nil {
		http.Error(w, "failed to parse form field", http.StatusInternalServerError)
		return
	}
	password := r.PostForm.Get(password)
	if password == "" {
		http.Error(w, "empty password", http.StatusBadRequest)
		return
	}

	h.service.CalculateHashAndDuration(startTime, fiveSecTimer, password)

	if err := json.NewEncoder(w).Encode(len(h.service.GetHashedPasswords()) + 1); err != nil {
		http.Error(w, "failed to encode json response", http.StatusInternalServerError)
		return
	}
}

func (h *Handler) GetHash(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	intID, err := strconv.Atoi(id)
	if err != nil {
		http.Error(w, "failed to convert ID to int", http.StatusBadRequest)
		return
	}

	if err := json.NewEncoder(w).Encode(h.service.GetHashedPasswords()[intID-1]); err != nil {
		http.Error(w, "failed to encode hashed password", http.StatusInternalServerError)
		return
	}
}

func (h *Handler) GetStats(w http.ResponseWriter, r *http.Request) {
	stats := h.service.CalculateStats()
	if err := json.NewEncoder(w).Encode(stats); err != nil {
		http.Error(w, "failed to encode stats", http.StatusInternalServerError)
		return
	}
}

func (h *Handler) Shutdown(w http.ResponseWriter, r *http.Request) {
	if err := syscall.Kill(syscall.Getpid(), syscall.SIGINT); err != nil {
		http.Error(w, "failed to execute shutdown signal", http.StatusInternalServerError)
	}
}
