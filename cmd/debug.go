package cmd

import (
	"context"
	"fmt"
	"geo-checker/internal/webpage"
	"strings"
	"time"

	"github.com/spf13/cobra"
)

var debugCmd = &cobra.Command{
	Use:   "debug [URL]",
	Short: "Debug webpage content extraction",
	Long:  "Debug and display detailed information about webpage content extraction",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		url := args[0]
		
		scraper := webpage.New()
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()
		
		fmt.Printf("🔍 Debugging content extraction for: %s\n", url)
		fmt.Println(strings.Repeat("=", 60))
		
		pageData, err := scraper.ScrapeURL(ctx, url)
		if err != nil {
			return fmt.Errorf("failed to scrape URL: %w", err)
		}
		
		fmt.Printf("📄 Title: %s\n", pageData.Title)
		fmt.Printf("📏 Content Length: %d characters\n", len(pageData.Content))
		fmt.Printf("🏷️  Meta Tags: %d found\n", len(pageData.MetaTags))
		fmt.Printf("📋 Headings: %d found\n", len(pageData.Headings))
		fmt.Println()
		
		if len(pageData.Headings) > 0 {
			fmt.Println("📋 Headings Found:")
			for _, heading := range pageData.Headings {
				fmt.Printf("  H%d: %s\n", heading.Level, heading.Text)
			}
			fmt.Println()
		}
		
		if len(pageData.MetaTags) > 0 {
			fmt.Println("🏷️  Meta Tags Found:")
			for key, value := range pageData.MetaTags {
				if len(value) > 100 {
					value = value[:100] + "..."
				}
				fmt.Printf("  %s: %s\n", key, value)
			}
			fmt.Println()
		}
		
		if pageData.Content == "" {
			fmt.Println("❌ No content extracted!")
			fmt.Println("This could mean:")
			fmt.Println("  - The page requires JavaScript to render content")
			fmt.Println("  - The page has an unusual HTML structure")
			fmt.Println("  - The content is in a format not recognized by the scraper")
			fmt.Println("  - The page is mostly images or interactive elements")
		} else {
			fmt.Println("✅ Content Successfully Extracted:")
			fmt.Println(strings.Repeat("-", 40))
			
			// Show first 500 characters of content
			content := pageData.Content
			if len(content) > 500 {
				content = content[:500] + "...\n[Content truncated for display]"
			}
			fmt.Println(content)
		}
		
		return nil
	},
}

func init() {
	rootCmd.AddCommand(debugCmd)
}