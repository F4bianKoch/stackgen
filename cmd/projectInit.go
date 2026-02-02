package cmd

import (
	projectInit "github.com/f4biankoch/stackgen/pkg/projectinit"
	"github.com/spf13/cobra"
)

var force bool
var template string
var projectInitCmd = &cobra.Command{
	Use:   "init <project-name>",
	Short: "Initializes a new stackgen project",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		projectName := args[0]
		return projectInit.Run(projectName, force, template)
	},
}

func init() {
	rootCmd.AddCommand(projectInitCmd)
	projectInitCmd.Flags().BoolVar(&force, "force", false, "overwrite existing files if the target directory is not empty")
	projectInitCmd.Flags().StringVar(&template, "template", "basic", "specify project template to use (default: basic)")
}
