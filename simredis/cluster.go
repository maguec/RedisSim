package simredis

import (
	"context"

	"github.com/go-redis/redis/v9"
)

func getMasterNodes(conf *redis.ClusterOptions) ([]string, error) {
	var ctx = context.Background()
	nodes := []string{}
	client := redis.NewClusterClient(conf)
	slots, err := client.ClusterSlots(ctx).Result()
	if err != nil {
		return nil, err
	}
	for _, s := range slots {
		nodes = append(nodes, s.Nodes[0].Addr)
	}
	return nodes, nil

}
