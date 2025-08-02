package cmd

import (
	"fmt"
	"geo-checker/pkg/analyzer"
	"geo-checker/pkg/config"
	"geo-checker/pkg/formatter"
	"geo-checker/pkg/llm"
	"geo-checker/pkg/ui"

	"github.com/spf13/cobra"
)

var analyzeCmd = &cobra.Command{
	Use:   "analyze [URL]",
	Short: "Analyze a single webpage for GEO optimization",
	Long:  "Analyze a single webpage using the specified LLM provider to assess GEO optimization opportunities",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		url := args[0]
		
		provider, _ := cmd.Flags().GetString("provider")
		model, _ := cmd.Flags().GetString("model")
		output, _ := cmd.Flags().GetString("output")
		mode, _ := cmd.Flags().GetString("mode")
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
			fmt.Printf("Mode: %s\n\n", mode)
		}
		
		cfg := &config.Config{
			LLMProvider:  provider,
			Model:        model,
			OutputFormat: output,
			Mode:         mode,
			MaxTokens:    4000,
			Temperature:  0.7,
			Timeout:      30,
		}
		
		analyzer := analyzer.New(cfg)
		result, err := analyzer.AnalyzeURL(url)
		if err != nil {
			return fmt.Errorf("failed to analyze URL: %w", err)
		}
		
		formatter := formatter.New(output)
		fmt.Print(formatter.FormatAnalysisResult(result))
		return nil
	},
}

func init() {
	analyzeCmd.Flags().StringP("provider", "p", "claude", "LLM provider (claude, openai, local)")
	analyzeCmd.Flags().StringP("model", "m", "", "Model to use (leave empty for recommended model)")
	analyzeCmd.Flags().StringP("output", "o", "text", "Output format (text, json, markdown)")
	analyzeCmd.Flags().StringP("mode", "", "auto", "Analysis mode (auto, local, llm, hybrid)")
	analyzeCmd.Flags().BoolP("interactive", "i", false, "Interactive model selection")
}