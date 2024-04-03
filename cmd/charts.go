package cmd

import (
	"fmt"

	"github.com/tillkuhn/letitgo/charts"

	"github.com/spf13/cobra"
)

// chartsCmd represents the charts command
var chartsCmd = &cobra.Command{
	Use:   "charts",
	Short: "A brief description of your command",
	Long:  `A longer description that spans multiple lines`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("charts called")
		charts.Run() // Delegate
	},
}

func init() {
	rootCmd.AddCommand(chartsCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// chartsCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// chartsCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
