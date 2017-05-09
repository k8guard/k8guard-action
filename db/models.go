package db

import (
	"time"
)

type VActionRow struct {
	Namespace string
	Type      string
	Source    string
	VType     string
	VSource   string
	Actions   map[string][]time.Time
	CreatedAt time.Time
	ExpiresAt time.Time
}
