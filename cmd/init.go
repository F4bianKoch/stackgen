package cmd

import (
	"github.com/f4biankoch/stackgen/pkg/projectInit"
	"github.com/spf13/cobra"
)

var force bool
var template string
var initCmd = &cobra.Command{
	Use:   "init <project-name>",
	Short: "Initializes a new stackgen project",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		projectName := args[0]
		return projectInit.Run(projectName, force, template)
	},
}

func init() {
	rootCmd.AddCommand(initCmd)
	initCmd.Flags().BoolVar(&force, "force", false, "overwrite existing files if the target directory is not empty")
	initCmd.Flags().StringVar(&template, "template", "basic", "specify project template to use (default: basic)")
}
