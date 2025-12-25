package client

import (
	"fmt"
	"net/http"
	"strconv"
)

func extractRetryAfter(header http.Header) (int, error) {
	retryAfter := header.Get("Retry-After")
	if retryAfter == "" {
		return 0, fmt.Errorf("no Retry-After header")
	}
	retryAfterInt, err := strconv.Atoi(retryAfter)
	if err == nil {
		return retryAfterInt, nil
	} else {
		return 0, fmt.Errorf("invalid Retry-After header: %w", err)
	}
}
