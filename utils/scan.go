package utils

import (
	"context"
	"fmt"
	"log"
	"sync"

	"github.com/go-redis/redis/v9"
	"github.com/maguec/RedisSim/simredis"
	"github.com/schollz/progressbar/v3"
)

func KeyCleanup(conf *redis.ClusterOptions, ctx context.Context, prefix string, dryrun bool) error {
	batch := 1
	nodes, err := getMasterNodes(conf, ctx)
	if err != nil {
		return err
	}
	keys, err := scanNodes(nodes, ctx, prefix)
	if err != nil {
		return err
	}
	bar := progressbar.Default(int64(len(keys)))
	client := simredis.ClusterClient(conf, ctx)
	pipe := client.Pipeline()
	for _, key := range keys {
		bar.Add(1)
		if dryrun {
			fmt.Printf("UNLINK %s\n", key)
		} else {
			_, e := pipe.Unlink(ctx, key).Result()
			if e != nil {
				return e
			}
		}

		if batch%300 == 0 {
			_, err = pipe.Exec(ctx)
			if err != nil {
				return err
			}
		}
		batch += 1
	}
	_, err = pipe.Exec(ctx)
	if err != nil {
		return err
	}
	return nil

}

func scanNodes(nodes []string, ctx context.Context, prefix string) ([]string, error) {
	var keys []string
	wg := new(sync.WaitGroup)
	txns := make(chan []string, len(nodes))
	for w := 0; w < len(nodes); w++ {
		wg.Add(1)
		client := redis.NewClient(&redis.Options{Addr: nodes[w]})
		go scanPrefix(client, ctx, prefix, wg, txns)
	}
	wg.Wait()
	for j := 0; j < len(nodes); j++ {
		w := <-txns
		for _, k := range w {
			keys = append(keys, k)
		}
	}
	return keys, nil

}

func scanPrefix(client *redis.Client, ctx context.Context, prefix string, wg *sync.WaitGroup, results chan<- []string) {
	var k, keys []string
	var cursor uint64
	var err error
	defer wg.Done()
	pre := fmt.Sprintf("%s:*", prefix)
	for {
		k, cursor, err = client.Scan(ctx, cursor, pre, 0).Result()
		if err != nil {
			log.Fatal(err)
		}

		keys = append(keys, k...)

		if cursor == 0 { // no more keys
			break
		}
	}
	results <- keys

	return

}
