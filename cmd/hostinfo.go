// Package cmd deals with Cobra CLI Entrypoints
package cmd

import (
	"fmt"
	"os"
	"runtime"

	"github.com/spf13/cobra"
)

// hostinfoCmd represents the hostinfo command
var hostinfoCmd = &cobra.Command{
	Use:   "hostinfo",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		name, err := os.Hostname()
		if err != nil {
			panic(err)
		}

		fmt.Printf("hostname: %s OS: %s\n", name, runtime.GOOS)
	},
}

func init() {
	rootCmd.AddCommand(hostinfoCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// hostinfoCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// hostinfoCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
