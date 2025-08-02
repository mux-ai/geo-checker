package scorer

import (
	"geo-checker/internal/webpage"
	"math"
	"regexp"
	"strings"
)

type LocalScorer struct {
	weights GEOWeights
}

type GEOWeights struct {
	ContentStructure float64
	SemanticClarity  float64
	ContextRichness  float64
	AuthoritySignals float64
	Accessibility    float64
}

type GEOScore struct {
	Overall          int                    `json:"overall_score"`
	Breakdown        ScoreBreakdown         `json:"breakdown"`
	Suggestions      []string               `json:"suggestions"`
	Strengths        []string               `json:"strengths"`
	Weaknesses       []string               `json:"weaknesses"`
	Metadata         map[string]interface{} `json:"metadata"`
}

type ScoreBreakdown struct {
	ContentStructure ScoreDetail `json:"content_structure"`
	SemanticClarity  ScoreDetail `json:"semantic_clarity"`
	ContextRichness  ScoreDetail `json:"context_richness"`
	AuthoritySignals ScoreDetail `json:"authority_signals"`
	Accessibility    ScoreDetail `json:"accessibility"`
}

type ScoreDetail struct {
	Score       int      `json:"score"`
	MaxScore    int      `json:"max_score"`
	Percentage  float64  `json:"percentage"`
	Issues      []string `json:"issues"`
	Positives   []string `json:"positives"`
}

func NewLocalScorer() *LocalScorer {
	return &LocalScorer{
		weights: GEOWeights{
			ContentStructure: 0.25,
			SemanticClarity:  0.25,
			ContextRichness:  0.20,
			AuthoritySignals: 0.15,
			Accessibility:    0.15,
		},
	}
}

func (ls *LocalScorer) AnalyzeContent(content string, pageData *webpage.PageData) *GEOScore {
	score := &GEOScore{
		Breakdown:   ScoreBreakdown{},
		Suggestions: []string{},
		Strengths:   []string{},
		Weaknesses:  []string{},
		Metadata:    make(map[string]interface{}),
	}

	// Analyze each component
	score.Breakdown.ContentStructure = ls.analyzeContentStructure(content, pageData)
	score.Breakdown.SemanticClarity = ls.analyzeSemanticClarity(content)
	score.Breakdown.ContextRichness = ls.analyzeContextRichness(content, pageData)
	score.Breakdown.AuthoritySignals = ls.analyzeAuthoritySignals(content, pageData)
	score.Breakdown.Accessibility = ls.analyzeAccessibility(content, pageData)

	// Calculate overall score
	score.Overall = ls.calculateOverallScore(score.Breakdown)

	// Generate suggestions and insights
	ls.generateInsights(score)

	// Add metadata
	score.Metadata["content_length"] = len(content)
	score.Metadata["word_count"] = len(strings.Fields(content))
	score.Metadata["heading_count"] = len(pageData.Headings)
	score.Metadata["meta_tags_count"] = len(pageData.MetaTags)

	return score
}

func (ls *LocalScorer) analyzeContentStructure(content string, pageData *webpage.PageData) ScoreDetail {
	detail := ScoreDetail{MaxScore: 100, Issues: []string{}, Positives: []string{}}
	score := 0

	// Check heading hierarchy (30 points)
	headingScore := ls.evaluateHeadingHierarchy(pageData.Headings)
	score += headingScore
	if headingScore >= 25 {
		detail.Positives = append(detail.Positives, "Good heading hierarchy structure")
	} else {
		detail.Issues = append(detail.Issues, "Improve heading hierarchy (H1 → H2 → H3)")
	}

	// Check content organization (25 points)
	orgScore := ls.evaluateContentOrganization(content)
	score += orgScore
	if orgScore >= 20 {
		detail.Positives = append(detail.Positives, "Well-organized content structure")
	} else {
		detail.Issues = append(detail.Issues, "Content could be better organized with clear sections")
	}

	// Check paragraph structure (25 points)
	paraScore := ls.evaluateParagraphStructure(content)
	score += paraScore
	if paraScore >= 20 {
		detail.Positives = append(detail.Positives, "Good paragraph structure")
	} else {
		detail.Issues = append(detail.Issues, "Use shorter, more focused paragraphs")
	}

	// Check list usage (20 points)
	listScore := ls.evaluateListUsage(content)
	score += listScore
	if listScore >= 15 {
		detail.Positives = append(detail.Positives, "Effective use of lists for organization")
	} else {
		detail.Issues = append(detail.Issues, "Consider using lists to organize key points")
	}

	detail.Score = score
	detail.Percentage = float64(score) / float64(detail.MaxScore) * 100
	return detail
}

