package cmd

import (
	"github.com/dorianim/downtimerobot/internal/announcements"
	"github.com/dorianim/downtimerobot/internal/crawler"
	"github.com/dorianim/downtimerobot/internal/frontend"
	"github.com/dorianim/downtimerobot/internal/statistics"
	"github.com/spf13/cobra"
)

// runCmd represents the run command
var frontendCmd = &cobra.Command{
	Use:   "frontend",
	Short: "Generate the frontend",
	Long:  "Generates the static frontend and statistics",
	Run: func(cmd *cobra.Command, args []string) {
		crawledServices, err := crawler.LoadServices()
		cobra.CheckErr(err)
		serviceList, serviceDetailList, err := statistics.Generate(crawledServices)
		cobra.CheckErr(err)
		announcementList, err := announcements.Generate()
		cobra.CheckErr(err)
		err = frontend.Generate(serviceList, serviceDetailList, announcementList)
		cobra.CheckErr(err)
	},
}

func init() {
	rootCmd.AddCommand(frontendCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// runCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// runCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
