package cmd

import (
	"fmt"
	"geo-checker/pkg/llm"
	"strings"

	"github.com/spf13/cobra"
)

var modelsCmd = &cobra.Command{
	Use:   "models [provider]",
	Short: "List available models for LLM providers",
	Long:  "List all available models for the specified provider, or all providers if none specified",
	Args:  cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		models := llm.GetAvailableModels()
		
		if len(args) == 1 {
			// Show models for specific provider
			provider := args[0]
			providerModels, exists := models[provider]
			if !exists {
				return fmt.Errorf("unknown provider: %s. Available providers: claude, openai, local", provider)
			}
			
			fmt.Printf("📋 %s Models\n", strings.ToUpper(provider))
			fmt.Println(strings.Repeat("=", 50))
			
			for _, model := range providerModels {
				indicator := " "
				if model.Recommended {
					indicator = "⭐"
				}
				fmt.Printf("%s %s\n", indicator, model.Name)
				fmt.Printf("   %s\n", model.Description)
				fmt.Printf("   Max tokens: %d\n", model.MaxTokens)
				fmt.Println()
			}
		} else {
			// Show all models
			fmt.Println("🤖 Available LLM Models")
			fmt.Println(strings.Repeat("=", 50))
			
			for provider, providerModels := range models {
				fmt.Printf("\n📋 %s:\n", strings.ToUpper(provider))
				for _, model := range providerModels {
					indicator := " "
					if model.Recommended {
						indicator = "⭐"
					}
					fmt.Printf("%s %s - %s\n", indicator, model.Name, model.Description)
				}
			}
			
			fmt.Println("\n⭐ = Recommended model")
			fmt.Println("\nUsage:")
			fmt.Println("  mux-geo analyze <url> --provider <provider> --model <model>")
			fmt.Println("  mux-geo analyze <url> --interactive")
		}
		
		return nil
	},
}

func init() {
	rootCmd.AddCommand(modelsCmd)
}