package cmd

import (
	"github.com/bewiwi/mta/consumer"
	"github.com/spf13/cobra"
)

// consumeCmd represents the consume command
var consumeCmd = &cobra.Command{
	Use:   "consume",
	Short: "Consume response",
	Run: func(cmd *cobra.Command, args []string) {
		consumer.Consume()
	},
}

func init() {
	RootCmd.AddCommand(consumeCmd)
}