func (ls *LocalScorer) analyzeSemanticClarity(content string) ScoreDetail {
	detail := ScoreDetail{MaxScore: 100, Issues: []string{}, Positives: []string{}}
	score := 0

	// Check readability (40 points)
	readScore := ls.evaluateReadability(content)
	score += readScore
	if readScore >= 30 {
		detail.Positives = append(detail.Positives, "Content is clear and readable")
	} else {
		detail.Issues = append(detail.Issues, "Simplify sentence structure for better readability")
	}

	// Check terminology consistency (30 points)
	termScore := ls.evaluateTerminologyConsistency(content)
	score += termScore
	if termScore >= 25 {
		detail.Positives = append(detail.Positives, "Consistent terminology usage")
	} else {
		detail.Issues = append(detail.Issues, "Use consistent terminology throughout")
	}

	// Check definition clarity (30 points)
	defScore := ls.evaluateDefinitionClarity(content)
	score += defScore
	if defScore >= 25 {
		detail.Positives = append(detail.Positives, "Clear definitions and explanations")
	} else {
		detail.Issues = append(detail.Issues, "Define technical terms and concepts clearly")
	}

	detail.Score = score
	detail.Percentage = float64(score) / float64(detail.MaxScore) * 100
	return detail
}

func (ls *LocalScorer) analyzeContextRichness(content string, pageData *webpage.PageData) ScoreDetail {
	detail := ScoreDetail{MaxScore: 100, Issues: []string{}, Positives: []string{}}
	score := 0

	// Check content depth (40 points)
	depthScore := ls.evaluateContentDepth(content)
	score += depthScore
	if depthScore >= 30 {
		detail.Positives = append(detail.Positives, "Rich, detailed content")
	} else {
		detail.Issues = append(detail.Issues, "Add more detailed explanations and examples")
	}

	// Check examples and specifics (35 points)
	exampleScore := ls.evaluateExamplesAndSpecifics(content)
	score += exampleScore
	if exampleScore >= 25 {
		detail.Positives = append(detail.Positives, "Good use of examples and specific details")
	} else {
		detail.Issues = append(detail.Issues, "Include more concrete examples and specific details")
	}

	// Check background information (25 points)
	backgroundScore := ls.evaluateBackgroundInfo(content)
	score += backgroundScore
	if backgroundScore >= 20 {
		detail.Positives = append(detail.Positives, "Adequate background information provided")
	} else {
		detail.Issues = append(detail.Issues, "Provide more context and background information")
	}

	detail.Score = score
	detail.Percentage = float64(score) / float64(detail.MaxScore) * 100
	return detail
}

func (ls *LocalScorer) analyzeAuthoritySignals(content string, pageData *webpage.PageData) ScoreDetail {
	detail := ScoreDetail{MaxScore: 100, Issues: []string{}, Positives: []string{}}
	score := 0

	// Check citations and references (40 points)
	citationScore := ls.evaluateCitations(content)
	score += citationScore
	if citationScore >= 30 {
		detail.Positives = append(detail.Positives, "Good use of citations and references")
	} else {
		detail.Issues = append(detail.Issues, "Add more citations and credible references")
	}

	// Check expertise indicators (35 points)
	expertiseScore := ls.evaluateExpertiseIndicators(content)
	score += expertiseScore
	if expertiseScore >= 25 {
		detail.Positives = append(detail.Positives, "Clear expertise and authority indicators")
	} else {
		detail.Issues = append(detail.Issues, "Include more expertise and credibility signals")
	}

	// Check factual accuracy indicators (25 points)
	factScore := ls.evaluateFactualAccuracy(content)
	score += factScore
	if factScore >= 20 {
		detail.Positives = append(detail.Positives, "Content appears factual and well-researched")
	} else {
		detail.Issues = append(detail.Issues, "Ensure factual accuracy and provide sources")
	}

	detail.Score = score
	detail.Percentage = float64(score) / float64(detail.MaxScore) * 100
	return detail
}

