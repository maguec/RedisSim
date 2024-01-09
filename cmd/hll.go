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

var hllCount, hllEntryCount, hllRuns int
var hllPrefix string
var hllMerge bool

// hllCmd represents the hll command
var hllCmd = &cobra.Command{
	Use:   "hll",
	Short: "Get timing information around HyperLoglog in Redis",
	Long: `Create HyperLoglog data structues in Redis and gather timing information around them
For more information on HyperLoglog see Redis documentation.`,
	Run: func(cmd *cobra.Command, args []string) {
		conf := simredis.RedisConf(server, password, 1, port, cluster)
		cluster := simredis.ClusterClient(conf, ctx)
		err := cluster.Ping(ctx).Err()
		if err != nil {
			log.Panic("Unable to connect to cluster: ", err.Error())
		}
		err = utils.HllBench(conf, hllCount, hllEntryCount, hllRuns, rps, statsHide, hllMerge, hllPrefix)
		if err != nil {
			log.Panic("Unable to connect to cluster: ", err.Error())
		}
	},
}

func init() {
	rootCmd.AddCommand(hllCmd)

	hllCmd.Flags().BoolVar(&hllMerge, "merge-test", false, "Test merging HLL data structures")
	hllCmd.Flags().IntVar(&hllCount, "hll-count", 100, "Number of HyperLogLog keys to create")
	hllCmd.Flags().IntVar(&hllEntryCount, "hll-entry-count", 100, "Number of HyperLogLog entries per key")
	hllCmd.Flags().IntVar(&hllRuns, "hll-runs", 10, "Number of times to run through the tests")
	hllCmd.Flags().StringVar(&hllPrefix, "hll-prefix", "hll", "Prefix to use in front of the HLL")
}
