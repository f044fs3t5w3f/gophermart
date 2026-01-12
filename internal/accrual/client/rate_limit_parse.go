package client

import (
	"fmt"
	"regexp"
	"strconv"
)

const rateLimitPattern = "No more than (\\d+) requests per minute allowed"

func extractRateLimit(body string) (int, error) {
	match := regexp.MustCompile(rateLimitPattern).FindStringSubmatch(body)
	if len(match) == 0 {
		return 0, fmt.Errorf("rate limit not found")
	}
	rateLimit, err := strconv.Atoi(match[1])
	if err != nil {
		return 0, err
	}
	return rateLimit, nil
}