func (ls *LocalScorer) analyzeAccessibility(content string, pageData *webpage.PageData) ScoreDetail {
	detail := ScoreDetail{MaxScore: 100, Issues: []string{}, Positives: []string{}}
	score := 0

	// Check meta information (30 points)
	metaScore := ls.evaluateMetaInformation(pageData)
	score += metaScore
	if metaScore >= 25 {
		detail.Positives = append(detail.Positives, "Good meta information for AI understanding")
	} else {
		detail.Issues = append(detail.Issues, "Add comprehensive meta descriptions and keywords")
	}

	// Check content parsing friendliness (35 points)
	parseScore := ls.evaluateParsingFriendliness(content)
	score += parseScore
	if parseScore >= 25 {
		detail.Positives = append(detail.Positives, "Content is easy to parse and understand")
	} else {
		detail.Issues = append(detail.Issues, "Structure content for better machine readability")
	}

	// Check information density (35 points)
	densityScore := ls.evaluateInformationDensity(content)
	score += densityScore
	if densityScore >= 25 {
		detail.Positives = append(detail.Positives, "Good information density")
	} else {
		detail.Issues = append(detail.Issues, "Balance information density - avoid being too sparse or dense")
	}

	detail.Score = score
	detail.Percentage = float64(score) / float64(detail.MaxScore) * 100
	return detail
}

// Helper functions for evaluation
func (ls *LocalScorer) evaluateHeadingHierarchy(headings []webpage.Heading) int {
	if len(headings) == 0 {
		return 0
	}

	score := 10 // Base score for having headings
	
	// Check for H1
	hasH1 := false
	for _, h := range headings {
		if h.Level == 1 {
			hasH1 = true
			break
		}
	}
	if hasH1 {
		score += 10
	}

	// Check hierarchy flow
	if len(headings) > 1 {
		hierarchyGood := true
		for i := 1; i < len(headings); i++ {
			if headings[i].Level > headings[i-1].Level+1 {
				hierarchyGood = false
				break
			}
		}
		if hierarchyGood {
			score += 10
		}
	}

	return min(score, 30)
}

func (ls *LocalScorer) evaluateContentOrganization(content string) int {
	sections := strings.Split(content, "\n\n")
	if len(sections) < 2 {
		return 5
	}

	score := 10
	if len(sections) >= 3 {
		score += 10
	}
	if len(sections) >= 5 {
		score += 5
	}

	return min(score, 25)
}

func (ls *LocalScorer) evaluateParagraphStructure(content string) int {
	paragraphs := strings.Split(content, "\n\n")
	score := 0
	
	goodParagraphs := 0
	for _, para := range paragraphs {
		words := len(strings.Fields(para))
		if words >= 20 && words <= 150 {
			goodParagraphs++
		}
	}

	if len(paragraphs) > 0 {
		ratio := float64(goodParagraphs) / float64(len(paragraphs))
		score = int(ratio * 25)
	}

	return score
}

func (ls *LocalScorer) evaluateListUsage(content string) int {
	// Simple check for list indicators
	listIndicators := []string{"•", "-", "*", "1.", "2.", "3.", "①", "②", "③"}
	listCount := 0
	
	for _, indicator := range listIndicators {
		listCount += strings.Count(content, indicator)
	}

	if listCount == 0 {
		return 5
	} else if listCount <= 3 {
		return 10
	} else if listCount <= 8 {
		return 20
	}
	return 20
}

func (ls *LocalScorer) evaluateReadability(content string) int {
	words := strings.Fields(content)
	if len(words) == 0 {
		return 0
	}

	// Simple readability metrics
	avgWordsPerSentence := ls.calculateAvgWordsPerSentence(content)
	avgSyllablesPerWord := ls.calculateAvgSyllablesPerWord(words)

	score := 20 // Base score

	// Prefer 15-20 words per sentence
	if avgWordsPerSentence >= 10 && avgWordsPerSentence <= 25 {
		score += 10
	}

	// Prefer 1-3 syllables per word average
	if avgSyllablesPerWord >= 1.0 && avgSyllablesPerWord <= 2.5 {
		score += 10
	}

	return min(score, 40)
}

func (ls *LocalScorer) evaluateTerminologyConsistency(content string) int {
	// Simple consistency check - could be enhanced
	words := strings.Fields(strings.ToLower(content))
	wordCount := make(map[string]int)
	
	for _, word := range words {
		if len(word) > 4 { // Focus on longer words
			wordCount[word]++
		}
	}

	// Check for consistent usage of key terms
	consistentTerms := 0
	totalKeyTerms := 0
	
	for _, count := range wordCount {
		if count >= 3 { // Word appears multiple times
			totalKeyTerms++
			if count >= 3 {
				consistentTerms++
			}
		}
	}

	if totalKeyTerms == 0 {
		return 15
	}

	ratio := float64(consistentTerms) / float64(totalKeyTerms)
	return int(ratio * 30)
}

