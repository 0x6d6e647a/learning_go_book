package main

import (
	"context"
	"fmt"
	"math/rand"
	"time"
)

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	total := 0
	count := 0
	var reason string

loop:
	for {
		select {
		case <-ctx.Done():
			reason = ctx.Err().Error()
			break loop
		default:
		}

		n := rand.Intn(100_000_000)
		if n == 1_234 {
			reason = "got 1,234"
			break loop
		}
		total += n
		count += 1

	}

	fmt.Println("total:", total, "number of iterations:", count, reason)
}
