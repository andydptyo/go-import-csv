package cmd

import (
	"github.com/spf13/cobra"
)

func init() {
}

var CSVCmd = &cobra.Command{
	Use:   "csv",
	Short: "Commands related to Csv",
	Long:  `Commands related to Csv`,
}
