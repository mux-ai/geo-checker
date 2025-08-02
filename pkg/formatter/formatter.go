package formatter

import (
	"encoding/json"
	"fmt"
	"geo-checker/internal/bulk"
	"geo-checker/pkg/analyzer"
	"geo-checker/pkg/scanner"
	"geo-checker/pkg/ui"
	"strings"
	"time"
)

type Formatter struct {
	format string
	ui     *ui.UI
}

func New(format string) *Formatter {
	return &Formatter{
		format: format,
		ui:     ui.New(),
	}
}

func (f *Formatter) FormatAnalysisResult(result *analyzer.Result) string {
	switch f.format {
	case "json":
		return f.formatJSON(result)
	case "markdown":
		return f.formatMarkdown(result)
	default:
		return f.formatText(result)
	}
}

func (f *Formatter) FormatBulkResults(results []*bulk.BulkResult) string {
	switch f.format {
	case "json":
		return f.formatBulkJSON(results)
	case "markdown":
		return f.formatBulkMarkdown(results)
	default:
		return f.formatBulkText(results)
	}
}

func (f *Formatter) FormatScanResults(results []*scanner.ScanResult) string {
	switch f.format {
	case "json":
		return f.formatScanJSON(results)
	case "markdown":
		return f.formatScanMarkdown(results)
	default:
		return f.formatScanText(results)
	}
}

func (f *Formatter) formatText(result *analyzer.Result) string {
	var sb strings.Builder
	
	// Set UI color mode
	f.ui.NoColor = false
	
	// Header
	f.ui.PrintHeader("GEO ANALYSIS REPORT")
	fmt.Println()
	
	// Basic info section
	f.ui.PrintSection("ANALYSIS DETAILS")
	if result.URL != "" {
		f.ui.PrintKeyValue("URL", result.URL)
	}
	if result.Title != "" {
		f.ui.PrintKeyValue("Title", result.Title)
	}
	f.ui.PrintKeyValue("Mode", strings.ToTitle(result.Mode))
	f.ui.PrintKeyValue("Analyzed", result.ProcessedAt.Format("2006-01-02 15:04:05"))
	if result.TokensUsed > 0 {
		f.ui.PrintKeyValue("Tokens", fmt.Sprintf("%d", result.TokensUsed))
	}
	
	// Overall score
	fmt.Println()
	f.ui.PrintSection("OVERALL SCORE")
	f.ui.PrintScore("GEO Score", result.Score, 100)
	
	// Add scoring method information
	if scoringMethod, exists := result.Metadata["scoring_method"]; exists {
		switch scoringMethod {
		case "hybrid_averaged":
			if localScore, hasLocal := result.Metadata["local_score"]; hasLocal {
				if llmScore, hasLLM := result.Metadata["llm_score"]; hasLLM {
					fmt.Printf("    📊 Hybrid Score (Local: %v, LLM: %v, Averaged)\n", localScore, llmScore)
				}
			}
		case "local_only":
			fmt.Printf("    📏 Local Rule-Based Scoring\n")
		case "local_only_fallback":
			fmt.Printf("    📏 Local Scoring (LLM unavailable)\n")
		case "llm_only":
			fmt.Printf("    🤖 LLM-Based Scoring\n")
		}
	}
	
	// Detailed breakdown
	if result.LocalScore != nil {
		fmt.Println()
		f.ui.PrintSection("DETAILED BREAKDOWN")
		f.ui.PrintScore("Content Structure", 
			result.LocalScore.Breakdown.ContentStructure.Score, 100)
		f.ui.PrintScore("Semantic Clarity", 
			result.LocalScore.Breakdown.SemanticClarity.Score, 100)
		f.ui.PrintScore("Context Richness", 
			result.LocalScore.Breakdown.ContextRichness.Score, 100)
		f.ui.PrintScore("Authority Signals", 
			result.LocalScore.Breakdown.AuthoritySignals.Score, 100)
		f.ui.PrintScore("Accessibility", 
			result.LocalScore.Breakdown.Accessibility.Score, 100)
		
		// Strengths
		if len(result.LocalScore.Strengths) > 0 {
			fmt.Println()
			f.ui.PrintSubsection("Strengths")
			for _, strength := range result.LocalScore.Strengths {
				f.ui.PrintListItem(strength, true)
			}
		}
		
		// Recommendations
		if len(result.Suggestions) > 0 {
			fmt.Println()
			f.ui.PrintSubsection("Recommendations")
			for i, suggestion := range result.Suggestions {
				fmt.Printf("    %2d. %s\n", i+1, suggestion)
			}
		}
	}
	
	// LLM Analysis and recommendations
	if result.Analysis != "" {
		// Check if this contains LLM insights or just local analysis
		if strings.Contains(result.Analysis, "Enhanced Analysis Recommendation") {
			// Split analysis into local part and recommendation part
			parts := strings.Split(result.Analysis, "## 🤖 Enhanced Analysis Recommendation")
			if len(parts) > 1 {
				fmt.Println()
				f.ui.PrintMarkdownContent("## 🤖 Enhanced Analysis Recommendation" + parts[1])
			}
		} else if result.Mode != "local" {
			// This is LLM analysis content - format it beautifully
			fmt.Println()
			f.ui.PrintSection("AI INSIGHTS")
			fmt.Println()
			
			// Format the LLM response as markdown
			f.ui.PrintMarkdownContent(result.Analysis)
		}
	}
	
	fmt.Println()
	
	return sb.String()
}

