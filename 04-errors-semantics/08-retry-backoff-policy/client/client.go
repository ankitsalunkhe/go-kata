package client

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"slices"
	"strconv"
	"time"
)

type ErrTransient struct {
	Message string
}

func (err *ErrTransient) Error() string {
	return err.Message
}

var (
	ErrTimeout = &ErrTransient{"timeout from callout"}
	ErrTooMany = &ErrTransient{"too many, unavailable"}
)

type Api struct {
	Timeout time.Duration
	Status  int
}

type Caller interface {
	Callout(ctx context.Context) error
}

func (c Api) Callout(ctx context.Context) error {
	ctx, cancel := context.WithTimeout(ctx, c.Timeout)
	defer cancel()

	slog.Info("Start callout")
	client := http.Client{}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, "https://httpbin.org/status/"+strconv.Itoa(c.Status), nil)
	if err != nil {
		return fmt.Errorf("creating http request :%v", err)
	}

	res, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("making http request :%w", ErrTimeout)
	}

	if slices.Contains([]int{http.StatusTooManyRequests, http.StatusServiceUnavailable}, res.StatusCode) {
		return ErrTooMany
	}

	return err
}
