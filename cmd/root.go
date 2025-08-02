package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "mux-geo",
	Short: "Generative Engine Optimization CLI tool",
	Long: `A powerful CLI tool for Generative Engine Optimization (GEO) that provides:
- Bulk URL checking and analysis
- Local project directory scanning  
- Webpage data analysis with multiple LLM providers
- Support for Claude, GPT, and local LLMs`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Welcome to GEO Checker! Use --help to see available commands.")
	},
}

func Execute() error {
	return rootCmd.Execute()
}

func init() {
	rootCmd.AddCommand(analyzeCmd)
	rootCmd.AddCommand(bulkCmd)
	rootCmd.AddCommand(scanCmd)
}