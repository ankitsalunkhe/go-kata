package retry_test

import (
	"context"
	"testing"
	"time"

	"github.com/ankitsalunkhe/go-kata/04-errors-semantics/08-retry-backoff-policy/client"
	"github.com/ankitsalunkhe/go-kata/04-errors-semantics/08-retry-backoff-policy/retry"
	"github.com/stretchr/testify/require"
)

func TestDo(t *testing.T) {

	t.Run("success", func(t *testing.T) {
		r := retry.Retryer{
			MaxAttempts: 3,
			Delay:       2 * time.Second,
		}

		mockedFunction := func(context.Context) error {
			return nil
		}

		err := r.Do(context.Background(), mockedFunction)
		require.NoError(t, err)
	})

	t.Run("retry untill retries exahausted", func(t *testing.T) {
		callCount := 0
		r := retry.Retryer{
			MaxAttempts: 3,
			Delay:       1 * time.Second,
		}

		mockedFunction := func(context.Context) error {
			callCount++
			return client.ErrTimeout
		}

		err := r.Do(context.Background(), mockedFunction)
		require.Error(t, err)
		require.EqualError(t, err, "exhausted all 3 retries: timeout from callout")
		require.Equal(t, 4, callCount)
	})

	t.Run("retry untill succes before exhausted", func(t *testing.T) {
		callCount := 0
		r := retry.Retryer{
			MaxAttempts: 3,
			Delay:       1 * time.Second,
		}

		mockedFunction := func(context.Context) error {
			if callCount == 2 {
				return nil
			}
			callCount++
			return client.ErrTimeout
		}

		err := r.Do(context.Background(), mockedFunction)
		require.NoError(t, err)
		require.Equal(t, 2, callCount)
	})
}
