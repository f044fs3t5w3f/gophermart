package models

import "time"

type Withdraw struct {
	ID          int64
	UserID      int64
	Order       string
	Sum         float64
	ProcessedAt time.Time
}
