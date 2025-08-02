package bulk

import (
	"bufio"
	"fmt"
	"geo-checker/pkg/analyzer"
	"geo-checker/pkg/config"
	"geo-checker/pkg/ui"
	"os"
	"strings"
	"sync"
)

type Processor struct {
	config   *config.Config
	analyzer *analyzer.Analyzer
	ui       *ui.UI
}

type BulkResult struct {
	URL     string             `json:"url"`
	Result  *analyzer.Result   `json:"result,omitempty"`
	Error   string             `json:"error,omitempty"`
}

func New(cfg *config.Config) *Processor {
	return &Processor{
		config:   cfg,
		analyzer: analyzer.New(cfg),
		ui:       ui.New(),
	}
}

func (p *Processor) ProcessFile(filename string) ([]*BulkResult, error) {
	urls, err := p.readURLsFromFile(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to read URLs from file: %w", err)
	}
	
	if len(urls) == 0 {
		return nil, fmt.Errorf("no URLs found in file")
	}
	
	return p.ProcessURLs(urls)
}

func (p *Processor) ProcessURLs(urls []string) ([]*BulkResult, error) {
	results := make([]*BulkResult, len(urls))
	
	// Show status messages for text output
	showProgress := p.config.OutputFormat != "json"
	
	var progress *ui.UI
	
	if showProgress {
		progress = ui.New()
		progress.PrintInfo(fmt.Sprintf("Processing %d URLs with %d concurrent workers...", len(urls), p.config.Concurrent))
	}
	
	// Create a semaphore to limit concurrent requests
	semaphore := make(chan struct{}, p.config.Concurrent)
	var wg sync.WaitGroup
	
	for i, url := range urls {
		wg.Add(1)
		go func(index int, u string) {
			defer wg.Done()
			
			// Acquire semaphore
			semaphore <- struct{}{}
			defer func() { <-semaphore }()
			
			result := &BulkResult{URL: u}
			
			analysisResult, err := p.analyzer.AnalyzeURL(u)
			if err != nil {
				result.Error = err.Error()
			} else {
				result.Result = analysisResult
			}
			
			results[index] = result
		}(i, url)
	}
	
	wg.Wait()
	
	if showProgress {
		progress.PrintSuccess(fmt.Sprintf("Completed analysis of %d URLs!", len(urls)))
	}
	
	return results, nil
}

func (p *Processor) readURLsFromFile(filename string) ([]string, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()
	
	var urls []string
	scanner := bufio.NewScanner(file)
	
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line != "" && !strings.HasPrefix(line, "#") {
			// Basic URL validation
			if strings.HasPrefix(line, "http://") || strings.HasPrefix(line, "https://") {
				urls = append(urls, line)
			}
		}
	}
	
	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error reading file: %w", err)
	}
	
	return urls, nil
}