func (ls *LocalScorer) evaluateDefinitionClarity(content string) int {
	// Look for definition patterns
	definitionPatterns := []string{
		" is ", " are ", " means ", " refers to ", " defined as ",
		"definition", "meaning", "explanation", "concept",
	}

	definitionCount := 0
	for _, pattern := range definitionPatterns {
		definitionCount += strings.Count(strings.ToLower(content), pattern)
	}

	if definitionCount == 0 {
		return 10
	} else if definitionCount <= 5 {
		return 20
	} else if definitionCount <= 10 {
		return 30
	}
	return 30
}

func (ls *LocalScorer) evaluateContentDepth(content string) int {
	wordCount := len(strings.Fields(content))
	
	if wordCount < 100 {
		return 5
	} else if wordCount < 300 {
		return 15
	} else if wordCount < 800 {
		return 30
	} else if wordCount < 1500 {
		return 40
	}
	return 35 // Very long content might be too dense
}

func (ls *LocalScorer) evaluateExamplesAndSpecifics(content string) int {
	examplePatterns := []string{
		"example", "for instance", "such as", "including", "like",
		"specifically", "particular", "namely", "e.g.", "i.e.",
	}

	exampleCount := 0
	for _, pattern := range examplePatterns {
		exampleCount += strings.Count(strings.ToLower(content), pattern)
	}

	if exampleCount == 0 {
		return 5
	} else if exampleCount <= 3 {
		return 15
	} else if exampleCount <= 8 {
		return 25
	} else if exampleCount <= 15 {
		return 35
	}
	return 30
}

func (ls *LocalScorer) evaluateBackgroundInfo(content string) int {
	backgroundPatterns := []string{
		"background", "context", "history", "overview", "introduction",
		"originally", "previously", "traditionally", "historically",
	}

	backgroundCount := 0
	for _, pattern := range backgroundPatterns {
		backgroundCount += strings.Count(strings.ToLower(content), pattern)
	}

	if backgroundCount == 0 {
		return 5
	} else if backgroundCount <= 3 {
		return 15
	} else if backgroundCount <= 6 {
		return 25
	}
	return 25
}

func (ls *LocalScorer) evaluateCitations(content string) int {
	citationPatterns := []string{
		"according to", "research shows", "study found", "source:",
		"reference", "cited", "published", "journal", "doi:",
		"http://", "https://", "www.", ".com", ".org", ".edu",
	}

	citationCount := 0
	for _, pattern := range citationPatterns {
		citationCount += strings.Count(strings.ToLower(content), pattern)
	}

	if citationCount == 0 {
		return 5
	} else if citationCount <= 3 {
		return 15
	} else if citationCount <= 8 {
		return 25
	} else if citationCount <= 15 {
		return 40
	}
	return 35
}

func (ls *LocalScorer) evaluateExpertiseIndicators(content string) int {
	expertisePatterns := []string{
		"expert", "professional", "certified", "experienced", "qualified",
		"research", "analysis", "methodology", "findings", "conclusion",
		"peer-reviewed", "academic", "scholarly", "evidence-based",
	}

	expertiseCount := 0
	for _, pattern := range expertisePatterns {
		expertiseCount += strings.Count(strings.ToLower(content), pattern)
	}

	if expertiseCount == 0 {
		return 10
	} else if expertiseCount <= 3 {
		return 20
	} else if expertiseCount <= 8 {
		return 35
	}
	return 35
}

func (ls *LocalScorer) evaluateFactualAccuracy(content string) int {
	// Look for hedging language that might indicate uncertainty
	uncertaintyPatterns := []string{
		"might", "could", "possibly", "perhaps", "maybe", "seems",
		"appears", "likely", "probably", "allegedly", "reportedly",
	}

	factualPatterns := []string{
		"fact", "proven", "demonstrated", "confirmed", "verified",
		"established", "documented", "evidence", "data", "statistics",
	}

	uncertaintyCount := 0
	factualCount := 0

	contentLower := strings.ToLower(content)
	for _, pattern := range uncertaintyPatterns {
		uncertaintyCount += strings.Count(contentLower, pattern)
	}
	for _, pattern := range factualPatterns {
		factualCount += strings.Count(contentLower, pattern)
	}

	// Prefer more factual language, less uncertainty
	score := 15 // Base score
	if factualCount > uncertaintyCount {
		score += 10
	}
	if factualCount >= 3 {
		score += 5
	}

	return min(score, 25)
}

