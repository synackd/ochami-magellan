package cmd

import (
	"encoding/json"
	"fmt"
	"log"

	urlx "github.com/OpenCHAMI/magellan/internal/url"
	"github.com/OpenCHAMI/magellan/pkg/crawler"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// The `crawl` command walks a collection of Redfish endpoints to collect
// specfic inventory detail. This command only expects host names and does
// not require a scan to be performed beforehand.
var crawlCmd = &cobra.Command{
	Use:   "crawl [uri]",
	Short: "Crawl a single BMC for inventory information",
	Long: "Crawl a single BMC for inventory information. This command does NOT store information\n" +
		"about the scan into cache after completion. To do so, use the 'collect' command instead\n\n" +
		"Examples:\n" +
		"  magellan crawl https://bmc.example.com\n" +
		"  magellan crawl https://bmc.example.com -i -u username -p password",
	Args: func(cmd *cobra.Command, args []string) error {
		// Validate that the only argument is a valid URI
		var err error
		if err := cobra.ExactArgs(1)(cmd, args); err != nil {
			return err
		}
		args[0], err = urlx.Sanitize(args[0])
		if err != nil {
			return fmt.Errorf("failed to sanitize URI: %w", err)
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		systems, err := crawler.CrawlBMCForSystems(crawler.CrawlerConfig{
			URI:      args[0],
			Username: cmd.Flag("username").Value.String(),
			Password: cmd.Flag("password").Value.String(),
			Insecure: cmd.Flag("insecure").Value.String() == "true",
		})
		if err != nil {
			log.Fatalf("Error crawling BMC: %v", err)
		}
		// Marshal the inventory details to JSON
		jsonData, err := json.MarshalIndent(systems, "", "  ")
		if err != nil {
			fmt.Println("Error marshalling to JSON:", err)
			return
		}

		// Print the pretty JSON
		fmt.Println(string(jsonData))
	},
}

func init() {
	crawlCmd.Flags().StringP("username", "u", "", "Set the username for the BMC")
	crawlCmd.Flags().StringP("password", "p", "", "Set the password for the BMC")
	crawlCmd.Flags().BoolP("insecure", "i", false, "Ignore SSL errors")

	checkBindFlagError(viper.BindPFlag("crawl.username", crawlCmd.Flags().Lookup("username")))
	checkBindFlagError(viper.BindPFlag("crawl.password", crawlCmd.Flags().Lookup("password")))
	checkBindFlagError(viper.BindPFlag("crawl.insecure", crawlCmd.Flags().Lookup("insecure")))

	rootCmd.AddCommand(crawlCmd)
}
