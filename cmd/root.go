/*
Copyright © 2023 Chris Mague github@mague.com

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

	http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var server, password, prefix, ratio string
var port, clients, rps int
var verbose, cluster bool

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "RedisSim",
	Short: "Simulate Load on your Redis server",
	Long: `Simulate various load scenarios on your Redis server

See the various sub commands for options`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	//Run: func(cmd *cobra.Command, args []string) { },
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVar(&server, "server", "localhost", "Redis Server to connect to")
	rootCmd.PersistentFlags().IntVar(&port, "port", 6379, "Redis Port to connect to")
	rootCmd.PersistentFlags().StringVar(&password, "password", "", "Redis password to connect with")
	rootCmd.PersistentFlags().IntVar(&clients, "clients", 10, "Number of clients to use")
	rootCmd.PersistentFlags().IntVar(&rps, "rps", 0, "Rate limit for number of requests per second (0 is disabled)")
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "Verbose output")
	rootCmd.PersistentFlags().BoolVarP(&cluster, "cluster", "c", false, "Enable Cluser API")

}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	viper.AutomaticEnv() // read in environment variables that match
}
