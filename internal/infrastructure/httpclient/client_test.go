package httpclient

import (
	"context"
	"errors"
	"fmt"
	"net"
	"testing"
	"time"
)

func TestIsRetryableHTTPError(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name string
		err  error
		want bool
	}{
		{name: "deadline exceeded", err: context.DeadlineExceeded, want: false},
		{name: "canceled", err: context.Canceled, want: false},
		{name: "wrapped deadline", err: fmt.Errorf("get: %w", context.DeadlineExceeded), want: false},
		{name: "net timeout", err: &net.DNSError{IsTimeout: true}, want: false},
		{name: "connection refused", err: errors.New("connection refused"), want: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			if got := isRetryableHTTPError(tt.err); got != tt.want {
				t.Fatalf("isRetryableHTTPError() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIsRetryableHTTPError_netTimeout(t *testing.T) {
	t.Parallel()
	err := &timeoutNetErr{}
	if isRetryableHTTPError(err) {
		t.Fatal("expected timeout net.Error to be non-retryable")
	}
}

type timeoutNetErr struct{}

func (e *timeoutNetErr) Error() string   { return "timeout" }
func (e *timeoutNetErr) Timeout() bool   { return true }
func (e *timeoutNetErr) Temporary() bool { return false }

func TestIsRetryableHTTPError_deadlineWithTime(t *testing.T) {
	t.Parallel()
	ctx, cancel := context.WithTimeout(context.Background(), time.Nanosecond)
	defer cancel()
	<-ctx.Done()
	if isRetryableHTTPError(ctx.Err()) {
		t.Fatal("expected context error after timeout to be non-retryable")
	}
}
