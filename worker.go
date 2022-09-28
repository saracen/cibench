package main

import (
	"time"
)

func worker(idx int, events chan int, fn func(int) error, start time.Time) {
	go func() {
		i := 0
		for {
			i++

			if err := fn(idx); err != nil {
				break
			}

			if time.Since(start) >= 10*time.Second {
				break
			}
		}
		events <- i
	}()
}

func do(threads int, fn func(idx int) error) (int, time.Duration) {
	events := make(chan int, threads)

	start := time.Now()
	for i := 0; i < threads; i++ {
		worker(i, events, fn, start)
	}

	total := 0
	for i := 0; i < threads; i++ {
		total += <-events
	}

	return total, time.Since(start)
}
