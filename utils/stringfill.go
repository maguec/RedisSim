package utils

import (
	"fmt"
	"sync"
	"time"

	"github.com/go-redis/redis/v9"
)

func fillworker(id int, wg *sync.WaitGroup, client *redis.ClusterClient, jobs <-chan int, results chan<- time.Duration) {
	defer wg.Done()
	for {
		select {
		case job, ok := <-jobs:
			if !ok {
				fmt.Printf("Error with thread %d\n", id)
			}
			startTime := time.Now()
			fmt.Printf("do stuff here: %d - %d\n", id, job)
			time.Sleep(400 * time.Millisecond)
			results <- time.Since(startTime)
		default:
			return
		}
	}
}

func Stringfill(client *redis.ClusterClient, size, count, threads int) error {
	fmt.Printf("size: %+v\n", size)
	fmt.Printf("count: %+v\n", count)
	fmt.Printf("threads: %+v\n", threads)
	fmt.Printf("client: %+v\n", client)
	results := make(chan time.Duration, count)
	txns := make(chan int, count)
	for t := 0; t < count; t++ {
		txns <- t
	}

	wg := new(sync.WaitGroup)
	for w := 0; w < threads; w++ {
		wg.Add(1)
		go fillworker(w, wg, client, txns, results)
	}
	wg.Wait()
	return nil
}
