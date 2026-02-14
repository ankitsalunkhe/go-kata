package main

import (
	"context"
	"io"
	"log/slog"
	"net/http"
	"sync"
	"time"
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
		"https://httpbin.org/delay/5",
	})

	for r := range result {
		if r.Err != nil {
			slog.Error(r.Err.Error())
		} else {
			slog.Info(string(r.URL))
		}
	}
}

type DataFetcher struct{}

type Result struct {
	URL  string
	Body []byte
	Err  error
}

func (f *DataFetcher) Fetch(ctx context.Context, urls []string) <-chan Result {
	resultChan := make(chan Result, len(urls))

	wg := sync.WaitGroup{}

	for i, url := range urls {
		func(idx int, u string) {
			wg.Go(func() {
				res, err := callAPI(ctx, u, idx)
				if err != nil {
					res.Err = err
				}
				resultChan <- res
			})
		}(i, url)
	}

	go func() {
		wg.Wait()
		close(resultChan)
	}()

	return resultChan
}

func callAPI(ctx context.Context, url string, key int) (Result, error) {
	slog.Info("calling", "id", key)
	result := Result{
		URL: url,
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return result, err
	}

	resByte, err := http.DefaultClient.Do(req)
	if err != nil {
		return result, err
	}
	defer resByte.Body.Close()

	res, err := io.ReadAll(resByte.Body)
	if err != nil {
		return result, err
	}

	result.Body = res
	return result, nil
}