func (ls *LocalScorer) evaluateMetaInformation(pageData *webpage.PageData) int {
	score := 0

	if pageData.Title != "" {
		score += 10
	}

	if desc, exists := pageData.MetaTags["description"]; exists && len(desc) > 50 {
		score += 10
	}

	if keywords, exists := pageData.MetaTags["keywords"]; exists && len(keywords) > 10 {
		score += 5
	}

	if len(pageData.MetaTags) >= 3 {
		score += 5
	}

	return min(score, 30)
}

func (ls *LocalScorer) evaluateParsingFriendliness(content string) int {
	score := 15 // Base score

	// Check for clear sentence structure
	sentences := strings.Split(content, ".")
	if len(sentences) > 3 {
		score += 10
	}

	// Check for consistent formatting
	if strings.Contains(content, "\n\n") {
		score += 10
	}

	return min(score, 35)
}

func (ls *LocalScorer) evaluateInformationDensity(content string) int {
	words := strings.Fields(content)
	sentences := strings.Split(content, ".")
	
	if len(sentences) == 0 {
		return 0
	}

	avgWordsPerSentence := float64(len(words)) / float64(len(sentences))

	// Optimal range: 12-20 words per sentence
	if avgWordsPerSentence >= 10 && avgWordsPerSentence <= 25 {
		return 35
	} else if avgWordsPerSentence >= 8 && avgWordsPerSentence <= 30 {
		return 25
	} else if avgWordsPerSentence >= 5 && avgWordsPerSentence <= 35 {
		return 15
	}
	return 10
}

// Utility functions
func (ls *LocalScorer) calculateAvgWordsPerSentence(content string) float64 {
	sentences := regexp.MustCompile(`[.!?]+`).Split(content, -1)
	words := strings.Fields(content)
	
	if len(sentences) == 0 {
		return 0
	}
	
	return float64(len(words)) / float64(len(sentences))
}

func (ls *LocalScorer) calculateAvgSyllablesPerWord(words []string) float64 {
	totalSyllables := 0
	for _, word := range words {
		totalSyllables += ls.countSyllables(word)
	}
	
	if len(words) == 0 {
		return 0
	}
	
	return float64(totalSyllables) / float64(len(words))
}

func (ls *LocalScorer) countSyllables(word string) int {
	word = strings.ToLower(word)
	vowels := "aeiouy"
	syllables := 0
	prevWasVowel := false
	
	for _, char := range word {
		isVowel := strings.ContainsRune(vowels, char)
		if isVowel && !prevWasVowel {
			syllables++
		}
		prevWasVowel = isVowel
	}
	
	// Handle silent e
	if strings.HasSuffix(word, "e") && syllables > 1 {
		syllables--
	}
	
	if syllables == 0 {
		syllables = 1
	}
	
	return syllables
}

func (ls *LocalScorer) calculateOverallScore(breakdown ScoreBreakdown) int {
	weightedScore := 0.0
	
	weightedScore += float64(breakdown.ContentStructure.Score) * ls.weights.ContentStructure
	weightedScore += float64(breakdown.SemanticClarity.Score) * ls.weights.SemanticClarity
	weightedScore += float64(breakdown.ContextRichness.Score) * ls.weights.ContextRichness
	weightedScore += float64(breakdown.AuthoritySignals.Score) * ls.weights.AuthoritySignals
	weightedScore += float64(breakdown.Accessibility.Score) * ls.weights.Accessibility
	
	return int(math.Round(weightedScore))
}

func (ls *LocalScorer) generateInsights(score *GEOScore) {
	// Collect all strengths and weaknesses
	allDetails := []ScoreDetail{
		score.Breakdown.ContentStructure,
		score.Breakdown.SemanticClarity,
		score.Breakdown.ContextRichness,
		score.Breakdown.AuthoritySignals,
		score.Breakdown.Accessibility,
	}

	for _, detail := range allDetails {
		score.Strengths = append(score.Strengths, detail.Positives...)
		for _, issue := range detail.Issues {
			score.Suggestions = append(score.Suggestions, issue)
		}
		if detail.Percentage < 50 {
			score.Weaknesses = append(score.Weaknesses, detail.Issues...)
		}
	}

	// Add overall suggestions based on score
	if score.Overall < 60 {
		score.Suggestions = append(score.Suggestions, "Consider comprehensive content restructuring for better GEO optimization")
	} else if score.Overall < 80 {
		score.Suggestions = append(score.Suggestions, "Focus on improving the lowest-scoring areas identified above")
	}
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}