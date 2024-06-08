package cmd

import (
	"log"
	"os"

	"github.com/spf13/cobra"

	"github.com/HardDie/LibraryHashCheck/internal/crawler"
	"github.com/HardDie/LibraryHashCheck/internal/validators"
)

// calculateCmd represents the calculate command
var calculateCmd = &cobra.Command{
	Use:   "calculate",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		rootDir, err := os.Getwd()
		if err != nil {
			log.Fatalf("check os.Getwd: %v", err)
		}

		var hash crawler.HashMethod
		hash = validators.NewMd5()

		err = crawler.
			New(hash).
			Calculate(rootDir)
		if err != nil {
			log.Fatalf("calculate: %v", err.Error())
		}
	},
}

func init() {
	rootCmd.AddCommand(calculateCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// calculateCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// calculateCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
