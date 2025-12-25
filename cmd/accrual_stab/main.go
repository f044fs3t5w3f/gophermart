package main

// Эта заглушка была сделана chatGPT исключительно для откладки

import (
	"encoding/json"
	"log"
	"math/rand"
	"net/http"
	"strings"
	"sync"
	"time"
)

type OrderStatus string

const (
	StatusRegistered OrderStatus = "REGISTERED"
	StatusProcessing OrderStatus = "PROCESSING"
	StatusProcessed  OrderStatus = "PROCESSED"
	StatusInvalid    OrderStatus = "INVALID"
)

type Order struct {
	Number  string
	Status  OrderStatus
	Accrual *int
}

var (
	orders = make(map[string]*Order)
	mu     sync.Mutex
)

func main() {

	http.HandleFunc("/api/orders/", getOrderHandler)

	log.Println("Stub service started on :9000")
	log.Fatal(http.ListenAndServe(":9000", nil))
}

func getOrderHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	number := strings.TrimPrefix(r.URL.Path, "/api/orders/")
	if number == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	mu.Lock()
	order, exists := orders[number]
	if !exists {
		order = &Order{
			Number: number,
			Status: StatusRegistered,
		}
		orders[number] = order
		scheduleTransitions(order)
	}
	mu.Unlock()

	if !exists {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	writeOrderResponse(w, order)
}

func scheduleTransitions(order *Order) {
	firstDelay := time.Duration(rand.Intn(6)+5) * time.Second   // 5–10
	secondDelay := time.Duration(rand.Intn(6)+10) * time.Second // 10–15

	time.AfterFunc(firstDelay, func() {
		mu.Lock()
		if order.Status == StatusRegistered {
			order.Status = StatusProcessing
		}
		mu.Unlock()
	})

	time.AfterFunc(firstDelay+secondDelay, func() {
		mu.Lock()
		defer mu.Unlock()

		if order.Status != StatusProcessing {
			return
		}

		if rand.Intn(5) == 0 {
			order.Status = StatusInvalid
			order.Accrual = nil
			return
		}

		accrual := rand.Intn(4901) + 100 // 100–5000
		order.Status = StatusProcessed
		order.Accrual = &accrual
	})
}

func writeOrderResponse(w http.ResponseWriter, order *Order) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	resp := map[string]interface{}{
		"order":  order.Number,
		"status": order.Status,
	}

	if order.Status == StatusProcessed && order.Accrual != nil {
		resp["accrual"] = *order.Accrual
	}

	_ = json.NewEncoder(w).Encode(resp)
}
