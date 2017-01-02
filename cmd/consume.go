package cmd

import (
	"github.com/bewiwi/mta/consumer"
	"github.com/spf13/cobra"
)

// consumeCmd represents the consume command
var consumeCmd = &cobra.Command{
	Use:   "consume",
	Short: "Consume answer from kafka",
}

// stdoutCmd represents the stdout command
var stdoutCmd = &cobra.Command{
	Use:   "stdout",
	Short: "Consume and display in stdout",
	Run: func(cmd *cobra.Command, args []string) {
		consumer.RunStdout()
	},
}

func init() {
	RootCmd.AddCommand(consumeCmd)
	consumeCmd.AddCommand(stdoutCmd)

}
