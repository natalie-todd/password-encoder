package service

import "time"

type Servicer interface {
	CalculateHashAndDuration(startTime time.Time, fiveSecTimer *time.Timer, password string)

	GetHashedPasswords() []string

	CalculateStats() *Stats
}

//mockgen -source=service/service_interface.go