package main

import (
	"context"
	"fmt"
	"time"
)

func main() {
	// schedul()
	// cmder()
	ctx := context.Background()
	// ctx, cancel := context.WithCancel(ctx)
	// defer cancel()
	ctx, cancel := context.WithTimeout(ctx, time.Second)
	defer cancel()

	// time.AfterFunc(time.Second, cancel)

	sleepAndTake(ctx, 5*time.Second, "hello")
}

func sleepAndTake(ctx context.Context, d time.Duration, s string) {
	select {
	case <-time.After(d):
		fmt.Println(s)
	case <-ctx.Done():
		fmt.Println("Done")
		fmt.Println(ctx.Err())
	}
}
