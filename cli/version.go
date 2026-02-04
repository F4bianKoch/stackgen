/*
Copyright © 2026 NAME HERE <EMAIL ADDRESS>
*/
package cli

import (
	"fmt"

	"github.com/spf13/cobra"
)

// versionCmd represents the version command
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "stackgen version information",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("stackgen %s\n", Version)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
