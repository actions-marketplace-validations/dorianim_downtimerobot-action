/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"github.com/dorianim/downtimerobot/internal/announcements"
	"github.com/dorianim/downtimerobot/internal/crawler"
	"github.com/dorianim/downtimerobot/internal/frontend"
	"github.com/dorianim/downtimerobot/internal/notifications"
	"github.com/dorianim/downtimerobot/internal/statistics"
	"github.com/spf13/cobra"
)

// runCmd represents the run command
var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Crawl all services and generate the frontend and send pending notifications",
	Long: `Checks all services and writes the result to the data files.
In addition to that, it generates the static frontend and statistics.
It also sends pending notifications.
It is equal to running crawl, frontend and notify.`,
	Run: func(cmd *cobra.Command, args []string) {
		crawledServices, err := crawler.CrawlServices()
		cobra.CheckErr(err)
		serviceList, serviceDetailList, err := statistics.Generate(crawledServices)
		cobra.CheckErr(err)
		announcementList, err := announcements.Generate()
		cobra.CheckErr(err)
		err = frontend.Generate(serviceList, serviceDetailList, announcementList)
		cobra.CheckErr(err)
		err = notifications.Notify(crawledServices)
		cobra.CheckErr(err)
	},
}

func init() {
	rootCmd.AddCommand(runCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// runCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// runCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
