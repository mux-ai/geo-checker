package ui

import (
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/fatih/color"
	"github.com/briandowns/spinner"
)

var (
	// Colors
	Success     = color.New(color.FgGreen, color.Bold)
	Error       = color.New(color.FgRed, color.Bold)
	Warning     = color.New(color.FgYellow, color.Bold)
	Info        = color.New(color.FgCyan, color.Bold)
	Header      = color.New(color.FgMagenta, color.Bold)
	Subtle      = color.New(color.FgHiBlack)
	Score       = color.New(color.FgWhite, color.Bold)
	
	// Themed colors
	Primary     = color.New(color.FgBlue, color.Bold)
	Secondary   = color.New(color.FgHiBlue)
	Accent      = color.New(color.FgHiCyan)
	
	// Markdown-style colors
	H1          = color.New(color.FgHiMagenta, color.Bold)
	H2          = color.New(color.FgHiBlue, color.Bold)
	H3          = color.New(color.FgHiCyan, color.Bold)
	CodeBlock   = color.New(color.BgHiBlack, color.FgWhite)
	InlineCode  = color.New(color.FgHiYellow)
	Quote       = color.New(color.FgHiBlack, color.Italic)
	Link        = color.New(color.FgBlue, color.Underline)
	Bold        = color.New(color.Bold)
	Italic      = color.New(color.Italic)
	ListItem    = color.New(color.FgHiGreen)
	Checkmark   = color.New(color.FgGreen, color.Bold)
	Cross       = color.New(color.FgRed, color.Bold)
)

type UI struct {
	spinner *spinner.Spinner
	NoColor bool
}

func New() *UI {
	s := spinner.New(spinner.CharSets[14], 100*time.Millisecond)
	s.Color("cyan")
	
	return &UI{
		spinner: s,
		NoColor: false,
	}
}

func (ui *UI) StartSpinner(message string) {
	if !ui.NoColor {
		ui.spinner.Suffix = " " + message
		ui.spinner.Start()
	} else {
		fmt.Print(message + "...")
	}
}

func (ui *UI) UpdateSpinner(message string) {
	if !ui.NoColor {
		ui.spinner.Suffix = " " + message
	}
}

func (ui *UI) StopSpinner() {
	if !ui.NoColor {
		ui.spinner.Stop()
	}
}

func (ui *UI) PrintHeader(title string) {
	if ui.NoColor {
		fmt.Printf("=== %s ===\n", title)
		return
	}
	
	// Calculate box width - minimum 60 characters for consistency
	minWidth := 60
	titleWidth := len(title)
	boxWidth := minWidth
	if titleWidth+6 > minWidth {
		boxWidth = titleWidth + 6
	}
	
	// Calculate padding to center the title
	totalPadding := boxWidth - titleWidth - 2 // -2 for the side borders
	leftPadding := totalPadding / 2
	rightPadding := totalPadding - leftPadding
	
	border := strings.Repeat("═", boxWidth-2)
	leftSpaces := strings.Repeat(" ", leftPadding)
	rightSpaces := strings.Repeat(" ", rightPadding)
	
	Header.Printf("╔%s╗\n", border)
	Header.Printf("║%s%s%s║\n", leftSpaces, title, rightSpaces)
	Header.Printf("╚%s╝\n", border)
	fmt.Println()
}

func (ui *UI) PrintSuccess(message string) {
	if ui.NoColor {
		fmt.Printf("✓ %s\n", message)
	} else {
		Success.Printf("✓ %s\n", message)
	}
}

func (ui *UI) PrintError(message string) {
	if ui.NoColor {
		fmt.Printf("✗ %s\n", message)
	} else {
		Error.Printf("✗ %s\n", message)
	}
}

func (ui *UI) PrintWarning(message string) {
	if ui.NoColor {
		fmt.Printf("⚠ %s\n", message)
	} else {
		Warning.Printf("⚠ %s\n", message)
	}
}

func (ui *UI) PrintInfo(message string) {
	if ui.NoColor {
		fmt.Printf("ℹ %s\n", message)
	} else {
		Info.Printf("ℹ %s\n", message)
	}
}

func (ui *UI) PrintScore(label string, score int, maxScore int) {
	percentage := float64(score) / float64(maxScore) * 100
	
	if ui.NoColor {
		fmt.Printf("  %-20s %3d/%-3d (%.1f%%)\n", label+":", score, maxScore, percentage)
		return
	}
	
	// Color based on score
	var scoreColor *color.Color
	if percentage >= 80 {
		scoreColor = Success
	} else if percentage >= 60 {
		scoreColor = Warning
	} else {
		scoreColor = Error
	}
	
	// Aligned output with consistent spacing
	fmt.Printf("  %-20s ", label+":")
	scoreColor.Printf("%3d", score)
	fmt.Printf("/")
	scoreColor.Printf("%-3d", maxScore)
	Subtle.Printf(" (%.1f%%)\n", percentage)
}


