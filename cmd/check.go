// Copyright © 2016 NAME HERE <EMAIL ADDRESS>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"github.com/spf13/cobra"
	"github.com/bewiwi/mta/checks"
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
