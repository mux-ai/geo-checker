package scanner

import (
	"fmt"
	"geo-checker/pkg/analyzer"
	"geo-checker/pkg/config"
	"geo-checker/pkg/ui"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

type Scanner struct {
	config   *config.Config
	analyzer *analyzer.Analyzer
	ui       *ui.UI
}

type ScanResult struct {
	FilePath string           `json:"file_path"`
	Result   *analyzer.Result `json:"result,omitempty"`
	Error    string           `json:"error,omitempty"`
}

func New(cfg *config.Config) *Scanner {
	return &Scanner{
		config:   cfg,
		analyzer: analyzer.New(cfg),
		ui:       ui.New(),
	}
}

func (s *Scanner) ScanDirectory(dirPath string) ([]*ScanResult, error) {
	var results []*ScanResult
	var filesToScan []string
	
	showProgress := s.config.OutputFormat != "json"
	
	if showProgress {
		s.ui.StartSpinner("Discovering files...")
	}
	
	// First pass: discover all files to scan
	err := filepath.WalkDir(dirPath, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		
		if d.IsDir() {
			return nil
		}
		
		if s.shouldScanFile(path) {
			filesToScan = append(filesToScan, path)
		}
		
		return nil
	})
	
	if err != nil {
		if showProgress {
			s.ui.StopSpinner()
		}
		return nil, fmt.Errorf("failed to walk directory: %w", err)
	}
	
	if showProgress {
		s.ui.StopSpinner()
		s.ui.PrintInfo(fmt.Sprintf("Found %d files to analyze", len(filesToScan)))
	}
	
	if len(filesToScan) == 0 {
		if showProgress {
			s.ui.PrintWarning("No matching files found")
		}
		return results, nil
	}
	
	// Second pass: analyze files
	for _, path := range filesToScan {
		result := s.scanFile(path)
		results = append(results, result)
	}
	
	if showProgress {
		successCount := 0
		errorCount := 0
		totalScore := 0
		
		for _, result := range results {
			if result.Error != "" {
				errorCount++
			} else if result.Result != nil {
				successCount++
				totalScore += result.Result.Score
			}
		}
		
		s.ui.PrintSuccess(fmt.Sprintf("Scan complete! Processed %d files", len(filesToScan)))
		
		if successCount > 0 {
			avgScore := totalScore / successCount
			s.ui.PrintInfo(fmt.Sprintf("Average GEO Score: %d/100", avgScore))
		}
		
		if errorCount > 0 {
			s.ui.PrintWarning(fmt.Sprintf("%d files had errors", errorCount))
		}
	}
	
	return results, nil
}

func (s *Scanner) scanFile(filePath string) *ScanResult {
	result := &ScanResult{FilePath: filePath}
	
	content, err := s.readHTMLFile(filePath)
	if err != nil {
		result.Error = fmt.Sprintf("failed to read file: %v", err)
		return result
	}
	
	title := s.extractTitleFromPath(filePath)
	analysisResult, err := s.analyzer.AnalyzeContent(content, title)
	if err != nil {
		result.Error = fmt.Sprintf("failed to analyze content: %v", err)
		return result
	}
	
	result.Result = analysisResult
	return result
}

func (s *Scanner) shouldScanFile(filePath string) bool {
	ext := strings.ToLower(filepath.Ext(filePath))
	
	for _, allowedExt := range s.config.Extensions {
		if ext == strings.ToLower(allowedExt) {
			return true
		}
	}
	
	return false
}

func (s *Scanner) readHTMLFile(filePath string) (string, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return "", err
	}
	
	// For HTML files, we might want to extract just the text content
	content := string(data)
	
	// Basic HTML content extraction (could be enhanced)
	content = s.extractTextFromHTML(content)
	
	return content, nil
}

func (s *Scanner) extractTextFromHTML(html string) string {
	// Simple text extraction - remove common HTML tags
	// This is a basic implementation; for better results, we could use goquery
	
	// Remove script and style content
	html = removeTagContent(html, "script")
	html = removeTagContent(html, "style")
	
	// Remove HTML tags but keep content
	html = removeTags(html)
	
	// Clean up whitespace
	lines := strings.Split(html, "\n")
	var cleanLines []string
	
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line != "" {
			cleanLines = append(cleanLines, line)
		}
	}
	
	return strings.Join(cleanLines, "\n")
}

func (s *Scanner) extractTitleFromPath(filePath string) string {
	base := filepath.Base(filePath)
	ext := filepath.Ext(base)
	return strings.TrimSuffix(base, ext)
}

func removeTagContent(html, tag string) string {
	startTag := fmt.Sprintf("<%s", tag)
	endTag := fmt.Sprintf("</%s>", tag)
	
	for {
		start := strings.Index(strings.ToLower(html), strings.ToLower(startTag))
		if start == -1 {
			break
		}
		
		// Find the end of the opening tag
		tagEnd := strings.Index(html[start:], ">")
		if tagEnd == -1 {
			break
		}
		tagEnd += start + 1
		
		// Find the closing tag
		end := strings.Index(strings.ToLower(html[tagEnd:]), strings.ToLower(endTag))
		if end == -1 {
			break
		}
		end += tagEnd + len(endTag)
		
		html = html[:start] + html[end:]
	}
	
	return html
}

func removeTags(html string) string {
	inTag := false
	var result strings.Builder
	
	for _, char := range html {
		if char == '<' {
			inTag = true
		} else if char == '>' {
			inTag = false
		} else if !inTag {
			result.WriteRune(char)
		}
	}
	
	return result.String()
}