package webpage

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
)

type Scraper struct {
	client *http.Client
}

type PageData struct {
	URL      string            `json:"url"`
	Title    string            `json:"title"`
	Content  string            `json:"content"`
	MetaTags map[string]string `json:"meta_tags"`
	Headings []Heading         `json:"headings"`
}

type Heading struct {
	Level int    `json:"level"`
	Text  string `json:"text"`
}

func New() *Scraper {
	return &Scraper{
		client: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

func (s *Scraper) ScrapeURL(ctx context.Context, url string) (*PageData, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	
	req.Header.Set("User-Agent", "GEO-Checker/1.0 (+https://github.com/your-repo/geo-checker)")
	
	resp, err := s.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch URL: %w", err)
	}
	defer resp.Body.Close()
	
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("HTTP error: %d", resp.StatusCode)
	}
	
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}
	
	return s.parseHTML(string(body), url)
}

func (s *Scraper) parseHTML(html, source string) (*PageData, error) {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(html))
	if err != nil {
		return nil, fmt.Errorf("failed to parse HTML: %w", err)
	}
	
	pageData := &PageData{
		URL:      source,
		MetaTags: make(map[string]string),
		Headings: []Heading{},
	}
	
	// Extract title
	pageData.Title = doc.Find("title").Text()
	
	// Extract meta tags
	doc.Find("meta").Each(func(i int, s *goquery.Selection) {
		if name, exists := s.Attr("name"); exists {
			content, _ := s.Attr("content")
			pageData.MetaTags[name] = content
		}
		if property, exists := s.Attr("property"); exists {
			content, _ := s.Attr("content")
			pageData.MetaTags[property] = content
		}
	})
	
	// Extract headings
	doc.Find("h1, h2, h3, h4, h5, h6").Each(func(i int, s *goquery.Selection) {
		level := getHeadingLevel(s.Get(0).Data)
		text := strings.TrimSpace(s.Text())
		if text != "" {
			pageData.Headings = append(pageData.Headings, Heading{
				Level: level,
				Text:  text,
			})
		}
	})
	
	// Extract main content
	content := s.extractContent(doc)
	pageData.Content = strings.TrimSpace(content)
	
	// Validate that we have some content
	if pageData.Content == "" {
		// If no content extracted, create minimal content from available data
		var fallbackContent strings.Builder
		if pageData.Title != "" {
			fallbackContent.WriteString("Page Title: " + pageData.Title + "\n\n")
		}
		
		if len(pageData.Headings) > 0 {
			fallbackContent.WriteString("Page Headings:\n")
			for _, heading := range pageData.Headings {
				fallbackContent.WriteString(fmt.Sprintf("H%d: %s\n", heading.Level, heading.Text))
			}
			fallbackContent.WriteString("\n")
		}
		
		if len(pageData.MetaTags) > 0 {
			if desc, exists := pageData.MetaTags["description"]; exists && desc != "" {
				fallbackContent.WriteString("Meta Description: " + desc + "\n\n")
			}
		}
		
		fallbackText := fallbackContent.String()
		if fallbackText != "" {
			pageData.Content = fallbackText
		} else {
			// Absolute fallback
			pageData.Content = fmt.Sprintf("Webpage at %s - Content extraction failed, only metadata available.", source)
		}
	}
	
	return pageData, nil
}

func (s *Scraper) extractContent(doc *goquery.Document) string {
	// Remove script and style elements
	doc.Find("script, style, nav, footer, header, aside").Remove()
	
	var content strings.Builder
	
	// Extract main content areas
	mainSelectors := []string{
		"main",
		"article",
		"[role=\"main\"]",
		".content",
		".main-content",
		"#content",
		"#main",
	}
	
	var mainContent *goquery.Selection
	for _, selector := range mainSelectors {
		if sel := doc.Find(selector); sel.Length() > 0 {
			mainContent = sel.First()
			break
		}
	}
	
	if mainContent == nil {
		mainContent = doc.Find("body")
	}
	
	// Extract text content
	mainContent.Find("h1, h2, h3, h4, h5, h6, p, li, td, th, blockquote, pre").Each(func(i int, s *goquery.Selection) {
		text := strings.TrimSpace(s.Text())
		if text != "" {
			content.WriteString(text)
			content.WriteString("\n\n")
		}
	})
	
	// Fallback: if no content found with specific selectors, try to get all text from body
	if content.Len() == 0 {
		bodyText := strings.TrimSpace(doc.Find("body").Text())
		// Clean up excessive whitespace
		bodyText = strings.Join(strings.Fields(bodyText), " ")
		if len(bodyText) > 50 { // Only use if substantial content
			content.WriteString(bodyText)
		}
	}
	
	// Final fallback: use title and headings if no other content
	if content.Len() == 0 {
		title := doc.Find("title").Text()
		if title != "" {
			content.WriteString("Title: " + title + "\n\n")
		}
		
		doc.Find("h1, h2, h3, h4, h5, h6").Each(func(i int, s *goquery.Selection) {
			heading := strings.TrimSpace(s.Text())
			if heading != "" {
				content.WriteString(heading + "\n")
			}
		})
	}
	
	return content.String()
}

func getHeadingLevel(tagName string) int {
	switch tagName {
	case "h1":
		return 1
	case "h2":
		return 2
	case "h3":
		return 3
	case "h4":
		return 4
	case "h5":
		return 5
	case "h6":
		return 6
	default:
		return 0
	}
}

func readFile(filePath string) (string, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return "", fmt.Errorf("failed to read file %s: %w", filePath, err)
	}
	return string(data), nil
}