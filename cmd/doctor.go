package cmd

import (
	"github.com/f4biankoch/stackgen/pkg/doctor"
	"github.com/spf13/cobra"
)

var doctorCmd = &cobra.Command{
	Use:   "doctor",
	Short: "Check whether this system can run stackgen (Docker, Compose, permissions, filesystem)",
	RunE: func(cmd *cobra.Command, args []string) error {
		return doctor.Run()
	},
}

func init() {
	rootCmd.AddCommand(doctorCmd)
}
