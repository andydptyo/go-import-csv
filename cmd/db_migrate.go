package cmd

import (
	"log"

	"github.com/andydptyo/go-import-csv/internal/config"
	"github.com/andydptyo/go-import-csv/internal/database/mysql"
	"github.com/spf13/cobra"
)

var rollback bool
var migrationDirection int

func init() {
	MigrateCmd.Flags().BoolVarP(&rollback, "rollback", "r", false, "Rollback database schema to previous versions")
}

var MigrateCmd = &cobra.Command{
	Use:   "migrate",
	Short: "migrate database schema",
	Long:  `migrate database schema`,
	Run: func(cmd *cobra.Command, args []string) {
		c, err := config.FromFile(Cfg)
		if err != nil {
			log.Fatalf("Error creating a new config %v", err)
		}

		if rollback {
			migrationDirection = 1
		}

		if c.Database != nil {
			applied, err := mysql.RunMigration(c.Database.GetDsn(), migrationDirection)

			if err != nil {
				log.Fatalf("error while running migration %v", err)
			} else {
				log.Printf("database migration ran applied %d", applied)
			}
		}
	},
}
