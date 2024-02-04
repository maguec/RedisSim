package utils

import (
	"context"
	"fmt"
	"log"
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

func simhotkey(
	id int, wg *sync.WaitGroup, conf *redis.UniversalOptions,
	jobs <-chan int,
	ctx context.Context, rps int,
	rl ratelimit.Limiter,
	prefix string,
	bar *progressbar.ProgressBar,
	tach *tachymeter.Tachymeter,
	mm *metermaid.Metermaid,
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
			if rps > 0 {
				rl.Take()
			}
			startTime := time.Now()
			_, err := client.Set(ctx,
				fmt.Sprintf("%s:%d", prefix, job),
				0, 0,
			).Result()
			if err != nil {
				fmt.Printf(err.Error())
			}
			_, err = client.Incr(ctx,
				fmt.Sprintf("%s:HOTKEY", prefix),
			).Result()
			tach.AddTime(time.Since(startTime))
			mm.Add()
			bar.Add(1)
			if err != nil {
				fmt.Printf(err.Error())
			}
		default:
			return
		}
	}
}

func HotkeySim(conf *redis.UniversalOptions, count, threads, rps, runs int, hide bool, prefix string) error {
	var ctx = context.Background()
	var ops []int
	client := simredis.ClusterClient(conf, ctx)
	// init the hotkey counter
	_, err := client.Set(ctx,
		fmt.Sprintf("%s:HOTKEY", prefix),
		0, 0,
	).Result()
	if err != nil {
		fmt.Printf("%+v\n", err)
	}
	log.Printf("Starting a run with hotkey: %s:HOTKEY, count: %d, runs: %d", prefix, count, runs)

	for r := 0; r < runs; r++ {
		for o := 1; o <= count; o++ {
			ops = append(ops, o)
		}
	}
	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(ops), func(i, j int) { ops[i], ops[j] = ops[j], ops[i] })
	jobs := make(chan int, len(ops))
	for _, x := range ops {
		jobs <- x
	}

	bar := progressbar.Default(int64(len(ops)))

	// rewrite this if the RPS > 0 since creating a ratelimiter with 0 will div by zero
	rl := ratelimit.New(1)
	if rps > 0 {
		rl = ratelimit.New(rps)
	}

	tach := tachymeter.New(&tachymeter.Config{Size: len(ops)})
	mm := metermaid.New(&metermaid.Config{Size: len(ops)})

	wg := new(sync.WaitGroup)
	for w := 0; w < threads; w++ {
		wg.Add(1)
		go simhotkey(w, wg, conf, jobs, ctx, rps, rl, prefix, bar, tach, mm)
	}
	wg.Wait()
	if !hide {
		ShowStats(tach, mm, true)
	}
	return nil
}
