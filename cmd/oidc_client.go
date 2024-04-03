package cmd

import (
	"fmt"

	"github.com/tillkuhn/letitgo/oidc"

	"github.com/spf13/cobra"
)

// oidcclientCmd represents the oidcclient command
var oidcclientCmd = &cobra.Command{
	Use:   "oidcclient",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. 

OIDC Client demonstrates https://github.com/zitadel/oidc.`,
	Run: func(cmd *cobra.Command, args []string) {
		// Call with go run main.go oidcclient --debug
		debugFlag, _ := cmd.Flags().GetBool("debug")
		fmt.Printf("verbose: %v", debugFlag)
		oidc.RunClient(args)
	},
	Args: cobra.MaximumNArgs(1),
}

func init() {
	rootCmd.AddCommand(oidcclientCmd)

	oidcclientCmd.Flags().BoolP("verbose", "v", false, "Verbose output for oidc client")

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// oidcclientCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// oidcclientCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
