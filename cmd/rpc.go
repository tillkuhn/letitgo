package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/tillkuhn/letitgo/rpc"
	"github.com/tillkuhn/letitgo/rpcclient"
)

var rpcClientMode bool

// rpcCmd represents the rpc command
var rpcCmd = &cobra.Command{
	Use:   "rpc",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		if !rpcClientMode {
			fmt.Println("rpc called server mode")
			rpc.Serve()
		} else {
			fmt.Println("rpc called client mode")
			rpcclient.Run()
		}

	},
}

func init() {
	rootCmd.AddCommand(rpcCmd)
	rpcCmd.PersistentFlags().BoolVar(&rpcClientMode, "client", false, "client mode")
	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// rpcCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// rpcCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
