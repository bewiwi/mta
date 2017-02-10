package cmd

import (
	"github.com/spf13/cobra"
	"github.com/bewiwi/mta/check"
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
		check := check.Ping{
			Host: hostPing,
		}
		response, _ := check.Run()
		response.Print()
	},
}

// httpCmd represents the ping command
var httpCmd = &cobra.Command{
	Use:   "http",
	Short: "Run a check http",
	Run: func(cmd *cobra.Command, args []string) {
		check := check.HttpCheck{
			Host: hostPing,
		}
		response, _ := check.Run()
		response.Print()
	},
}

func init() {
	RootCmd.AddCommand(checkCmd)

	checkCmd.AddCommand(pingCmd)
	pingCmd.Flags().StringVarP(&hostPing, "host", "H", "127.0.0.1", "Host to ping")

	checkCmd.AddCommand(httpCmd)
	httpCmd.Flags().StringVarP(&hostPing, "host", "H", "127.0.0.1", "HttpCheck host")
}
