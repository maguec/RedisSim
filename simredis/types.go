package simredis

import (
	"context"
	"fmt"
	"time"

	"github.com/go-redis/redis/v9"
)

func RedisConf(server, password string, clients, port int) *redis.ClusterOptions {
	var discovery_ports = []string{}
	discovery_ports = append(discovery_ports, fmt.Sprintf("%s:%d", server, port))
	conf := &redis.ClusterOptions{
		Addrs:        discovery_ports,
		Password:     password,
		PoolSize:     clients,
		MinIdleConns: clients,
		PoolTimeout:  0,
		DialTimeout:  2 * time.Second,
	}

	return conf
}

func ClusterClient(conf *redis.ClusterOptions, ctx context.Context) *redis.ClusterClient {
	// Setup Redis Connection pool
	client := redis.NewClusterClient(conf)
	// update all the slots from the discovery port
	client.ReloadState(ctx)

	return client
}
