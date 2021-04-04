package handlers

import (
	"crypto/sha512"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
	"time"
)

var hashedPasswords []string
var durations []time.Duration

const password = "password"

type Handler struct {}

type Stats struct {
	Total   int64 `json:"total"`
	Average int64 `json:"average"`
}

func InitializeHandler() *Handler {
	return &Handler{}
}

func (h *Handler) CreateHash(w http.ResponseWriter, r *http.Request) {
	startTime := time.Now()
	fiveSecTimer := time.NewTimer(5 * time.Second)

	if err := r.ParseForm(); err != nil {
		fmt.Errorf("500 failed to parse form - this could probably be a better error")
	}
	password := r.PostForm.Get(password)

	go func() {
		<-fiveSecTimer.C
		h.hashPassword(password)
		duration := time.Since(startTime)
		durations = append(durations, duration)
	}()

	json.NewEncoder(w).Encode(len(hashedPasswords) + 1)
}

func (h *Handler) GetHash(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	intID, err := strconv.Atoi(id)
	if err != nil {
		fmt.Errorf("failed to convert path param to int")
	}

	json.NewEncoder(w).Encode(hashedPasswords[intID-1])
}

func (h *Handler) GetStats(w http.ResponseWriter, r *http.Request) {
	totalDuration := time.Duration(0)
	for _, duration := range durations {
		totalDuration = totalDuration + duration
	}

	averageDuration := int(totalDuration*time.Microsecond) / len(hashedPasswords)

	stats := &Stats{
		Total:   int64(len(durations)),
		Average: int64(averageDuration),
	}
	json.NewEncoder(w).Encode(stats)
}

func (h *Handler) hashPassword(password string) () {
	pwBytes := []byte(password)
	sha := sha512.Sum512(pwBytes)
	hashedPasswords = append(hashedPasswords, base64.StdEncoding.EncodeToString(sha[:]))
}