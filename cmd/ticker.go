package cmd

import (
	"fmt"

	"github.com/tillkuhn/letitgo/ticker"

	"github.com/spf13/cobra"
)

// tickerCmd represents the ticker command
var tickerCmd = &cobra.Command{
	Use:   "ticker",
	Short: "Run Task periodically",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("ticker command called")
		ticker.RunTickerWithChannel()
	},
}

func init() {
	rootCmd.AddCommand(tickerCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// tickerCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// tickerCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
