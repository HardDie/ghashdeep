package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Version of app",
	Run: func(_ *cobra.Command, _ []string) {
		_, _ = fmt.Println(Version)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
