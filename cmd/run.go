/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"github.com/spf13/cobra"
	"test.local/prometheus-cli-exporter/internal/app"
	domain "test.local/prometheus-cli-exporter/internal/repository"
)

var c domain.Config

var runCmd = &cobra.Command{
	Use:   "run",
	Short: "run vstorage prometheus exporter",
	Long:  `That application run vstorage command with specified cluster name to provides R-Storage cluster metrics to prometheus server`,
	Run: func(cmd *cobra.Command, args []string) {

		app.Run(c)

	},
}

func init() {

	rootCmd.AddCommand(runCmd)

	runCmd.Flags().StringVarP(&c.Debug, "debug", "d", "info", "message level for traces")
	runCmd.Flags().StringVarP(&c.Port, "port", "p", "9000", "port to listen requests")
	runCmd.Flags().StringVarP(&c.IP, "ip", "i", "127.0.0.1", "listen on specified address only")
	runCmd.Flags().StringVarP(&c.Name, "cluster", "c", "SKALA-R", "name of cluster connect to")
	runCmd.MarkFlagsRequiredTogether("ip", "port", "cluster")
}
