package cli

import (
	"github.com/f4biankoch/stackgen/internal/templates"
	"github.com/spf13/cobra"
)

var list bool
var templatesCmd = &cobra.Command{
	Use:   "templates",
	Short: "interact with stackgens stack templates",
	RunE: func(cmd *cobra.Command, args []string) error {
		if list {
			return templates.ListTemplates()
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(templatesCmd)
	templatesCmd.Flags().BoolVar(&list, "list", false, "list available templates")
}
