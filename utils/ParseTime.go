package utils

import (
	"time"
)

// ParseTime convierte un string en formato RFC3339 a un objeto time.Time
func ParseTime(timestamp string) time.Time {
	t, _ := time.Parse(time.RFC3339, timestamp)
	return t
}
