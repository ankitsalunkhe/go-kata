package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"strconv"
	"time"

	"golang.org/x/sync/errgroup"
	"golang.org/x/time/rate"
)

type FanOutClient struct {
	MaxInFlight int
	client      http.Client
	limiter     *rate.Limiter
}

type user struct {
	ID        int    `json:"id"`
	FirstName string `json:"firstName"`
}

func (f FanOutClient) FetchAll(ctx context.Context, userIDs []int) (map[int][]byte, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	eg, ctx := errgroup.WithContext(ctx)
	eg.SetLimit(f.MaxInFlight)

	userChan := make(chan user, len(userIDs))

	for _, userId := range userIDs {
		eg.Go(func() error {
			return f.fetch(ctx, userId, userChan)
		})
	}

	if err := eg.Wait(); err != nil {
		return nil, fmt.Errorf("error during fetching : %v", err)
	}

	close(userChan)

	result := make(map[int][]byte, len(userIDs))
	for res := range userChan {
		userBytes, err := json.Marshal(res)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal user: %v", err)
		}
		result[res.ID] = userBytes
	}

	return result, nil
}

func (f FanOutClient) fetch(ctx context.Context, userID int, userChan chan user) error {
	slog.Info("fetching", "user", userID)

	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, "http://dummyjson.com/users/"+strconv.Itoa(userID), nil)
	if err != nil {
		return fmt.Errorf("creating http request: %v", err)
	}

	if err := f.limiter.Wait(ctx); err != nil {
		return fmt.Errorf("waiting on limiter :%v", err)
	}

	httpRes, err := f.client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to callout :%v", err)
	}

	resByte, err := io.ReadAll(httpRes.Body)
	defer httpRes.Body.Close()

	response := user{}
	if err := json.Unmarshal(resByte, &response); err != nil {
		return fmt.Errorf("failed to unmarshal response :%v", err)
	}

	userChan <- response

	return nil
}
