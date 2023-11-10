package utils

import (
	"context"
	"fmt"
	"time"

	"github.com/go-redis/redis/v9"
)

func RedisConf(server, password string, clients, port int) *redis.ClusterClient {
	var discovery_ports = []string{}
	var ctx = context.Background()
	discovery_ports = append(discovery_ports, fmt.Sprintf("%s:%d", server, port))
	// Setup Redis Connection pool
	client := redis.NewClusterClient(&redis.ClusterOptions{
		Addrs:        discovery_ports,
		Password:     password,
		PoolSize:     clients,
		MinIdleConns: clients,
		PoolTimeout:  0,
		DialTimeout:  2 * time.Second,
	})

	// update all the slots from the discovery port
	client.ReloadState(ctx)

	return client
}
