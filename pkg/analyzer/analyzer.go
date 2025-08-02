package analyzer

import (
	"context"
	"fmt"
	"geo-checker/internal/webpage"
	"geo-checker/pkg/config"
	"geo-checker/pkg/llm"
	"geo-checker/pkg/scorer"
	"geo-checker/pkg/ui"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"
)

type Analyzer struct {
	config        *config.Config
	provider      llm.Provider
	scraper       *webpage.Scraper
	localScorer   *scorer.LocalScorer
	ui            *ui.UI
	initError     error // Store initialization errors for LLM mode
	originalMode  string // Store original mode before auto-detection
}

type Result struct {
	URL           string            `json:"url"`
	Title         string            `json:"title"`
	Analysis      string            `json:"analysis,omitempty"`
	LocalScore    *scorer.GEOScore  `json:"local_score,omitempty"`
	Score         int               `json:"score"`
	Suggestions   []string          `json:"suggestions"`
	Metadata      map[string]any    `json:"metadata"`
	ProcessedAt   time.Time         `json:"processed_at"`
	TokensUsed    int               `json:"tokens_used"`
	Mode          string            `json:"mode"` // "local", "llm", or "hybrid"
}

// loadSystemPrompt loads the system prompt from SYSTEM_PROMPT.md file
func loadSystemPrompt() (string, error) {
	// Get the executable directory to find SYSTEM_PROMPT.md relative to it
	execDir, err := os.Executable()
	if err != nil {
		return "", fmt.Errorf("failed to get executable path: %w", err)
	}
	
	// Look for SYSTEM_PROMPT.md in the project root (relative to executable)
	projectRoot := filepath.Dir(filepath.Dir(filepath.Dir(execDir)))
	promptPath := filepath.Join(projectRoot, "SYSTEM_PROMPT.md")
	
	// If not found there, try current working directory
	if _, err := os.Stat(promptPath); os.IsNotExist(err) {
		wd, _ := os.Getwd()
		promptPath = filepath.Join(wd, "SYSTEM_PROMPT.md")
	}
	
	// Read the file
	content, err := ioutil.ReadFile(promptPath)
	if err != nil {
		return "", fmt.Errorf("failed to read SYSTEM_PROMPT.md from %s: %w", promptPath, err)
	}
	
	return string(content), nil
}

// getGeoPrompt returns a condensed analysis prompt for LLM calls
func getGeoPrompt() string {
	// Use a condensed version optimized for LLM analysis to avoid timeouts
	return `Analyze this webpage content for Generative Engine Optimization (GEO).

You are an expert SEO/GEO auditor evaluating content for AI systems like ChatGPT, Claude, Gemini, and Perplexity. Assess how well this content would be cited, quoted, or referenced in AI answers.

## Core GEO Factors (Score each 0-100):
1. **Content Structure** (25%): Heading hierarchy, logical flow, organization
2. **Semantic Clarity** (25%): Clear language, defined concepts, unambiguous terms
3. **Context Richness** (20%): Sufficient detail, examples, background information
4. **Authority Signals** (15%): Citations, expertise indicators, credibility markers
5. **Accessibility** (15%): Meta tags, machine-readable structure, AI-friendly formatting

## Required Output Format:

**Overall Score: [X]/100**

### Analysis Summary
Brief assessment of GEO readiness for AI citation and reference.

### Factor Scores
| Factor | Score | Key Finding |
|--------|-------|-------------|
| Content Structure | X/100 | Brief note |
| Semantic Clarity | X/100 | Brief note |
| Context Richness | X/100 | Brief note |
| Authority Signals | X/100 | Brief note |
| Accessibility | X/100 | Brief note |

### Key Recommendations
- **High Impact**: Most important improvement
- **Quick Win**: Easy implementation with good ROI
- **Long-term**: Strategic optimization for AI visibility

### AI Optimization Priority
Focus area for maximizing citation potential in AI responses.

CRITICAL: Start response with "Overall Score: [number]/100" for score extraction.`
}

