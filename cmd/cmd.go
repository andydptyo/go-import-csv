package cmd

import (
	"github.com/spf13/cobra"
)

var Cfg string
var TotalWorker int

func init() {
	RootCmd.PersistentFlags().StringVarP(&Cfg, "config", "c", "config.yaml", "Config file in yaml to be used")
	RootCmd.PersistentFlags().IntVarP(&TotalWorker, "worker", "w", 100, "total worker to be used")

	CSVCmd.AddCommand(ImportCmd)
	CSVCmd.AddCommand(PopulateCmd)
	DBCmd.AddCommand(MigrateCmd)
	RootCmd.AddCommand(DBCmd)
	RootCmd.AddCommand(CSVCmd)
}

var RootCmd = &cobra.Command{
	Use:   "app",
	Short: "app",
	Long:  `app`,
	Run: func(cmd *cobra.Command, args []string) {

	},
}
