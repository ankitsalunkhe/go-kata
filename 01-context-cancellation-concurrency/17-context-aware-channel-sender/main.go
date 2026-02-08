package main

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"runtime"
	"time"

	"golang.org/x/sync/errgroup"
)

func main() {
	df := DataFetcher{}
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	result := df.Fetch(ctx, []string{
		"https://dummyjson.com/users/1",
		"https://dummyjson.com/users/2",
		"https://dummyjson.com/users/3",
		"https://dummyjson.com/users/4",
		"https://dummyjson.com/users/5",
	})

	for r := range result {
		slog.Info(string(r.Body))
	}

	fmt.Println("count", runtime.NumGoroutine())
}

type DataFetcher struct{}

type Result struct {
	URL  string
	Body []byte
	Err  error
}

func (f *DataFetcher) Fetch(ctx context.Context, urls []string) <-chan Result {
	errgrp, ctx := errgroup.WithContext(ctx)

	resultChan := make(chan Result)

	for i, url := range urls {
		errgrp.Go(func() error {
			return singleFetch(ctx, url, i, resultChan)
		})
	}

	go func() {
		defer close(resultChan)
		if err := errgrp.Wait(); err != nil {
			fmt.Printf("Encountered an error: %v\n", err)
			r := Result{
				Err: err,
			}
			resultChan <- r
		}
	}()

	return resultChan
}

func singleFetch(ctx context.Context, url string, key int, resultChan chan<- Result) error {
	singleResultChan := make(chan Result)
	go callAPI(url, key, singleResultChan)

	select {
	case <-ctx.Done():
		return ctx.Err()
	case result := <-singleResultChan:
		resultChan <- result
		return nil
	}
}

func callAPI(url string, key int, singleResultChan chan<- Result) {
	slog.Info("calling", "id", key)
	if key == 2 {
		time.Sleep(5 * time.Second)
	}

	result := Result{}
	resByte, err := http.DefaultClient.Get(url)
	if err != nil {
		result.Err = err
		singleResultChan <- result
		return
	}
	defer resByte.Body.Close()

	res, err := io.ReadAll(resByte.Body)
	if err != nil {
		result.Err = err
		singleResultChan <- result
		return
	}

	result.Body = res
	singleResultChan <- result
}