func (f *Formatter) formatJSON(result *analyzer.Result) string {
	data, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		return fmt.Sprintf("Error formatting JSON: %v", err)
	}
	return string(data)
}

func (f *Formatter) formatMarkdown(result *analyzer.Result) string {
	var sb strings.Builder
	
	sb.WriteString("# GEO Analysis Report\n\n")
	if result.URL != "" {
		sb.WriteString(fmt.Sprintf("**URL:** %s\n", result.URL))
	}
	if result.Title != "" {
		sb.WriteString(fmt.Sprintf("**Title:** %s\n", result.Title))
	}
	sb.WriteString(fmt.Sprintf("**Analyzed:** %s\n", result.ProcessedAt.Format(time.RFC3339)))
	if result.TokensUsed > 0 {
		sb.WriteString(fmt.Sprintf("**Tokens Used:** %d\n", result.TokensUsed))
	}
	sb.WriteString("\n## Analysis\n\n")
	sb.WriteString(result.Analysis)
	sb.WriteString("\n")
	
	return sb.String()
}

func (f *Formatter) formatBulkText(results []*bulk.BulkResult) string {
	var sb strings.Builder
	
	f.ui.PrintHeader("GEO BULK ANALYSIS REPORT")
	
	successCount := 0
	errorCount := 0
	totalScore := 0
	
	for i, result := range results {
		f.ui.PrintSection(fmt.Sprintf("RESULT %d", i+1))
		f.ui.PrintKeyValue("URL", result.URL)
		
		if result.Error != "" {
			fmt.Println()
			f.ui.PrintError(fmt.Sprintf("Analysis failed: %s", result.Error))
			errorCount++
		} else if result.Result != nil {
			f.ui.PrintKeyValue("Title", result.Result.Title)
			if result.Result.TokensUsed > 0 {
				f.ui.PrintKeyValue("Tokens", fmt.Sprintf("%d", result.Result.TokensUsed))
			}
			fmt.Println()
			f.ui.PrintScore("GEO Score", result.Result.Score, 100)
			
			// Show all recommendations
			if len(result.Result.Suggestions) > 0 {
				fmt.Println()
				f.ui.PrintSubsection("Recommendations")
				for _, suggestion := range result.Result.Suggestions {
					f.ui.PrintListItem(suggestion, false)
				}
			}
			
			successCount++
			totalScore += result.Result.Score
		}
		fmt.Println()
	}
	
	// Summary
	f.ui.PrintSection("SUMMARY")
	f.ui.PrintKeyValue("Total URLs", fmt.Sprintf("%d", len(results)))
	f.ui.PrintKeyValue("Successful", fmt.Sprintf("%d", successCount))
	f.ui.PrintKeyValue("Errors", fmt.Sprintf("%d", errorCount))
	
	if successCount > 0 {
		avgScore := totalScore / successCount
		f.ui.PrintKeyValue("Average", fmt.Sprintf("%d/100", avgScore))
		fmt.Println()
		
		if avgScore >= 80 {
			f.ui.PrintSuccess("Excellent overall GEO performance! 🎉")
		} else if avgScore >= 60 {
			f.ui.PrintWarning("Good GEO performance with room for improvement")
		} else {
			f.ui.PrintError("GEO performance needs significant improvement")
		}
	}
	
	return sb.String()
}

func (f *Formatter) formatBulkJSON(results []*bulk.BulkResult) string {
	data, err := json.MarshalIndent(results, "", "  ")
	if err != nil {
		return fmt.Sprintf("Error formatting JSON: %v", err)
	}
	return string(data)
}

