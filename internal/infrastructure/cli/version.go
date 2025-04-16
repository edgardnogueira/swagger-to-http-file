package cli

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	version     string
	commitHash  string
	buildDate   string
)

// SetVersionInfo sets the version information
func SetVersionInfo(v, commit, date string) {
	version = v
	commitHash = commit
	buildDate = date
}

// versionCmd represents the version command
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version information",
	Long:  `Print the version, commit hash, and build date of swagger-to-http-file.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("swagger-to-http-file v%s\n", version)
		fmt.Printf("Commit: %s\n", commitHash)
		fmt.Printf("Built: %s\n", buildDate)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
