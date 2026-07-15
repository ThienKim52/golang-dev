package main

import (
	"context"
	"time"

	"github.com/ThienKim52/golang-dev/pkg/redis"
)

func main() {
	rclient, err := redis.NewClient("")
	if err != nil {
		panic(err)
	}
	rclient.Set(context.Background(), "1235", "google.com", time.Hour)

	rclient2, err := redis.NewClient("CACHE")
	if err != nil {
		panic(err)
	}
	rclient2.Set(context.Background(), "1999", "google.com", time.Hour)
}