func (f *Formatter) formatBulkMarkdown(results []*bulk.BulkResult) string {
	var sb strings.Builder
	
	sb.WriteString("# GEO Bulk Analysis Report\n\n")
	
	successCount := 0
	errorCount := 0
	
	for i, result := range results {
		sb.WriteString(fmt.Sprintf("## Result %d\n\n", i+1))
		sb.WriteString(fmt.Sprintf("**URL:** %s\n\n", result.URL))
		
		if result.Error != "" {
			sb.WriteString(fmt.Sprintf("**ERROR:** %s\n\n", result.Error))
			errorCount++
		} else if result.Result != nil {
			sb.WriteString(fmt.Sprintf("**Title:** %s\n", result.Result.Title))
			sb.WriteString(fmt.Sprintf("**Tokens Used:** %d\n\n", result.Result.TokensUsed))
			sb.WriteString("### Analysis\n\n")
			sb.WriteString(result.Result.Analysis)
			sb.WriteString("\n\n")
			successCount++
		}
	}
	
	sb.WriteString("## Summary\n\n")
	sb.WriteString(fmt.Sprintf("- **Total URLs:** %d\n", len(results)))
	sb.WriteString(fmt.Sprintf("- **Successful:** %d\n", successCount))
	sb.WriteString(fmt.Sprintf("- **Errors:** %d\n", errorCount))
	
	return sb.String()
}

func (f *Formatter) formatScanText(results []*scanner.ScanResult) string {
	var sb strings.Builder
	
	f.ui.PrintHeader("GEO DIRECTORY SCAN REPORT")
	fmt.Println()
	
	successCount := 0
	errorCount := 0
	totalScore := 0
	
	for i, result := range results {
		f.ui.PrintSection(fmt.Sprintf("FILE %d", i+1))
		f.ui.PrintKeyValue("Path", result.FilePath)
		
		if result.Error != "" {
			fmt.Println()
			f.ui.PrintError(fmt.Sprintf("Analysis failed: %s", result.Error))
			errorCount++
		} else if result.Result != nil {
			f.ui.PrintKeyValue("Title", result.Result.Title)
			if result.Result.TokensUsed > 0 {
				f.ui.PrintKeyValue("Tokens", fmt.Sprintf("%d", result.Result.TokensUsed))
			}
			fmt.Println()
			f.ui.PrintScore("GEO Score", result.Result.Score, 100)
			
			// Show all recommendations if available
			if len(result.Result.Suggestions) > 0 {
				fmt.Println()
				f.ui.PrintSubsection("Recommendations")
				for _, suggestion := range result.Result.Suggestions {
					f.ui.PrintListItem(suggestion, false)
				}
			}
			
			successCount++
			totalScore += result.Result.Score
		}
		fmt.Println()
	}
	
	// Summary
	f.ui.PrintSection("SUMMARY")
	f.ui.PrintKeyValue("Total Files", fmt.Sprintf("%d", len(results)))
	f.ui.PrintKeyValue("Successful", fmt.Sprintf("%d", successCount))
	f.ui.PrintKeyValue("Errors", fmt.Sprintf("%d", errorCount))
	
	if successCount > 0 {
		avgScore := totalScore / successCount
		f.ui.PrintKeyValue("Average", fmt.Sprintf("%d/100", avgScore))
		fmt.Println()
		
		if avgScore >= 80 {
			f.ui.PrintSuccess("Excellent directory GEO performance! 🎉")
		} else if avgScore >= 60 {
			f.ui.PrintWarning("Good directory GEO performance with room for improvement")
		} else {
			f.ui.PrintError("Directory GEO performance needs significant improvement")
		}
	}
	
	return sb.String()
}

func (f *Formatter) formatScanJSON(results []*scanner.ScanResult) string {
	data, err := json.MarshalIndent(results, "", "  ")
	if err != nil {
		return fmt.Sprintf("Error formatting JSON: %v", err)
	}
	return string(data)
}

func (f *Formatter) formatScanMarkdown(results []*scanner.ScanResult) string {
	var sb strings.Builder
	
	sb.WriteString("# GEO Directory Scan Report\n\n")
	
	successCount := 0
	errorCount := 0
	
	for i, result := range results {
		sb.WriteString(fmt.Sprintf("## File %d\n\n", i+1))
		sb.WriteString(fmt.Sprintf("**Path:** `%s`\n\n", result.FilePath))
		
		if result.Error != "" {
			sb.WriteString(fmt.Sprintf("**ERROR:** %s\n\n", result.Error))
			errorCount++
		} else if result.Result != nil {
			sb.WriteString(fmt.Sprintf("**Title:** %s\n", result.Result.Title))
			sb.WriteString(fmt.Sprintf("**Tokens Used:** %d\n\n", result.Result.TokensUsed))
			sb.WriteString("### Analysis\n\n")
			sb.WriteString(result.Result.Analysis)
			sb.WriteString("\n\n")
			successCount++
		}
	}
	
	sb.WriteString("## Summary\n\n")
	sb.WriteString(fmt.Sprintf("- **Total Files:** %d\n", len(results)))
	sb.WriteString(fmt.Sprintf("- **Successful:** %d\n", successCount))
	sb.WriteString(fmt.Sprintf("- **Errors:** %d\n", errorCount))
	
	return sb.String()
}