func New(cfg *config.Config) *Analyzer {
	analyzer := &Analyzer{
		config:      cfg,
		scraper:     webpage.New(),
		localScorer: scorer.NewLocalScorer(),
		ui:          ui.New(),
	}

	// Intelligent mode selection based on available API keys
	originalMode := cfg.Mode
	if cfg.Mode == "auto" || cfg.Mode == "" {
		cfg.Mode = determineOptimalMode(cfg.LLMProvider)
		
		// Auto-select provider if the specified one doesn't have a valid API key
		if cfg.Mode == "hybrid" && !hasValidAPIKey(cfg.LLMProvider) {
			if hasValidAPIKey("openai") {
				cfg.LLMProvider = "openai"
			} else if hasValidAPIKey("claude") {
				cfg.LLMProvider = "claude"
			}
		}
	}
	analyzer.originalMode = originalMode
	
	// Only initialize LLM provider if not in local-only mode
	if cfg.Mode != "local" {
		providerConfig := &llm.ProviderConfig{
			APIKey:      getAPIKey(cfg.LLMProvider),
			Model:       cfg.Model,
			MaxTokens:   cfg.MaxTokens,
			Temperature: cfg.Temperature,
			BaseURL:     cfg.LocalLLMURL,
		}
		
		provider, err := llm.NewProvider(cfg.LLMProvider, providerConfig)
		if err != nil {
			if cfg.Mode == "llm" {
				// Store the error to be returned when trying to analyze
				analyzer.initError = fmt.Errorf("failed to initialize LLM provider: %w", err)
			} else {
				// In hybrid mode, continue without LLM if initialization fails
				fmt.Printf("Warning: LLM provider initialization failed, falling back to local-only mode: %v\n", err)
				cfg.Mode = "local"
			}
		} else {
			analyzer.provider = provider
		}
	}
	
	return analyzer
}

func (a *Analyzer) AnalyzeURL(url string) (*Result, error) {
	// Don't show animations for JSON output
	showAnimations := a.config.OutputFormat != "json"
	
	if showAnimations {
		a.ui.StartSpinner("Fetching webpage content...")
	}
	
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(a.config.Timeout)*time.Second)
	defer cancel()
	
	pageData, err := a.scraper.ScrapeURL(ctx, url)
	if err != nil {
		if showAnimations {
			a.ui.StopSpinner()
		}
		return nil, fmt.Errorf("failed to scrape URL: %w", err)
	}
	
	if showAnimations {
		a.ui.UpdateSpinner("Analyzing content...")
	}
	
	// Debug: Check if content was extracted successfully
	if strings.TrimSpace(pageData.Content) == "" {
		if showAnimations {
			a.ui.StopSpinner()
		}
		return nil, fmt.Errorf("no content could be extracted from the webpage - the page may be empty, require JavaScript, or have unusual structure")
	}
	
	result, err := a.analyzePageData(pageData, url)
	
	if showAnimations {
		a.ui.StopSpinner()
		if err == nil {
			successMsg := a.formatSuccessMessage(result)
			a.ui.PrintSuccess(successMsg)
		}
	}
	
	return result, err
}

func (a *Analyzer) analyzePageData(pageData *webpage.PageData, source string) (*Result, error) {
	result := &Result{
		URL:         source,
		Title:       pageData.Title,
		ProcessedAt: time.Now(),
		Mode:        a.config.Mode,
		Metadata: map[string]any{
			"content_size": len(pageData.Content),
			"meta_tags":    pageData.MetaTags,
			"headings":     pageData.Headings,
		},
	}

	// Always calculate local score
	localScore := a.localScorer.AnalyzeContent(pageData.Content, pageData)
	result.LocalScore = localScore
	result.Score = localScore.Overall
	result.Suggestions = localScore.Suggestions

	switch a.config.Mode {
	case "local":
		// Local-only mode - just use local scoring
		result.Analysis = a.formatLocalAnalysis(localScore)
		result.Metadata["scoring_method"] = "local_only"
		
		// Add LLM recommendation if no API key is available and this was auto mode
		if (a.originalMode == "auto" || a.originalMode == "") && !hasValidAPIKey(a.config.LLMProvider) {
			result.Analysis += a.formatLLMRecommendation()
		}
		
	case "llm":
		// LLM-only mode
		if a.initError != nil {
			return nil, a.initError
		}
		if a.provider == nil {
			return nil, fmt.Errorf("LLM provider not available")
		}
		ctx, cancel := context.WithTimeout(context.Background(), time.Duration(a.config.Timeout)*time.Second)
		defer cancel()
		
		response, err := a.provider.Analyze(ctx, pageData.Content, getGeoPrompt())
		if err != nil {
			return nil, fmt.Errorf("LLM analysis failed: %w", err)
		}
		
		// Extract LLM score and average with local score
		llmScore := extractScoreFromLLMResponse(response.Content)
		if llmScore > 0 {
			// Average local and LLM scores
			result.Score = (localScore.Overall + llmScore) / 2
			result.Metadata["local_score"] = localScore.Overall
			result.Metadata["llm_score"] = llmScore
			result.Metadata["scoring_method"] = "llm_averaged"
		} else {
			// Keep local score if LLM score extraction fails
			result.Metadata["scoring_method"] = "llm_no_score_fallback"
		}
		
		result.Analysis = response.Content
		result.TokensUsed = response.TokensUsed
		result.Metadata["model"] = response.Model
		result.Metadata["provider"] = a.provider.Name()
		
	case "hybrid":
		// Hybrid mode - combine local scoring with LLM insights
		result.Analysis = a.formatLocalAnalysis(localScore)
		
		if a.provider != nil {
			ctx, cancel := context.WithTimeout(context.Background(), time.Duration(a.config.Timeout)*time.Second)
			defer cancel()
			
			hybridPrompt := a.createHybridPrompt(localScore, pageData.Content)
			response, err := a.provider.Analyze(ctx, pageData.Content, hybridPrompt)
			if err == nil {
				// Parse LLM score if available and average with local score
				llmScore := extractScoreFromLLMResponse(response.Content)
				if llmScore > 0 {
					// Average local and LLM scores
					result.Score = (localScore.Overall + llmScore) / 2
					result.Metadata["local_score"] = localScore.Overall
					result.Metadata["llm_score"] = llmScore
					result.Metadata["scoring_method"] = "hybrid_averaged"
				}
				
				result.Analysis += "\n\n" + response.Content
				result.TokensUsed = response.TokensUsed
				result.Metadata["model"] = response.Model
				result.Metadata["provider"] = a.provider.Name()
			} else {
				// In hybrid mode, log LLM errors but don't fail the analysis
				result.Metadata["llm_error"] = err.Error()
				result.Metadata["scoring_method"] = "local_only_fallback"
			}
		} else {
			result.Metadata["scoring_method"] = "local_only"
		}
	}
	
	return result, nil
}

