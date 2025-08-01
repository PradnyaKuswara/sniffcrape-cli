package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "sniffcrape-cli",
	Short: "Sniffcrape CLI is a simple web scraper tool",
	Long:  `Sniffcrape CLI is a command-line utility to scrape websites using Go and Colly.`,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
