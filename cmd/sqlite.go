package cmd

import (
	"fmt"
	"tillkuhn/goplay/sqlite"

	"github.com/spf13/cobra"
)

// sqliteCmd represents the sqlite command
var sqliteCmd = &cobra.Command{
	Use:   "sqlite",
	Short: "A brief description of your command",
	Long:  `A longer description that spans multiple lines`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("sqlite called")
		sqlite.Run()
	},
}

func init() {
	rootCmd.AddCommand(sqliteCmd)
}
