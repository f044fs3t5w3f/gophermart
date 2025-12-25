package handler

import "encoding/json"

type orderResponse struct {
	Number     string  `json:"number"`
	Status     string  `json:"status"`
	Accrual    float64 `json:"accrual"`
	UploadedAt string  `json:"uploaded_at,omitempty"`
}

func (order orderResponse) MarshalJSON() ([]byte, error) {
	type orderResponseAlias orderResponse
	if order.Status == "PROCESSED" {
		return json.Marshal(&struct {
			orderResponseAlias
		}{orderResponseAlias(order)})
	}

	return json.Marshal(&struct {
		Number     string `json:"number"`
		Status     string `json:"status"`
		UploadedAt string `json:"uploaded_at,omitempty"`
	}{
		Number:     order.Number,
		Status:     order.Status,
		UploadedAt: order.UploadedAt,
	})
}