func (a *Analyzer) AnalyzeContent(content, title string) (*Result, error) {
	// Create a minimal PageData for local scoring
	pageData := &webpage.PageData{
		Title:    title,
		Content:  content,
		MetaTags: make(map[string]string),
		Headings: []webpage.Heading{},
	}
	
	return a.analyzePageData(pageData, title)
}

func (a *Analyzer) formatLocalAnalysis(score *scorer.GEOScore) string {
	analysis := fmt.Sprintf("=== Local GEO Analysis ===\n\n")
	analysis += fmt.Sprintf("Overall Score: %d/100\n\n", score.Overall)
	
	analysis += "=== Score Breakdown ===\n"
	analysis += fmt.Sprintf("Content Structure: %d/100 (%.1f%%)\n", 
		score.Breakdown.ContentStructure.Score, score.Breakdown.ContentStructure.Percentage)
	analysis += fmt.Sprintf("Semantic Clarity: %d/100 (%.1f%%)\n", 
		score.Breakdown.SemanticClarity.Score, score.Breakdown.SemanticClarity.Percentage)
	analysis += fmt.Sprintf("Context Richness: %d/100 (%.1f%%)\n", 
		score.Breakdown.ContextRichness.Score, score.Breakdown.ContextRichness.Percentage)
	analysis += fmt.Sprintf("Authority Signals: %d/100 (%.1f%%)\n", 
		score.Breakdown.AuthoritySignals.Score, score.Breakdown.AuthoritySignals.Percentage)
	analysis += fmt.Sprintf("Accessibility: %d/100 (%.1f%%)\n\n", 
		score.Breakdown.Accessibility.Score, score.Breakdown.Accessibility.Percentage)
	
	if len(score.Strengths) > 0 {
		analysis += "=== Strengths ===\n"
		for _, strength := range score.Strengths {
			analysis += fmt.Sprintf("✓ %s\n", strength)
		}
		analysis += "\n"
	}
	
	if len(score.Suggestions) > 0 {
		analysis += "=== Recommendations ===\n"
		for i, suggestion := range score.Suggestions {
			analysis += fmt.Sprintf("%d. %s\n", i+1, suggestion)
		}
		analysis += "\n"
	}
	
	return analysis
}

func (a *Analyzer) createHybridPrompt(localScore *scorer.GEOScore, content string) string {
	prompt := fmt.Sprintf(`Based on the local GEO analysis below, provide additional insights and detailed recommendations for optimizing this content for AI systems.

Local Analysis Results:
- Overall Score: %d/100
- Content Structure: %d/100
- Semantic Clarity: %d/100  
- Context Richness: %d/100
- Authority Signals: %d/100
- Accessibility: %d/100

Key Issues Identified:
`, localScore.Overall, 
	localScore.Breakdown.ContentStructure.Score,
	localScore.Breakdown.SemanticClarity.Score,
	localScore.Breakdown.ContextRichness.Score,
	localScore.Breakdown.AuthoritySignals.Score,
	localScore.Breakdown.Accessibility.Score)

	for _, suggestion := range localScore.Suggestions {
		prompt += fmt.Sprintf("- %s\n", suggestion)
	}

	prompt += `
Please provide:
1. Validation or refinement of the local analysis
2. Specific, actionable recommendations for improvement
3. Examples of how to implement the suggestions
4. Advanced GEO strategies not covered by local analysis

Focus on practical advice for optimizing content for AI understanding and reference.`

	return prompt
}

