package cmd

import (
	"github.com/bewiwi/mta/checks"
	"github.com/spf13/cobra"
)

// checkCmd represents the check command
var checkCmd = &cobra.Command{
	Use:   "check",
	Short: "Run a check",
}

var hostPing string

// pingCmd represents the ping command
var pingCmd = &cobra.Command{
	Use:   "ping",
	Short: "Run a check ping ",
	Run: func(cmd *cobra.Command, args []string) {
		checkPing := checks.NewPing(hostPing)
		answer, _ := checkPing.Run()
		answer.Print()
	},
}

func init() {
	RootCmd.AddCommand(checkCmd)

	checkCmd.AddCommand(pingCmd)
	pingCmd.Flags().StringVarP(&hostPing, "host", "H", "127.0.0.1", "Host to ping")
}
