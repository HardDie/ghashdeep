/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package main

import (
	"runtime/debug"

	"github.com/HardDie/LibraryHashCheck/cmd"
)

var Version = "dev"

func main() {
	if info, available := debug.ReadBuildInfo(); available {
		switch info.Main.Version {
		case "", "(devel)":
			// skip
		default:
			Version = info.Main.Version
		}
	}
	cmd.Execute(Version)
}
