package client

const StatusRegistered = "REGISTERED"
const StatusProcessed = "PROCESSED"
const StatusInvalid = "INVALID"
const StatusProcessing = "PROCESSING"

type AccrualServiceResponse struct {
	Order   string  `json:"order"`
	Status  string  `json:"status"`
	Accrual float64 `json:"accrual"`
}
