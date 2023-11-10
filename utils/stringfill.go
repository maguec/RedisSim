package utils

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/go-redis/redis/v9"
	"github.com/jamiealquiza/tachymeter"
	"github.com/maguec/RedisSim/simredis"
	"github.com/maguec/metermaid"
)

func fillworker(
	id int, wg *sync.WaitGroup, conf *redis.ClusterOptions,
	jobs <-chan int,
	ctx context.Context, size int,
	mm *metermaid.Metermaid, tach *tachymeter.Tachymeter,
) {
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
			tach.AddTime(time.Since(startTime))
			mm.Add()
		default:
			return
		}
	}
}

func Stringfill(conf *redis.ClusterOptions, size, count, threads int, hide bool) error {
	var ctx = context.Background()
	txns := make(chan int, count)
	for t := 0; t < count; t++ {
		txns <- t
	}

	tach := tachymeter.New(&tachymeter.Config{Size: count})
	mm := metermaid.New(&metermaid.Config{Size: count})
	wg := new(sync.WaitGroup)
	for w := 0; w < threads; w++ {
		wg.Add(1)
		go fillworker(w, wg, conf, txns, ctx, size, mm, tach)
	}
	wg.Wait()
	if !hide {
		ShowStats(tach, mm)
	}
	return nil
}
