package cmd

import (
	"encoding/json"
	"fmt"

	"github.com/PradnyaKuswara/sniffcrape-cli/internal/services"
	"github.com/PradnyaKuswara/sniffcrape/pkg/logger"
	"github.com/spf13/cobra"
)

var url string

var scrapeCmd = &cobra.Command{
	Use:   "scrape",
	Short: "Scrape a website URL",
	Run: func(cmd *cobra.Command, args []string) {
		logger.InitLogger()
		logger.Log.Info("Starting Sniffcrape...")

		if url == "" {
			logger.Log.Error("No URL provided. Use --url flag.")
			return
		}

		scrapperService := services.NewScrapperService()
		result, err := scrapperService.ScrapeColly(url)
		if err != nil {
			logger.Log.Error("Error scraping URL: ", err)
			return
		}
		data, err := json.MarshalIndent(result, "", "  ")
		if err != nil {
			logger.Log.Error("Error formatting result:", err)
			return
		}
		fmt.Println("Scraped content:\n", string(data))
		logger.Log.Info("Scraping completed successfully")
	},
}

func init() {
	scrapeCmd.Flags().StringVarP(&url, "url", "u", "", "URL to scrape (required)")
	rootCmd.AddCommand(scrapeCmd)
}
