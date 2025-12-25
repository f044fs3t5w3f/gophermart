package client

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"golang.org/x/time/rate"
)

const initialRateLimit = 30

type client struct {
	limiter *rate.Limiter
	client  *http.Client
	baseURL string
}

func (c *client) GetInfo(ctx context.Context, orderID string) (*AccrualServiceResponse, error) {
	if err := c.limiter.Wait(ctx); err != nil {
		return nil, fmt.Errorf("limiter: %w", err)
	}
	url := fmt.Sprintf("http://%s/api/orders/%s", c.baseURL, orderID)
	response, err := c.client.Get(url)
	if err != nil {
		return nil, fmt.Errorf("GetInfo: %w", err)
	}
	defer response.Body.Close()
	if response.StatusCode == http.StatusTooManyRequests {
		retryAfter, err := extractRetryAfter(response.Header)
		if err != nil {
			fmt.Println(retryAfter)
			// set retry after
		}
		body, err := io.ReadAll(response.Body)
		if err != nil {
			rateLimit, err := extractRateLimit(string(body))
			if err != nil {
				c.limiter.SetLimit(rate.Limit(rateLimit))
				c.limiter.SetBurst(rateLimit)
			}
		}
		return nil, fmt.Errorf("rate limit error")
	}
	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("request status: %s", response.Status)
	}
	decoder := json.NewDecoder(response.Body)
	accrual := &AccrualServiceResponse{}
	err = decoder.Decode(accrual)
	if err != nil {
		return nil, fmt.Errorf("decode error: %s", response.Status)
	}
	return accrual, nil
}

func NewAccrualClient(url string) *client {
	return &client{
		client: &http.Client{
			Transport: &http.Transport{
				MaxIdleConns:    10,
				IdleConnTimeout: 90 * time.Second,
			},
			Timeout: 5 * time.Second,
		},
		baseURL: url,
		limiter: rate.NewLimiter(initialRateLimit, initialRateLimit),
	}
}
