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

func cpukillworker(
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
			_, err := client.SUnion(ctx,
				fmt.Sprintf("{%s:%d}:a", prefix, job),
				fmt.Sprintf("{%s:%d}:b", prefix, job),
			).Result()
			if err != nil {
				fmt.Printf(err.Error())
			}
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

func CPUKill(conf *redis.UniversalOptions, count, threads, rps, runs, keylength int, hide bool, prefix string) error {
	var ctx = context.Background()
	var ops []int
	client := simredis.ClusterClient(conf, ctx)
	log.Printf("Seeding data")
	for i := 0; i < count/2; i++ {
		pipe := client.Pipeline()
		for j := 0; j < keylength; j++ {
			pipe.SAdd(ctx, fmt.Sprintf("{%s:%d}:a", prefix, i), j)
			pipe.SAdd(ctx, fmt.Sprintf("{%s:%d}:b", prefix, i), j)
		}
		_, err := pipe.Exec(ctx)
		if err != nil {
			return err
		}
		_, err = pipe.Exec(ctx)
		if err != nil {
			panic(err)
		}
	}
	log.Printf("Starting a run with %d keys of length %d", count, keylength)

	for r := 0; r < runs; r++ {
		// we halve then add the same number twice so we get all of the Sets added
		for o := 1; o <= count/2; o++ {
			ops = append(ops, o)
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
		go cpukillworker(w, wg, conf, jobs, ctx, rps, rl, prefix, bar, tach, mm)
	}

	wg.Wait()

	if !hide {
		ShowStats(tach, mm, true)
	}

	return nil
}
