package models

import "time"

type OrderStatus string

const (
	OrderStatusNew        = OrderStatus("NEW")
	OrderStatusProcessing = OrderStatus("PROCESSING")
	OrderStatusInvalid    = OrderStatus("INVALID")
	OrderStatusProcessed  = OrderStatus("PROCESSED")
)

type Order struct {
	Id         int64
	UserId     int64
	Status     OrderStatus
	Number     string
	Accural    float64
	UploadedAt time.Time
}
