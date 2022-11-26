package cmd

import (
	"github.com/spf13/cobra"
)

func init() {
}

var DBCmd = &cobra.Command{
	Use:   "db",
	Short: "Commands related to Database",
	Long:  `Commands related to Database`,
}
