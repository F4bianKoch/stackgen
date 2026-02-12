package cli

import (
	projectReinit "github.com/f4biankoch/stackgen/internal/projectreinit"
	"github.com/spf13/cobra"
)

var projectReinitCmd = &cobra.Command{
	Use:   "reinit",
	Short: "Reinitializes the stackgen project it's executed in",
	RunE: func(cmd *cobra.Command, args []string) error {
		defaults, _ := cmd.Flags().GetBool("defaults")

		return projectReinit.Run(defaults)
	},
}

func init() {
	rootCmd.AddCommand(projectReinitCmd)
	projectReinitCmd.Flags().Bool("defaults", false, "only use template default options (recommended)")
}
