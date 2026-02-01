package main

import (
	"context"
	"log/slog"
	"time"

	"github.com/ankitsalunkhe/go-kata/04-errors-semantics/08-retry-backoff-policy/client"
	"github.com/ankitsalunkhe/go-kata/04-errors-semantics/08-retry-backoff-policy/retry"
)

func main() {
	r := retry.Retryer{
		MaxAttempts: 3,
		Delay:       2 * time.Second,
	}

	c := client.Api{
		Status: 503,
	}

	if err := r.Do(context.Background(), c.Callout); err != nil {
		slog.Error("Error from service", "err", err)
	}
}
