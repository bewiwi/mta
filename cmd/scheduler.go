package cmd

import (
	"github.com/bewiwi/mta/scheduler"
	"github.com/spf13/cobra"
)

// schedulerCmd represents the scheduler command
var schedulerCmd = &cobra.Command{
	Use:   "scheduler",
	Short: "scheduler check",
	Run: func(cmd *cobra.Command, args []string) {
		scheduler.Run()
	},
}

func init() {
	RootCmd.AddCommand(schedulerCmd)
}
