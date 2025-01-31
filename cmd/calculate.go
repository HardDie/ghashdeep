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

	"github.com/HardDie/ghashdeep/internal/crawler"
	"github.com/HardDie/ghashdeep/internal/logger"
)

// calculateCmd represents the calculate command
var calculateCmd = &cobra.Command{
	Use:   "calculate",
	Short: "Recursive hash calculation for all files",
	Run: func(cmd *cobra.Command, _ []string) {
		slogger := logger.New()

		hash, err := chooseHashAlgCmd(cmd)
		if err != nil {
			log.Fatal(err)
		}
		slogger.Info("Hash algorithm: " + hash.Name())

		cfg := getConfig(cmd)
		if cfg.Verbose {
			slogger.Debug(fmt.Sprintf("Config: %+v", cfg))
		}

		// listen app termination signals.
		signalChan := make(chan os.Signal, 1)
		signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)

		// init run group.
		var g run.Group

		// signal handler.
		g.Add(func() error {
			sig := <-signalChan
			slogger.Info(
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
				return fmt.Errorf("operation has been declined")
			}

			rootDir, err := os.Getwd()
			if err != nil {
				return fmt.Errorf("check os.Getwd: %w", err)
			}

			err = crawler.
				New(hash, cfg, slogger).
				Calculate(rootDir)
			if err != nil {
				return fmt.Errorf("calculate: %w", err)
			}

			fmt.Println("All checksums has been calculated!")
			return nil
		}, func(_ error) {
			os.Exit(0)
		})

		if err := g.Run(); err != nil {
			slogger.Info(
				"application stopped due to",
				slog.String("reason", err.Error()),
			)
		} else {
			slogger.Info(
				"application stopped",
			)
		}
	},
}

func init() {
	rootCmd.AddCommand(calculateCmd)
}
