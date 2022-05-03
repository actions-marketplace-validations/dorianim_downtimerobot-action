package cmd

import (
	"github.com/dorianim/downtimerobot/internal/crawler"
	"github.com/spf13/cobra"
)

// crawlCmd represents the crawl command
var crawlCmd = &cobra.Command{
	Use:   "crawl",
	Short: "Crawl all configured services",
	Long:  `Crawll all services configured in downtimerobot.yml and generate their statistics`,
	Run: func(cmd *cobra.Command, args []string) {
		_, err := crawler.CrawlServices()
		cobra.CheckErr(err)
	},
}

func init() {
	rootCmd.AddCommand(crawlCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// crawlCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// crawlCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
