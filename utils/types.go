package utils

import (
	"github.com/go-redis/redis/v9"
)

func RedisConf(server, password string, clients, port int) redis.ClusterClient {
	var discovery_ports = []string{}
	discovery_ports = append(discovery_ports, fmt.Sprintf("%s:%d", server, port))
	fmt.Printf("%+v\n", discovery_ports)
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
