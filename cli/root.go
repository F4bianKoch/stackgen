package cli

import (
	"os"

	"github.com/spf13/cobra"
)

var Version = "0.1.2"
var rootCmd = &cobra.Command{
	Use:   "stackgen",
	Short: "Stackgen is a CLI tool to generate infrastructure stacks",
	Long:  `Stackgen is a command-line interface (CLI) tool to generate infrastructure stacks quickly and efficiently.`,
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
}
