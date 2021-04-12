package service

import (
	"crypto/sha512"
	"encoding/base64"
	"time"
)

var durations []time.Duration
var hashedPasswords []string

type Stats struct {
	Total   int64 `json:"total"`
	Average int64 `json:"average"`
}

type Service struct {}

func InitializeService() *Service {
	return &Service{}
}

func (s *Service) CalculateHashAndDuration(startTime time.Time, fiveSecTimer *time.Timer, password string) {
	go func() {
		<-fiveSecTimer.C
		s.hashPassword(password)
		duration := time.Since(startTime)
		durations = append(durations, duration)
	}()
}

// TODO: consider putting this method on the passwords after moving them to a struct
func (s *Service) GetHashedPasswords() []string {
	return hashedPasswords
}

func (s *Service) CalculateStats() *Stats {
	totalDuration := time.Duration(0)
	for _, duration := range durations {
		totalDuration = totalDuration + duration
	}

	averageDuration := int(totalDuration*time.Microsecond) / len(s.GetHashedPasswords())

	return &Stats{
		Total:   int64(len(durations)),
		Average: int64(averageDuration),
	}
}

func (s *Service) hashPassword(password string) () {
	pwBytes := []byte(password)
	sha := sha512.Sum512(pwBytes)
	hashedPasswords = append(hashedPasswords, base64.StdEncoding.EncodeToString(sha[:]))
}
