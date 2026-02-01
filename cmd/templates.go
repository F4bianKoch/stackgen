package cmd

import (
	"github.com/f4biankoch/stackgen/pkg/templates"
	"github.com/spf13/cobra"
)

var list bool
var templatesCmd = &cobra.Command{
	Use:   "templates",
	Short: "interact with stackgens stack templates",
	RunE: func(cmd *cobra.Command, args []string) error {
		return templates.Run(list)
	},
}

func init() {
	rootCmd.AddCommand(templatesCmd)
	templatesCmd.Flags().BoolVar(&list, "list", false, "list available templates")
}
