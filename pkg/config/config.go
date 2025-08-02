package config

type Config struct {
	LLMProvider   string
	Model         string
	OutputFormat  string
	Mode          string // "local", "llm", "hybrid"
	Concurrent    int
	Extensions    []string
	
	// API Keys
	ClaudeAPIKey  string
	OpenAIAPIKey  string
	LocalLLMURL   string
	
	// Analysis settings
	MaxTokens     int
	Temperature   float64
	Timeout       int
}

