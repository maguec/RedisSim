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

var runs int

// exerciseCmd represents the exercise command
var exerciseCmd = &cobra.Command{
	Use:   "exercise",
	Short: "Exercise your keys based on a prefix",
	Long:  `The exercise command grabs the key prefix:1 and attempts to determine it's type then goes ahead and updates/reads these keys to generate load so you can see how your server performs.`,
	Run: func(cmd *cobra.Command, args []string) {
		if prefix == "" {
			log.Fatal("Please set a prefix")
		}
		// We're going to overwrite the pool size here as we spin a new connection for each go routine
		conf := simredis.RedisConf(server, password, 1, port, cluster)
		cluster := simredis.ClusterClient(conf, ctx)
		err := cluster.Ping(ctx).Err()
		if err != nil {
			log.Panic("Unable to connect to cluster: ", err.Error())
		}
		err = utils.Exercise(conf, size, totalSize, clients, rps, runs, statsHide, prefix, ratio)
		if err != nil {
			log.Panic("Unable to connect to cluster: ", err.Error())
		}

	},
}

func init() {
	rootCmd.AddCommand(exerciseCmd)
	exerciseCmd.Flags().StringVar(&prefix, "prefix", "", "Prefix all keys with this string:")
	exerciseCmd.Flags().StringVar(&ratio, "ratio", "1:1", "ratio of reads to writes")
	exerciseCmd.Flags().IntVar(&totalSize, "key-count", 100, "number of keys to exercise")
	exerciseCmd.Flags().IntVar(&runs, "runs", 1, "number of times to run through the exercise")
}
