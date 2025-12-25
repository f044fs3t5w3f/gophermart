package models

import "time"

type Withdraw struct {
	Id          int64
	UserId      int64
	Order       string
	Sum         float64
	ProcessedAt time.Time
}
