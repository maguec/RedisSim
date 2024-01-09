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
	"github.com/schollz/progressbar/v3"
	"go.uber.org/ratelimit"
)

func hllworker(
	id int, wg *sync.WaitGroup, conf *redis.UniversalOptions,
	ctx context.Context,
	entries, runs, rps int,
	rl ratelimit.Limiter,
	addMm *metermaid.Metermaid, addTach *tachymeter.Tachymeter,
	countMm *metermaid.Metermaid, countTach *tachymeter.Tachymeter,
	prefix string, bar *progressbar.ProgressBar,
) {
	client := simredis.ClusterClient(conf, ctx)
	client.Ping(ctx)
	defer wg.Done()

	// add all of the entries
	for run := 0; run < runs; run++ {
		for entry := 0; entry < entries; entry++ {
			if rps > 0 {
				rl.Take()
			}
			bar.Add(1)
			startTime := time.Now()
			err := client.PFAdd(
				ctx,
				fmt.Sprintf("%s:%d", prefix, id),
				fmt.Sprintf("%s:%d", prefix, entry),
			).Err()
			if err != nil {
				fmt.Printf("Error with PFAdd thread %d and job %d-%d\n", id, entry, run)
			}
			addTach.AddTime(time.Since(startTime))
			addMm.Add()
		}
		bar.Add(1)
		startTime := time.Now()
		err := client.PFCount(
			ctx,
			fmt.Sprintf("%s:%d", prefix, id),
		).Err()
		if err != nil {
			fmt.Printf("Error with PFCount thread %d and job %d\n", id, run)
		}
		countTach.AddTime(time.Since(startTime))
		countMm.Add()
	}
}

func HllBench(
	conf *redis.UniversalOptions,
	hllCount, hllEntryCount, hllRuns, rps int,
	hide, mergeTest bool,
	prefix string,
) error {
	var ctx = context.Background()

	// rewrite this if the RPS > 0 since creating a ratelimiter with 0 will div by zero
	rl := ratelimit.New(1)
	if rps > 0 {
		rl = ratelimit.New(rps)
	}

	addTach := tachymeter.New(&tachymeter.Config{Size: hllCount * hllEntryCount * hllRuns})
	countTach := tachymeter.New(&tachymeter.Config{Size: hllCount * hllRuns})
	addMm := metermaid.New(&metermaid.Config{Size: hllCount * hllEntryCount * hllRuns})
	countMm := metermaid.New(&metermaid.Config{Size: hllCount * hllRuns})
	wg := new(sync.WaitGroup)
	bar := progressbar.Default(int64(hllCount*hllEntryCount*hllRuns + hllRuns*hllCount))
	for w := 0; w < hllCount; w++ {
		wg.Add(1)
		go hllworker(
			w, wg, conf, ctx,
			hllEntryCount, hllRuns, rps, rl,
			addMm, addTach, countMm, countTach,
			prefix, bar)
	}
	wg.Wait()
	if !hide {
		fmt.Println("============================ PFAdd ============================")
		ShowStats(addTach, addMm, false)
		fmt.Println("============================ PFCount ============================")
		ShowStats(countTach, countMm, false)
	}
	return nil
}
