/*
Copyright © 2023 Chris Mague github@mague.com
*/
package cmd

import (
	"log"

	"github.com/maguec/RedisSim/simredis"
	"github.com/maguec/RedisSim/utils"
	"github.com/spf13/cobra"
)

var cpukillruns, cpukillcount, cpukillsize int
var loop bool

// cpukillCmd represents the cpukill command
var cpukillCmd = &cobra.Command{
	Use:   "cpukill",
	Short: "Run heavy CPU operations",
	Long:  `This command simulates a heavy CPU load on a redis server.  Creates a number of keys with a lot of elements then attempts to find the union of the two keys`,
	Run: func(cmd *cobra.Command, args []string) {
		conf := simredis.RedisConf(server, password, clients, port, cluster)
		cluster := simredis.ClusterClient(conf, ctx)
		err := cluster.Ping(ctx).Err()
		if err != nil {
			log.Panic("Unable to connect to cluster: ", err.Error())
		}
		err = utils.CPUKill(conf, cpukillcount, clients, rps, cpukillruns, cpukillsize, statsHide, prefix, loop)
		if err != nil {
			log.Panic("Unable to connect to cluster: ", err.Error())
		}

	},
}

func init() {
	rootCmd.AddCommand(cpukillCmd)
	cpukillCmd.Flags().IntVar(&cpukillruns, "runs", 100, "number of times to run through the simulation")
	cpukillCmd.Flags().IntVar(&cpukillcount, "count", 1000, "number of keys to simulate")
	cpukillCmd.Flags().IntVar(&cpukillsize, "size", 10000, "size of keys to simulate")
	cpukillCmd.Flags().StringVar(&prefix, "prefix", "CPUKILL", "Prefix all keys with this string:")
	cpukillCmd.Flags().BoolVarP(&loop, "loop-forever", "l", false, "Loop forever simulating statistics")
}
