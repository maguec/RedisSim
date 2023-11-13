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

var csvfile, keyfield, xprefix string
var batchsize int
var csvstatsHide bool

// loadcsvCmd represents the loadcsv command
var loadcsvCmd = &cobra.Command{
	Use:   "loadcsv",
	Short: "Load a CSV file into Redis",
	Long: `Takes a CSV file and loads it into a Redis hash datas structure.
User needs to define the CSV field to load in as the key.`,
	Run: func(cmd *cobra.Command, args []string) {
		conf := simredis.RedisConf(server, password, clients, port)
		cluster := simredis.ClusterClient(conf, ctx)
		err := cluster.Ping(ctx).Err()
		if err != nil {
			log.Panic("Unable to connect to cluster: ", err.Error())
		}
		err = utils.CSVwrite(conf, ctx, clients, rps, csvfile, keyfield, xprefix, csvstatsHide)
		if err != nil {
			log.Panic("Unable to connect to cluster: ", err.Error())
		}
	},
}

func init() {
	rootCmd.AddCommand(loadcsvCmd)
	loadcsvCmd.Flags().StringVar(&csvfile, "csv-file", "", "CSV file to load into Redis")
	loadcsvCmd.Flags().StringVar(&keyfield, "key-field", "", "CSV field to set as the key name")
	loadcsvCmd.Flags().StringVar(&xprefix, "csv-prefix", "", "prefix the key name")
	loadcsvCmd.Flags().IntVar(&batchsize, "batch-size", 300, "Size of write batches")
	loadcsvCmd.Flags().BoolVarP(&csvstatsHide, "stats-hide", "x", false, "Hide statistics")
}
