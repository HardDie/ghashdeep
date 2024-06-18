package cmd

import (
	"fmt"
	"log"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/oklog/run"
	"github.com/spf13/cobra"

	"github.com/HardDie/LibraryHashCheck/internal/crawler"
	"github.com/HardDie/LibraryHashCheck/internal/logger"
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
		// listen app termination signals.
		signalChan := make(chan os.Signal, 1)
		signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)

		// init run group.
		var g run.Group

		// signal handler.
		g.Add(func() error {
			sig := <-signalChan
			logger.Info(
				"signal",
				slog.String(logger.LogValueSignal, sig.String()),
			)
			return fmt.Errorf("interrupted by signal: %s", sig.String())
		}, func(error) {
			signal.Stop(signalChan)
		})

		g.Add(func() error {
			fmt.Println("All existing checksums will be replaced with new checksums.")
			fmt.Println("Are you sure? [YES/NO]:")
			fmt.Println()
			var answer string
			_, err := fmt.Scanf("%s", &answer)
			if err != nil {
				log.Fatalf("fmt.Scanf: %v", err)
			}
			if answer != "YES" {
				return fmt.Errorf("Operation has been declined!")
			}

			rootDir, err := os.Getwd()
			if err != nil {
				return fmt.Errorf("check os.Getwd: %w", err)
			}

			var hash crawler.HashMethod
			hash = validators.NewMd5()

			err = crawler.
				New(hash).
				Calculate(rootDir)
			if err != nil {
				return fmt.Errorf("calculate: %w", err)
			}

			fmt.Println("All checksums has been calculated!")
			return nil
		}, func(err error) {
			os.Exit(0)
		})

		if err := g.Run(); err != nil {
			logger.Info(
				"application stopped due to",
				slog.String("reason", err.Error()),
			)
		} else {
			logger.Info(
				"application stopped",
			)
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
