package simredis

import (
	"context"
	"fmt"
	"time"

	"github.com/go-redis/redis/v9"
)

func RedisConf(server, password string, clients, port int, cluster bool) *redis.UniversalOptions {
	// Note: the Universal options determines if we use cluster or standalone options
	// by the lenghty of the Addrs - > 1 = enable cluster mode
	// so here we go ahead and pull the total number of nodes ahead of time to
	// ensure the correct setting
	var discovery_ports = []string{}
	if cluster {
		discovery_ports, _ = getMasterNodes(&redis.ClusterOptions{
			Addrs:    []string{fmt.Sprintf("%s:%d", server, port)},
			Password: password,
			PoolSize: 1,
		})

	} else {
		discovery_ports = append(discovery_ports, fmt.Sprintf("%s:%d", server, port))
	}
	conf := &redis.UniversalOptions{
		Addrs:        discovery_ports,
		Password:     password,
		PoolSize:     clients,
		MinIdleConns: clients,
		PoolTimeout:  0,
		DialTimeout:  2 * time.Second,
	}

	return conf
}

func ClusterClient(conf *redis.UniversalOptions, ctx context.Context) redis.UniversalClient {
	// Setup Redis Connection pool
	client := redis.NewUniversalClient(conf)
	// update all the slots from the discovery port
	//client.ReloadState(ctx)

	return client
}