func getAPIKey(provider string) string {
	switch provider {
	case "claude":
		return os.Getenv("CLAUDE_API_KEY")
	case "gpt", "openai":
		return os.Getenv("OPENAI_API_KEY")
	default:
		return ""
	}
}

// hasValidAPIKey checks if a valid API key is available for the provider
func hasValidAPIKey(provider string) bool {
	apiKey := getAPIKey(provider)
	if apiKey == "" {
		return false
	}
	
	// Basic format validation
	switch provider {
	case "claude":
		return strings.HasPrefix(apiKey, "sk-ant-")
	case "gpt", "openai":
		return strings.HasPrefix(apiKey, "sk-")
	case "local":
		return true // Local doesn't require API key
	default:
		return false
	}
}

// determineOptimalMode automatically selects the best mode based on available API keys
func determineOptimalMode(provider string) string {
	// Check if the specified provider has a valid API key
	if hasValidAPIKey(provider) {
		return "hybrid"
	}
	
	// If specified provider doesn't have a key, check for any available API keys
	availableProviders := []string{"openai", "claude"}
	for _, p := range availableProviders {
		if hasValidAPIKey(p) {
			return "hybrid" // Use hybrid mode with any available provider
		}
	}
	
	return "local" // Fallback to local when no API key is available
}

// extractScoreFromLLMResponse attempts to extract a numerical score from LLM response
func extractScoreFromLLMResponse(content string) int {
	// Look for patterns like "Score: 75/100" or "Overall: 80"
	patterns := []string{
		`(?i)(?:overall|total|final)\s*(?:score|rating)?:?\s*(\d+)(?:/100|%)?`,
		`(?i)(?:score|rating):?\s*(\d+)(?:/100|%)?`,
		`(?i)(\d+)(?:/100|%)\s*(?:overall|total|final)?`,
	}
	
	for _, pattern := range patterns {
		re := regexp.MustCompile(pattern)
		if matches := re.FindStringSubmatch(content); len(matches) > 1 {
			if score, err := strconv.Atoi(matches[1]); err == nil {
				// Ensure score is within valid range
				if score >= 0 && score <= 100 {
					return score
				}
			}
		}
	}
	return 0 // No valid score found
}

// formatLLMRecommendation adds a recommendation to use LLM for better analysis
func (a *Analyzer) formatLLMRecommendation() string {
	return `

## 🤖 Enhanced Analysis Recommendation

This analysis used **local rule-based scoring only**. For more accurate GEO optimization insights, consider using **AI-powered analysis**:

### 🔑 Quick Setup

1. **Set up an API key**:
   - **OpenAI**: ` + "`" + `export OPENAI_API_KEY="your-key"` + "`" + `
   - **Claude**: ` + "`" + `export CLAUDE_API_KEY="your-key"` + "`" + `

2. **Re-run with enhanced analysis**:
   - ` + "`" + `geo-checker analyze <url> --mode hybrid` + "`" + `
   - ` + "`" + `geo-checker analyze <url> --interactive` + "`" + `

### 🎯 Benefits of AI Analysis

- **More accurate content quality assessment**
- **Better understanding of context and meaning**  
- **Averaged scoring** between local + AI for balanced results
- **Specific recommendations** based on actual content value

> **Key Insight**: Local scoring focuses on technical structure, while AI analysis evaluates real content quality that matters for GEO optimization.

### 📊 Analysis Comparison

- **Local Only**: Technical structure → Best for quick audits
- **AI + Local**: Quality + structure → **Best for comprehensive optimization**
- **AI Only**: Content meaning → Best for deep content analysis

> **Recommended**: Use **hybrid mode** for the most balanced and accurate results.
`
}

// formatSuccessMessage creates an informative success message based on scoring method
func (a *Analyzer) formatSuccessMessage(result *Result) string {
	scoringMethod, exists := result.Metadata["scoring_method"].(string)
	if !exists {
		scoringMethod = "unknown"
	}
	
	switch scoringMethod {
	case "hybrid_averaged", "llm_averaged":
		localScore := result.Metadata["local_score"].(int)
		llmScore := result.Metadata["llm_score"].(int)
		return fmt.Sprintf("Analysis complete! Score: %d/100 (Local: %d + AI: %d, averaged)", 
			result.Score, localScore, llmScore)
			
	case "local_only_fallback", "llm_no_score_fallback":
		return fmt.Sprintf("Analysis complete! Score: %d/100 (Local only - AI analysis failed)", 
			result.Score)
			
	case "local_only":
		return fmt.Sprintf("Analysis complete! Score: %d/100 (Local analysis)", 
			result.Score)
			
	default:
		return fmt.Sprintf("Analysis complete! Score: %d/100", result.Score)
	}
}