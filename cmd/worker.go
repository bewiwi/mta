package cmd

import (
	"github.com/bewiwi/mta/worker"
	"github.com/spf13/cobra"
)

// consumerCmd represents the kafka command
var consumerCmd = &cobra.Command{
	Use:   "worker",
	Short: "consume and make checks",
	Long:  "blabla",
	Run: func(cmd *cobra.Command, args []string) {
		worker.Run()
	},
}

func init() {
	RootCmd.AddCommand(consumerCmd)

}
