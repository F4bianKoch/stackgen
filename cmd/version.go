/*
Copyright © 2026 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"github.com/f4biankoch/stackgen/pkg/version"
	"github.com/spf13/cobra"
)

// versionCmd represents the version command
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "stackgen version information",
	Run: func(cmd *cobra.Command, args []string) {
		version.PrintVersion()
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
