/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// stringfillCmd represents the stringfill command
var stringfillCmd = &cobra.Command{
	Use:   "stringfill",
	Short: "Add Datatype string to a specific memory size",
	Long:  `This is used to simulate a datasize being stored on a Redis server`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("stringfill called")
	},
}

func init() {
	rootCmd.AddCommand(stringfillCmd)

	stringfillCmd.Flags().String("size", "10", "size in bytes per record")
	stringfillCmd.Flags().String("total-size", "1000", "total size of  records in memory")
	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// stringfillCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
