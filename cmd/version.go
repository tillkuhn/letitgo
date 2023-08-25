package cmd

import (
	"fmt"
	"runtime"

	"github.com/spf13/cobra"
)

// CommitHash Expected to be used by ldflags, e.g.  -X github.com/my-app/version.AppName=${APP}
var (
	CommitHash string
	CommitTag  string
	CommitDate string
	BuildDate  string
)

// versionCmd represents the version command
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Display Version Info",
	Long: `Display Git Information such as CommitHash
and Build date.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Version Info:")
		fmt.Printf("tag=%s hash=%s last_commit=%s build_date=%s\n", CommitTag, CommitHash, CommitDate, BuildDate)
		fmt.Printf("GOOS=%s GOOARCH=%s\n", runtime.GOOS, runtime.GOARCH)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// versionCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// versionCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
