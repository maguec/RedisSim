package utils

import (
	"context"

	"github.com/go-redis/redis/v9"
	"github.com/maguec/RedisSim/simredis"
)

func getMasterNodes(conf *redis.ClusterOptions, ctx context.Context) ([]string, error) {
	nodes := []string{}
	client := simredis.ClusterClient(conf, ctx)
	slots, err := client.ClusterSlots(ctx).Result()
	if err != nil {
		return nil, err
	}
	for _, s := range slots {
		nodes = append(nodes, s.Nodes[0].Addr)
	}
	return nodes, nil

}
