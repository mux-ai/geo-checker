package cmd

import (
	"fmt"
	"geo-checker/internal/bulk"
	"geo-checker/pkg/config"
	"geo-checker/pkg/formatter"
	"geo-checker/pkg/llm"
	"geo-checker/pkg/ui"

	"github.com/spf13/cobra"
)

var bulkCmd = &cobra.Command{
	Use:   "bulk [file]",
	Short: "Analyze multiple URLs from a file",
	Long:  "Analyze multiple URLs provided in a file (one URL per line) for GEO optimization",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		file := args[0]
		
		provider, _ := cmd.Flags().GetString("provider")
		model, _ := cmd.Flags().GetString("model")
		output, _ := cmd.Flags().GetString("output")
		mode, _ := cmd.Flags().GetString("mode")
		concurrent, _ := cmd.Flags().GetInt("concurrent")
		interactive, _ := cmd.Flags().GetBool("interactive")
		
		// Interactive model selection
		if interactive {
			selectedProvider, selectedModel, err := llm.InteractiveModelSelection(provider)
			if err != nil {
				return fmt.Errorf("interactive selection failed: %w", err)
			}
			provider = selectedProvider
			model = selectedModel
		} else {
			// Validate model for provider if specified
			if model != "" && provider != "" {
				if err := llm.ValidateModelForProvider(provider, model); err != nil {
					return fmt.Errorf("model validation failed: %w", err)
				}
			}
			
			// Set recommended model if not specified
			if model == "" {
				model = llm.GetRecommendedModel(provider)
				if model == "" {
					return fmt.Errorf("no default model available for provider: %s", provider)
				}
			}
		}
		
		// Show banner for text output
		if output == "text" {
			ui := ui.New()
			ui.PrintBanner()
			
			// Display selected configuration
			fmt.Printf("Provider: %s\n", provider)
			fmt.Printf("Model: %s\n", model)
			fmt.Printf("Mode: %s\n", mode)
			fmt.Printf("Concurrent requests: %d\n\n", concurrent)
		}
		
		cfg := &config.Config{
			LLMProvider:  provider,
			Model:        model,
			OutputFormat: output,
			Mode:         mode,
			Concurrent:   concurrent,
			MaxTokens:    4000,
			Temperature:  0.7,
			Timeout:      30,
		}
		
		processor := bulk.New(cfg)
		results, err := processor.ProcessFile(file)
		if err != nil {
			return fmt.Errorf("failed to process bulk URLs: %w", err)
		}
		
		formatter := formatter.New(output)
		fmt.Print(formatter.FormatBulkResults(results))
		return nil
	},
}

func init() {
	bulkCmd.Flags().StringP("provider", "p", "claude", "LLM provider (claude, openai, local)")
	bulkCmd.Flags().StringP("model", "m", "", "Model to use (leave empty for recommended model)")
	bulkCmd.Flags().StringP("output", "o", "text", "Output format (text, json, markdown)")
	bulkCmd.Flags().StringP("mode", "", "auto", "Analysis mode (auto, local, llm, hybrid)")
	bulkCmd.Flags().IntP("concurrent", "c", 5, "Number of concurrent requests")
	bulkCmd.Flags().BoolP("interactive", "i", false, "Interactive model selection")
}