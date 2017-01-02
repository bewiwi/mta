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
	"encoding/json"

	log "github.com/Sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/bewiwi/mta/kafka"
	"github.com/spf13/viper"
	"github.com/bewiwi/mta/models"
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
		consumer := kafka.GetConsumer(viper.GetStringSlice("KAFKA.TOPIC_ANSWER"))
		for msg := range consumer.Messages() {
			var checkAnswer models.CheckAnswer
			err := json.Unmarshal(msg.Value, &checkAnswer)
			if err != nil {
				log.WithError(err).Error("error unmarchal")
			}
		}
	},
}

func init() {
	RootCmd.AddCommand(consumeCmd)
	consumeCmd.AddCommand(stdoutCmd)

}
