package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"sync"
	"time"

	"golang.org/x/sync/singleflight"
)

func main() {
	var wg sync.WaitGroup
	sf := singleflight.Group{}

	c := &Cache[string, string]{
		sf: &sf,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	resultChan := make(chan string)

	for range 5 {
		wg.Go(func() {
			result, err := c.Get(ctx, "1", func(ctx context.Context) (string, error) {
				return callUserAPI("1")
			})
			if err != nil {
				slog.Error("receieved error for 1", "error", err)
				resultChan <- ""
			}

			resultChan <- result
		})
	}

	for range 5 {
		wg.Go(func() {
			result, err := c.Get(ctx, "2", func(ctx context.Context) (string, error) {
				return callUserAPI("2")
			})
			if err != nil {
				slog.Error("receieved error for 2", "error", err)
				resultChan <- ""
			}

			resultChan <- result
		})
	}

	go func() {
		wg.Wait()
		close(resultChan)
	}()

	for result := range resultChan {
		fmt.Println(result)
	}

}

type Cache[K comparable, V any] struct {
	sf *singleflight.Group
	wg *sync.WaitGroup
}

func (c *Cache[K, V]) Get(ctx context.Context, key K, loader func(context.Context) (V, error)) (V, error) {
	doKey := fmt.Sprintf("%v", key)
	var zero V
	select {
	case result := <-c.sf.DoChan(doKey, func() (any, error) {
		return loader(ctx)
	}):
		if result.Err != nil {
			return zero, result.Err
		}

		return result.Val.(V), nil
	case <-ctx.Done():
		return zero, ctx.Err()
	}

}

type user struct {
	Id   int    `json:"id"`
	Name string `json:"firstName"`
}

func callUserAPI(key string) (string, error) {
	slog.Info("calling user API", "user", key)
	time.Sleep(10 * time.Second)
	req, err := http.NewRequest("GET", "https://dummyjson.com/users/"+key, nil)
	if err != nil {
		return "", fmt.Errorf("creating http request: %w", err)
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("making http request: %w", err)
	}

	resBytes, err := io.ReadAll(res.Body)
	if err != nil {
		return "", fmt.Errorf("reading http response: %w", err)
	}

	var user user
	err = json.Unmarshal(resBytes, &user)
	if err != nil {
		return "", fmt.Errorf("marshalling response: %w", err)
	}

	return user.Name, nil
}
