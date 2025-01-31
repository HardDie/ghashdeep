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

// checkCmd represents the check command
var checkCmd = &cobra.Command{
	Use:   "check",
	Short: "Recursive search for checksum.* files and checksum verification",
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
			rootDir, err := os.Getwd()
			if err != nil {
				return fmt.Errorf("check os.Getwd: %w", err)
			}

			err = crawler.
				New(hash, cfg, slogger).
				Check(rootDir)
			if err != nil {
				return fmt.Errorf("check: %w", err)
			}
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
	rootCmd.AddCommand(checkCmd)
}
