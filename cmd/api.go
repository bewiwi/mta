package cmd

import (
	"github.com/bewiwi/mta/api"
	"github.com/spf13/cobra"
)

// apiCmd represents the scheduler command
var apiCmd = &cobra.Command{
	Use:   "api",
	Short: "Run api",
	Run: func(cmd *cobra.Command, args []string) {
		api.Run()
	},
}

func init() {
	RootCmd.AddCommand(apiCmd)
}
