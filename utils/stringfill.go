package utils

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/go-redis/redis/v9"
	"github.com/maguec/RedisSim/simredis"
)

func fillworker(
	id int, wg *sync.WaitGroup, conf *redis.ClusterOptions,
	jobs <-chan int, results chan<- time.Duration,
	ctx context.Context, size int) {
	client := simredis.ClusterClient(conf, ctx)
	client.Ping(ctx)
	defer wg.Done()
	for {
		select {
		case job, ok := <-jobs:
			if !ok {
				fmt.Printf("Error with thread %d\n", id)
			}
			startTime := time.Now()
			err := client.Set(ctx, fmt.Sprintf("string:%d", job), RandStringBytes(size), 0).Err()
			if err != nil {
				fmt.Printf("Error with thread %d and job %d\n", id, job)
			}
			results <- time.Since(startTime)
		default:
			return
		}
	}
}

func Stringfill(conf *redis.ClusterOptions, size, count, threads int) error {
	var ctx = context.Background()
	results := make(chan time.Duration, count)
	txns := make(chan int, count)
	for t := 0; t < count; t++ {
		txns <- t
	}

	wg := new(sync.WaitGroup)
	for w := 0; w < threads; w++ {
		wg.Add(1)
		go fillworker(w, wg, conf, txns, results, ctx, size)
	}
	wg.Wait()
	return nil
}
