package main

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"golang.org/x/time/rate"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		client := http.Client{
			Timeout: 10 * time.Second,
		}

		limiter := rate.NewLimiter(rate.Limit(1), 10)
		f := FanOutClient{
			MaxInFlight: 8,
			client:      client,
			limiter:     limiter,
		}

		res, err := f.FetchAll(context.Background(), []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12})
		w.Header().Set("Content-Type", "application/json")
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(err.Error())
			return
		}
		json.NewEncoder(w).Encode(res)
	})

	http.ListenAndServe(":8081", nil)
}
