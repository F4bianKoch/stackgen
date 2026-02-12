package cli

import (
	projectInit "github.com/f4biankoch/stackgen/internal/projectinit"
	"github.com/spf13/cobra"
)

var projectInitCmd = &cobra.Command{
	Use:   "init <project-name>",
	Short: "Initializes a new stackgen project",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		projectName := args[0]

		force, _ := cmd.Flags().GetBool("force")
		template, _ := cmd.Flags().GetString("template")
		defaults, _ := cmd.Flags().GetBool("defaults")

		return projectInit.Run(projectName, force, template, defaults)
	},
}

func init() {
	rootCmd.AddCommand(projectInitCmd)
	projectInitCmd.Flags().Bool("force", false, "overwrite existing files if the target directory is not empty")
	projectInitCmd.Flags().String("template", "basic", "specify project template to use (default: basic)")
	projectInitCmd.Flags().Bool("defaults", false, "only use template default options (recommended)")
}
