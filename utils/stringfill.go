package utils

import (
	"context"
	"fmt"
	"math/rand"
	"sync"
	"time"

	"github.com/go-redis/redis/v9"
	"github.com/jamiealquiza/tachymeter"
	"github.com/maguec/RedisSim/simredis"
	"github.com/maguec/metermaid"
	"github.com/schollz/progressbar/v3"
	"go.uber.org/ratelimit"
)

func fillworker(
	id int, wg *sync.WaitGroup, conf *redis.ClusterOptions,
	jobs <-chan int,
	ctx context.Context, size, rps, minttl, maxttl int,
	rl ratelimit.Limiter,
	mm *metermaid.Metermaid, tach *tachymeter.Tachymeter,
	prefix string, bar *progressbar.ProgressBar,
) {
	client := simredis.ClusterClient(conf, ctx)
	client.Ping(ctx)
	defer wg.Done()
	for {
		if rps > 0 {
			rl.Take()
		}
		select {
		case job, ok := <-jobs:
			if !ok {
				fmt.Printf("Error with thread %d\n", id)
			}
			ttl := 0
			if maxttl > minttl {
				ttl = rand.Intn(maxttl-minttl) + minttl
			}
			bar.Add(1)
			startTime := time.Now()
			err := client.Set(
				ctx,
				fmt.Sprintf("%s:%d", prefix, job),
				RandStringBytes(size),
				time.Duration(ttl)*time.Second,
			).Err()
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

func Stringfill(
	conf *redis.ClusterOptions,
	size, count, threads, rps, minttl, maxttl int,
	hide bool,
	prefix string,
) error {
	var ctx = context.Background()
	txns := make(chan int, count)
	for t := 0; t < count; t++ {
		txns <- t
	}

	// rewrite this if the RPS > 0 since creating a ratelimiter with 0 will div by zero
	rl := ratelimit.New(1)
	if rps > 0 {
		rl = ratelimit.New(rps)
	}

	tach := tachymeter.New(&tachymeter.Config{Size: count})
	mm := metermaid.New(&metermaid.Config{Size: count})
	wg := new(sync.WaitGroup)
	bar := progressbar.Default(int64(count))
	for w := 0; w < threads; w++ {
		wg.Add(1)
		go fillworker(w, wg, conf, txns, ctx, size, rps, minttl, maxttl, rl, mm, tach, prefix, bar)
	}
	wg.Wait()
	if !hide {
		ShowStats(tach, mm, true)
	}
	return nil
}
