package cmd

import (
	"fmt"
	"geo-checker/pkg/config"
	"geo-checker/pkg/formatter"
	"geo-checker/pkg/scanner"
	"geo-checker/pkg/ui"

	"github.com/spf13/cobra"
)

var scanCmd = &cobra.Command{
	Use:   "scan [directory]",
	Short: "Scan local project directory for HTML files and analyze them",
	Long:  "Recursively scan a local directory for HTML files and analyze them for GEO optimization opportunities",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		directory := args[0]
		
		provider, _ := cmd.Flags().GetString("provider")
		model, _ := cmd.Flags().GetString("model")
		output, _ := cmd.Flags().GetString("output")
		mode, _ := cmd.Flags().GetString("mode")
		extensions, _ := cmd.Flags().GetStringSlice("ext")
		
		// Show banner for text output
		if output == "text" {
			ui := ui.New()
			ui.PrintBanner()
		}
		
		cfg := &config.Config{
			LLMProvider:  provider,
			Model:        model,
			OutputFormat: output,
			Mode:         mode,
			Extensions:   extensions,
			MaxTokens:    4000,
			Temperature:  0.7,
			Timeout:      30,
		}
		
		scanner := scanner.New(cfg)
		results, err := scanner.ScanDirectory(directory)
		if err != nil {
			return fmt.Errorf("failed to scan directory: %w", err)
		}
		
		formatter := formatter.New(output)
		fmt.Print(formatter.FormatScanResults(results))
		return nil
	},
}

func init() {
	scanCmd.Flags().StringP("provider", "p", "claude", "LLM provider (claude, gpt, local)")
	scanCmd.Flags().StringP("model", "m", "claude-3-sonnet", "Model to use")
	scanCmd.Flags().StringP("output", "o", "text", "Output format (text, json, markdown)")
	scanCmd.Flags().StringP("mode", "", "local", "Analysis mode (local, llm, hybrid)")
	scanCmd.Flags().StringSliceP("ext", "e", []string{".html", ".htm"}, "File extensions to scan")
}