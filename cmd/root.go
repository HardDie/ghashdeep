/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/HardDie/ghashdeep/internal/crawler"
)

var Version string

var rootCmd = &cobra.Command{
	Use:   "ghashdeep",
	Short: "This utility will help you easily calculate or check previously calculated hash sums of the entire library recursively with a single command",
}

func Execute(v string) {
	Version = v
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().StringP("algorithm", "a", "md5", "The hashing algorithm you prefer to use. Possible algorithms: md5, sha1, sha224, sha256, sha384, sha512, xxhash, blake3")
}

func chooseHashAlgCmd(cmd *cobra.Command) (crawler.HashMethod, error) {
	alg, _ := cmd.Flags().GetString("algorithm")
	hash := crawler.ChooseHashAlg(alg)
	if hash == nil {
		return nil, fmt.Errorf("unknown flag --alg value %q", alg)
	}
	return hash, nil
}
