package cri

import (
	"context"
	"time"
)

// maxMsgSize use 16MB as the default message size limit.
// grpc library default is 4MB
const maxMsgSize = 1024 * 1024 * 16

// getContextWithTimeout returns a context with timeout.
func getContextWithTimeout(timeout time.Duration) (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), timeout)
}

func getContextWithCancel() (context.Context, context.CancelFunc) {
	return context.WithCancel(context.Background())
}
