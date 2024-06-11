package cmd

import (
	"fmt"
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
		fmt.Println("All existing checksums will be replaced with new checksums.")
		fmt.Println("Are you sure? [YES/NO]:")
		fmt.Println()
		var answer string
		_, err := fmt.Scanf("%s", &answer)
		if err != nil {
			log.Fatalf("fmt.Scanf: %v", err)
		}
		if answer != "YES" {
			fmt.Println("Operation has been declined!")
			return
		}

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

		fmt.Println("All checksums has been calculated!")
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
