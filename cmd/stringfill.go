/*
Copyright © 2023 Chris Mague github@mague.com
*/
package cmd

import (
	"context"
	"log"

	"github.com/maguec/RedisSim/simredis"
	"github.com/maguec/RedisSim/utils"
	"github.com/spf13/cobra"
)

var ctx = context.Background()
var size, totalSize int
var statsHide bool
var prefix string

// stringfillCmd represents the stringfill command
var stringfillCmd = &cobra.Command{
	Use:   "stringfill",
	Short: "Add Datatype string to a specific memory size",
	Long:  `This is used to simulate a datasize being stored on a Redis server`,
	Run: func(cmd *cobra.Command, args []string) {
		conf := simredis.RedisConf(server, password, clients, port)
		cluster := simredis.ClusterClient(conf, ctx)
		err := cluster.Ping(ctx).Err()
		if err != nil {
			log.Panic("Unable to connect to cluster: ", err.Error())
		}
		err = utils.Stringfill(conf, size, totalSize, clients, rps, statsHide, prefix)
		if err != nil {
			log.Panic("Unable to connect to cluster: ", err.Error())
		}

	},
}

func init() {
	rootCmd.AddCommand(stringfillCmd)

	stringfillCmd.Flags().IntVar(&size, "size", 32, "size in bytes per record")
	stringfillCmd.Flags().IntVar(&totalSize, "string-count", 1000, "total size of  records in memory")
	stringfillCmd.Flags().BoolVarP(&statsHide, "stats-hide", "x", false, "Hide statistics")
	stringfillCmd.Flags().StringVar(&prefix, prefix, "string", "Prefix all keys with this string:")
}
