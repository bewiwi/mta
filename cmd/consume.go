package cmd

import (
	"github.com/bewiwi/mta/consumer/influx"
	"github.com/spf13/cobra"
	"github.com/bewiwi/mta/consumer/db"
	"github.com/bewiwi/mta/consumer/stdout"
)

// consumeCmd represents the consume command
var consumeCmd = &cobra.Command{
	Use:   "consume",
	Short: "Consume response from kafka",
}

// stdoutCmd represents the stdout command
var stdoutCmd = &cobra.Command{
	Use:   "stdout",
	Short: "Consume and display in stdout",
	Run: func(cmd *cobra.Command, args []string) {
		stdout_consumer.Run()
	},
}

// dbCmd represents the stdout command
var dbCmd = &cobra.Command{
	Use:   "db",
	Short: "Consume and insert it in database",
	Run: func(cmd *cobra.Command, args []string) {
		db_consumer.Run()
	},
}

// influxCmd represents the stdout command
var influxCmd = &cobra.Command{
	Use:   "influx",
	Short: "Consume and insert it in influxdb",
	Run: func(cmd *cobra.Command, args []string) {
		influx_consumer.Run()
	},
}

func init() {
	RootCmd.AddCommand(consumeCmd)
	consumeCmd.AddCommand(stdoutCmd)
	consumeCmd.AddCommand(dbCmd)
	consumeCmd.AddCommand(influxCmd)

}
