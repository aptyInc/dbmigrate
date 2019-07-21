
package cmd

import (
	"github.com/spf13/cobra"
)

// downCmd represents the down command
var downCmd = &cobra.Command{
	Use:   "down",
	Short: "Down grades all the scripts that have run in the last migration",
	Long: `This is used to rollback the server. All the migrations in last up or down graded in this`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return migration.MigrateDown()
	},
}

func init() {
	rootCmd.AddCommand(downCmd)
}