func (ui *UI) PrintSection(title string) {
	if ui.NoColor {
		fmt.Printf("\n--- %s ---\n", title)
	} else {
		fmt.Println()
		Primary.Printf("▶ %s\n", title)
		Secondary.Println(strings.Repeat("─", len(title)+2))
	}
}

func (ui *UI) PrintSubsection(title string) {
	if ui.NoColor {
		fmt.Printf("\n%s:\n", title)
	} else {
		fmt.Println()
		Accent.Printf("● %s\n", title)
	}
}

func (ui *UI) PrintListItem(item string, positive bool) {
	if ui.NoColor {
		if positive {
			fmt.Printf("    + %s\n", item)
		} else {
			fmt.Printf("    • %s\n", item)
		}
	} else {
		if positive {
			fmt.Printf("    ")
			Success.Printf("✓ ")
			fmt.Printf("%s\n", item)
		} else {
			fmt.Printf("    ")
			Warning.Printf("• ")
			fmt.Printf("%s\n", item)
		}
	}
}

func (ui *UI) PrintKeyValue(key, value string) {
	if ui.NoColor {
		fmt.Printf("  %-12s %s\n", key+":", value)
	} else {
		fmt.Printf("  ")
		Secondary.Printf("%-12s", key+":")
		fmt.Printf(" %s\n", value)
	}
}

func (ui *UI) PrintBanner() {
	if ui.NoColor {
		fmt.Println("Mux AI - Generative Engine Optimization Tool")
		fmt.Println("============================================")
		return
	}
	
	banner := `                           MUX AI                            
             Generative Engine Optimization Tool             
                                                             
  🚀 Local Analysis  🤖 LLM Integration  📊 Smart Reporting  `
	
	Primary.Println(banner)
	fmt.Println()
}

// FormatMarkdownContent formats markdown-like content for beautiful terminal display
func (ui *UI) FormatMarkdownContent(content string) string {
	if ui.NoColor {
		return content
	}
	
	lines := strings.Split(content, "\n")
	var result strings.Builder
	
	for _, line := range lines {
		formatted := ui.formatMarkdownLine(line)
		result.WriteString(formatted + "\n")
	}
	
	return result.String()
}

// formatMarkdownLine formats a single line with markdown-style syntax
func (ui *UI) formatMarkdownLine(line string) string {
	if ui.NoColor {
		return line
	}
	
	trimmed := strings.TrimSpace(line)
	
	// Headers
	if strings.HasPrefix(trimmed, "### ") {
		return "  " + H3.Sprint("▶ "+strings.TrimPrefix(trimmed, "### "))
	}
	if strings.HasPrefix(trimmed, "## ") {
		return "\n" + H2.Sprint("◆ "+strings.TrimPrefix(trimmed, "## ")) + "\n" + strings.Repeat("─", 50)
	}
	if strings.HasPrefix(trimmed, "# ") {
		return "\n" + H1.Sprint("■ "+strings.TrimPrefix(trimmed, "# ")) + "\n" + strings.Repeat("═", 60)
	}
	
	// Lists
	if strings.HasPrefix(trimmed, "- ") {
		return "  " + ListItem.Sprint("•") + " " + ui.formatInlineMarkdown(strings.TrimPrefix(trimmed, "- "))
	}
	if strings.HasPrefix(trimmed, "* ") {
		return "  " + ListItem.Sprint("•") + " " + ui.formatInlineMarkdown(strings.TrimPrefix(trimmed, "* "))
	}
	
	// Numbered lists
	if matched := regexp.MustCompile(`^(\d+)\. `).FindStringSubmatch(trimmed); len(matched) > 0 {
		return "  " + ListItem.Sprint(matched[1]+".") + " " + ui.formatInlineMarkdown(strings.TrimPrefix(trimmed, matched[0]))
	}
	
	// Checkboxes
	if strings.HasPrefix(trimmed, "- [x] ") || strings.HasPrefix(trimmed, "- [X] ") {
		return "  " + Checkmark.Sprint("✓") + " " + ui.formatInlineMarkdown(strings.TrimPrefix(trimmed, "- [x] "))
	}
	if strings.HasPrefix(trimmed, "- [ ] ") {
		return "  " + Subtle.Sprint("□") + " " + ui.formatInlineMarkdown(strings.TrimPrefix(trimmed, "- [ ] "))
	}
	
	// Tables
	if strings.Contains(trimmed, "|") && strings.Count(trimmed, "|") >= 2 {
		return ui.formatTableRow(trimmed)
	}
	
	// Blockquotes
	if strings.HasPrefix(trimmed, "> ") {
		return "  " + Quote.Sprint("│ "+strings.TrimPrefix(trimmed, "> "))
	}
	
	// Code blocks (simplified - just detect lines with lots of backticks or indentation)
	if strings.HasPrefix(trimmed, "```") {
		return CodeBlock.Sprint("  " + trimmed)
	}
	if strings.HasPrefix(line, "    ") && len(strings.TrimSpace(line)) > 0 {
		return CodeBlock.Sprint(line)
	}
	
	// Regular paragraphs
	if trimmed == "" {
		return ""
	}
	
	return "  " + ui.formatInlineMarkdown(trimmed)
}

