package retry

import (
	"context"
	"errors"
	"fmt"
	"math"
	"time"

	"github.com/ankitsalunkhe/go-kata/04-errors-semantics/08-retry-backoff-policy/client"
)

type Retryer struct {
	MaxAttempts int
	Delay       time.Duration
}

func (r *Retryer) Do(ctx context.Context, fn func(context.Context) error) error {
	err := fn(ctx)
	if err == nil {
		return nil
	}

	var errTransient *client.ErrTransient
	if !errors.As(err, &errTransient) {
		return err
	}

	timer := time.NewTimer(r.Delay)

	for i := 1; i <= r.MaxAttempts; i++ {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-timer.C:
			delay := r.Delay * time.Duration(math.Pow(2, float64(i)))
			timer.Reset(delay)
			err := fn(ctx)
			if err == nil {
				return nil
			}

			if !errors.As(err, &errTransient) {
				return err
			}
		}
	}

	return fmt.Errorf("exhausted all %d retries: %w", r.MaxAttempts, err)
}
