/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"github.com/dorianim/downtimerobot/internal/crawler"
	"github.com/dorianim/downtimerobot/internal/notifications"
	"github.com/spf13/cobra"
)

// notifyCmd represents the notify command
var notifyCmd = &cobra.Command{
	Use:   "notify",
	Short: "Sends pending notifications",
	Long: `Will send a notification to all configured targets of services
which have changed their state in the last crawl.`,
	Run: func(cmd *cobra.Command, args []string) {
		crawledServices, err := crawler.LoadServices()
		cobra.CheckErr(err)
		err = notifications.Notify(crawledServices)
		cobra.CheckErr(err)
	},
}

func init() {
	rootCmd.AddCommand(notifyCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// notifyCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// notifyCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
