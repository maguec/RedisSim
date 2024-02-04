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

var hotkeyruns, hotkeycount int

// hotkeyCmd represents the hotkey command
var hotkeyCmd = &cobra.Command{
	Use:   "hotkey",
	Short: "Simulate a hot key",
	Long:  `This command simulates a hot key on a redis server`,
	Run: func(cmd *cobra.Command, args []string) {
		conf := simredis.RedisConf(server, password, clients, port, cluster)
		cluster := simredis.ClusterClient(conf, ctx)
		err := cluster.Ping(ctx).Err()
		if err != nil {
			log.Panic("Unable to connect to cluster: ", err.Error())
		}
		err = utils.HotkeySim(conf, hotkeycount, clients, rps, hotkeyruns, statsHide, prefix)
		if err != nil {
			log.Panic("Unable to connect to cluster: ", err.Error())
		}

	},
}

func init() {
	rootCmd.AddCommand(hotkeyCmd)
	hotkeyCmd.Flags().IntVar(&hotkeyruns, "hotkey-runs", 100, "number of times to run through the simulate")
	hotkeyCmd.Flags().IntVar(&hotkeycount, "hotkey-count", 100000, "number of keys to simulate")
}
