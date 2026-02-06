package cli

import (
	projectInit "github.com/f4biankoch/stackgen/internal/projectinit"
	"github.com/spf13/cobra"
)

var force bool
var template string
var defaults bool
var projectInitCmd = &cobra.Command{
	Use:   "init <project-name>",
	Short: "Initializes a new stackgen project",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		projectName := args[0]
		return projectInit.Run(projectName, force, template, defaults)
	},
}

func init() {
	rootCmd.AddCommand(projectInitCmd)
	projectInitCmd.Flags().BoolVar(&force, "force", false, "overwrite existing files if the target directory is not empty")
	projectInitCmd.Flags().StringVar(&template, "template", "basic", "specify project template to use (default: basic)")
	projectInitCmd.Flags().BoolVar(&defaults, "defaults", false, "only use template default options (recommended)")
}
