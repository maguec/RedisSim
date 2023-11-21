package utils

import (
	"context"
	"errors"
	"fmt"
	"log"
	"math/rand"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/go-redis/redis/v9"
	"github.com/jamiealquiza/tachymeter"
	"github.com/maguec/RedisSim/simredis"
	"github.com/maguec/metermaid"
	"go.uber.org/ratelimit"
)

type exercisetype struct {
	keyname   string
	operation string
}

type exerciseStats struct {
	tach *tachymeter.Tachymeter
	mm   *metermaid.Metermaid
}

func exercisekey(
	ctx context.Context,
	client *redis.ClusterClient,
	job exercisetype,
) error {
	switch job.operation {
	case "read":
		err := client.Get(ctx, job.keyname).Err()
		if err != nil {
			log.Fatal(err)
		}
	case "write":
		//TODO: maybe writing something better would be a good idea here
		err := client.Set(ctx, job.keyname, "EXERCISERUN", 0).Err()
		if err != nil {
			log.Fatal(err)
		}
	default:
		return errors.New(fmt.Sprintf("I don't know what to do with key: %s  operation: %s\n", job.keyname, job.operation))
	}
	return nil
}

func exerciseworker(
	id int, wg *sync.WaitGroup, conf *redis.ClusterOptions,
	jobs <-chan exercisetype,
	ctx context.Context, rps int,
	rl ratelimit.Limiter,
	stats map[string]exerciseStats,
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
			err := exercisekey(ctx, client, job)
			if err != nil {
				fmt.Printf(err.Error())
			}
			stats[job.operation].tach.AddTime(time.Since(startTime))
			stats[job.operation].mm.Add()
		default:
			return
		}
	}
}

func ratio2rw(ratio string) []int {
	var res []int
	rs := strings.Split(ratio, ":")
	for _, s := range rs {
		i, err := strconv.Atoi(s)
		if err != nil {
			log.Fatal("Couldn't parse the ratio")
		}
		res = append(res, i)
	}
	return res
}

func Exercise(conf *redis.ClusterOptions, size, count, threads, rps, runs int, hide bool, prefix, ratio string) error {
	var ctx = context.Background()
	var ops []exercisetype
	client := simredis.ClusterClient(conf, ctx)
	keytype, err := client.Type(ctx,
		fmt.Sprintf("%s:1", prefix),
	).Result()
	if err != nil {
		fmt.Printf("%+v\n", err)
	}
	log.Printf("Starting a run for type: %s with prefix: %s", keytype, prefix)
	log.Printf("runs: %d size: %d ratio:%s", runs, count, ratio)
	rw := ratio2rw(ratio)
	for r := 0; r < runs; r++ {
		for o := 1; o <= count*rw[0]; o++ {
			ops = append(ops,
				exercisetype{operation: "read", keyname: fmt.Sprintf("%s:%d", prefix, o)},
			)
		}
		for x := 1; x <= count*rw[1]; x++ {
			ops = append(ops,
				exercisetype{operation: "write", keyname: fmt.Sprintf("%s:%d", prefix, x)},
			)
		}
	}

	// Shuffle the jobs for some random fun
	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(ops), func(i, j int) { ops[i], ops[j] = ops[j], ops[i] })

	txns := make(chan exercisetype, len(ops))
	for t := 0; t < len(ops); t++ {
		txns <- ops[t]
	}

	// rewrite this if the RPS > 0 since creating a ratelimiter with 0 will div by zero
	rl := ratelimit.New(1)
	if rps > 0 {
		rl = ratelimit.New(rps)
	}

	var stats = map[string]exerciseStats{
		"read": {
			tach: tachymeter.New(&tachymeter.Config{Size: count * runs * rw[0]}),
			mm:   metermaid.New(&metermaid.Config{Size: count * runs * rw[0]}),
		},
		"write": {
			tach: tachymeter.New(&tachymeter.Config{Size: count * runs * rw[1]}),
			mm:   metermaid.New(&metermaid.Config{Size: count * runs * rw[1]}),
		},
	}

	//tach := tachymeter.New(&tachymeter.Config{Size: count})
	//mm := metermaid.New(&metermaid.Config{Size: count})
	wg := new(sync.WaitGroup)
	for w := 0; w < threads; w++ {
		wg.Add(1)
		//	go exerciseworker(w, wg, conf, txns, ctx, size, rps, rl, mm, tach, prefix)
		go exerciseworker(w, wg, conf, txns, ctx, rps, rl, stats)
	}
	wg.Wait()
	if !hide {
		for k, v := range stats {
			if v.tach.Size > 0 {
				fmt.Printf("========================== %s ==========================\n", k)
				ShowStats(v.tach, v.mm, true)
			}
		}
	}
	return nil
}
