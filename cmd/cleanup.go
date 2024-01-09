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

var cleanupRm, cleanupDry bool
var cleanupPrefix string

// cleanupCmd represents the cleanup command
var cleanupCmd = &cobra.Command{
	Use:   "cleanup",
	Short: "Cleanup keys by prefix",
	Long:  `Allows you to mass set of keys by prefix by either deleting the keys or setting a TTL.`,
	Run: func(cmd *cobra.Command, args []string) {
		// We're going to overwrite the pool size here as we spin a new connection for each go routine
		conf := simredis.RedisConf(server, password, 1, port, cluster)
		cluster := simredis.ClusterClient(conf, ctx)
		err := cluster.Ping(ctx).Err()
		if err != nil {
			log.Panic("Unable to connect to cluster: ", err.Error())
		}
		if cleanupRm {
			err = utils.KeyCleanup(conf, ctx, cleanupPrefix, cleanupDry)
			if err != nil {
				log.Panic("Unable to cleanup: ", err.Error())
			}
		} else {
			log.Println("No action is set, please use --rm to delete keys try with --dry-run first")
		}
	},
}

func init() {
	rootCmd.AddCommand(cleanupCmd)

	cleanupCmd.Flags().BoolVar(&cleanupRm, "rm", false, "Delete all keys matching a prefix")
	cleanupCmd.Flags().BoolVar(&cleanupDry, "dry-run", false, "Show what would happen but do not perform delete")
	cleanupCmd.Flags().StringVar(&cleanupPrefix, "prefix", "", "Remove keys matching this prefix")
}