// formatInlineMarkdown handles inline formatting like **bold**, *italic*, `code`
func (ui *UI) formatInlineMarkdown(text string) string {
	if ui.NoColor {
		return text
	}
	
	// Bold **text**
	boldRegex := regexp.MustCompile(`\*\*([^*]+)\*\*`)
	text = boldRegex.ReplaceAllStringFunc(text, func(match string) string {
		content := strings.Trim(match, "*")
		return Bold.Sprint(content)
	})
	
	// Italic *text*
	italicRegex := regexp.MustCompile(`\*([^*]+)\*`)
	text = italicRegex.ReplaceAllStringFunc(text, func(match string) string {
		content := strings.Trim(match, "*")
		return Italic.Sprint(content)
	})
	
	// Inline code `text`
	codeRegex := regexp.MustCompile("`([^`]+)`")
	text = codeRegex.ReplaceAllStringFunc(text, func(match string) string {
		content := strings.Trim(match, "`")
		return InlineCode.Sprint(content)
	})
	
	// Links [text](url) - simplified to just show as links
	linkRegex := regexp.MustCompile(`\[([^\]]+)\]\([^)]+\)`)
	text = linkRegex.ReplaceAllStringFunc(text, func(match string) string {
		parts := regexp.MustCompile(`\[([^\]]+)\]\([^)]+\)`).FindStringSubmatch(match)
		if len(parts) > 1 {
			return Link.Sprint(parts[1])
		}
		return match
	})
	
	return text
}

// formatTableRow formats a markdown table row for terminal display
func (ui *UI) formatTableRow(line string) string {
	if ui.NoColor {
		return "  " + line
	}
	
	// Check if this is a separator row (contains only |, -, :, spaces)
	separatorRegex := regexp.MustCompile(`^[\|\-:\s]+$`)
	if separatorRegex.MatchString(line) {
		// Format separator row with box drawing characters
		width := len(line)
		if width < 60 {
			width = 60
		}
		return "  " + Subtle.Sprint(strings.Repeat("─", width))
	}
	
	// Split the row into cells
	cells := strings.Split(line, "|")
	
	// Clean up cells (remove leading/trailing spaces)
	for i, cell := range cells {
		cells[i] = strings.TrimSpace(cell)
	}
	
	// Skip empty first/last cells (common in markdown tables)
	if len(cells) > 0 && cells[0] == "" {
		cells = cells[1:]
	}
	if len(cells) > 0 && cells[len(cells)-1] == "" {
		cells = cells[:len(cells)-1]
	}
	
	// Format cells with consistent width
	var formattedCells []string
	widths := []int{20, 8, -1} // Factor, Score, Key Finding (no limit for last column)
	
	for i, cell := range cells {
		if i < len(widths) && widths[i] > 0 {
			// Pad fixed-width columns (Factor and Score)
			formatted := ui.formatInlineMarkdown(cell)
			formattedCells = append(formattedCells, fmt.Sprintf("%-*s", widths[i], formatted))
		} else {
			// No width limit for Key Finding column - show full content
			formattedCells = append(formattedCells, ui.formatInlineMarkdown(cell))
		}
	}
	
	// Join with styled separators
	separator := Subtle.Sprint(" │ ")
	return "  " + strings.Join(formattedCells, separator)
}

// PrintMarkdownContent prints formatted markdown content
func (ui *UI) PrintMarkdownContent(content string) {
	formatted := ui.formatMarkdownContent(content)
	fmt.Print(formatted)
}

// formatMarkdownContent is the main formatting function
func (ui *UI) formatMarkdownContent(content string) string {
	return ui.FormatMarkdownContent(content)
}