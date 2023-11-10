/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"context"
	"log"

	"github.com/maguec/RedisSim/utils"
	"github.com/spf13/cobra"
)

var ctx = context.Background()
var size, totalSize int

// stringfillCmd represents the stringfill command
var stringfillCmd = &cobra.Command{
	Use:   "stringfill",
	Short: "Add Datatype string to a specific memory size",
	Long:  `This is used to simulate a datasize being stored on a Redis server`,
	Run: func(cmd *cobra.Command, args []string) {
		cluster := utils.RedisConf(server, password, clients, port)
		err := cluster.Ping(ctx).Err()
		if err != nil {
			log.Panic("Unable to connect to cluster: ", err.Error())
		}
		err = utils.Stringfill(cluster, size, totalSize, clients)
		if err != nil {
			log.Panic("Unable to connect to cluster: ", err.Error())
		}

	},
}

func init() {
	rootCmd.AddCommand(stringfillCmd)

	stringfillCmd.Flags().IntVar(&size, "size", 10, "size in bytes per record")
	stringfillCmd.Flags().IntVar(&totalSize, "totalsize", 1000, "total size of  records in memory")
	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// stringfillCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
