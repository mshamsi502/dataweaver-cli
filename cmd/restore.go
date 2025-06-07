// فایل: cmd/restore.go
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var restoreCmd = &cobra.Command{
	Use:   "restore",
	Short: "Restore data from a backup",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Please specify a subcommand, e.g., 'mongo'.")
	},
}

func init() {
	rootCmd.AddCommand(restoreCmd)
	// دستور restoreMongoCmd (که در فایل restore_mongo.go تعریف شده) را به این دستور اضافه می‌کنیم
	restoreCmd.AddCommand(restoreMongoCmd)
}
