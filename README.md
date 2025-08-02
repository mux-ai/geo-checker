# Mux AI - Generative Engine Optimization Tool

A powerful CLI tool for Generative Engine Optimization (GEO) that helps optimize content for AI-powered search engines and language models like ChatGPT, Claude, and other LLMs.

## Features

### 🆕 **Latest Updates**
- **🎯 Auto Provider Detection** - Automatically selects OpenAI if only OpenAI key available
- **📊 Universal Score Averaging** - LLM and Hybrid modes now always average Local + AI scores  
- **📋 Enhanced Table Rendering** - Beautiful formatted tables in terminal output
- **🔧 Improved Model Selection** - Better OpenAI model descriptions and recommendations
- **⚡ Smart Fallbacks** - More intelligent handling of API key availability

### 🧠 **Intelligent Analysis**
- **🎯 Auto Mode** - Automatically detects API keys and selects optimal provider/model
- **📊 Hybrid Scoring** - Always averages local + LLM scores for maximum accuracy 
- **🔄 Smart Fallbacks** - Gracefully handles API failures with local analysis backup
- **🎨 Beautiful Output** - Enhanced terminal formatting with table rendering and markdown support

### 🚀 **Core Capabilities**
- **🏠 Local Analysis** - Comprehensive GEO scoring without API requirements
- **🤖 LLM Integration** - Support for Claude, OpenAI GPT, and local LLMs
- **⚡ Bulk Processing** - Analyze multiple URLs concurrently
- **📁 Directory Scanning** - Scan local HTML files in project directories
- **🎯 Multiple Analysis Modes** - Auto, Local, LLM, or Hybrid analysis
- **📄 Multiple Output Formats** - Text, JSON, and Markdown formats

### 🎛️ **User Experience**
- **🔧 Zero Setup** - Works out of the box with intelligent mode selection
- **🎪 Interactive Mode** - Beautiful CLI for model and provider selection
- **📋 Model Management** - List and validate available models for each provider
- **🔍 Debug Tools** - Content extraction debugging and troubleshooting
- **💡 Smart Recommendations** - Context-aware suggestions for optimization

## Installation

### Prerequisites
- Go 1.19 or later

### Build from Source
```bash
git clone <repository-url>
cd mux-geo-checker
go mod tidy
go build -o mux-geo main.go
```

## Configuration

### 🎯 **Intelligent Setup (Recommended)**

**No configuration required!** The tool automatically detects your setup and chooses the best analysis method:

- ✅ **API Key Available** → Uses hybrid mode (Local + LLM averaged scores)
- ❌ **No API Key** → Uses local mode + shows upgrade recommendations

### 🔑 **API Keys (Optional for Enhanced Analysis)**

```bash
# For Claude - Latest models with excellent analysis
export CLAUDE_API_KEY="sk-ant-your-claude-api-key"

# For OpenAI - Most popular with cost-effective options  
export OPENAI_API_KEY="sk-proj-your-openai-api-key"

# For local LLM - Privacy-focused, no API costs, run locally
export OLLAMA_BASE_URL="http://localhost:11434"  # Optional, defaults to this
# No API key required - just install Ollama and pull models!
```

### ✅ **API Key Validation**

The tool automatically validates API key formats:
- **Claude**: Must start with `sk-ant-`
- **OpenAI**: Must start with `sk-` or `sk-proj-`
- **Local**: No API key required - just needs Ollama running locally

### 🏠 **Why Choose Local LLM?**

- **🔒 Privacy**: Your data never leaves your machine
- **💰 Cost-effective**: No API usage fees - unlimited analysis
- **🚀 Speed**: No network latency, direct local processing
- **🔌 Offline**: Works without internet connection
- **🛡️ Security**: Perfect for sensitive content analysis
- **⚙️ Control**: Choose your own models and configurations

### 🎪 **Interactive Setup & Auto-Detection**

```bash
# Let the tool guide you through setup
./mux-geo analyze <url> --interactive

# See available models and providers  
./mux-geo models
./mux-geo models openai
./mux-geo models claude

# Auto-provider detection (NEW!)
export OPENAI_API_KEY="your-key"  # Only OpenAI key set
./mux-geo analyze <url>           # Auto-detects and uses OpenAI
# ✓ Analysis complete! Score: 65/100 (Local: 29 + AI: 85, averaged)
```

## Usage

### 🚀 **Quick Start (Zero Setup)**

Get intelligent GEO analysis automatically:

```bash
# Auto mode - uses best available method
./mux-geo analyze https://example.com

# Interactive mode - guided setup
./mux-geo analyze https://example.com --interactive

# Debug content extraction issues
./mux-geo debug https://example.com

# Scan local files
./mux-geo scan ./website

# Bulk analyze URLs
./mux-geo bulk urls.txt
```

### 🤖 **Enhanced Analysis with LLMs**

For deeper insights with AI-powered analysis:

```bash
# Automatic hybrid analysis (when API key available)
export OPENAI_API_KEY="your-key"
./mux-geo analyze https://example.com  # Auto-detects and uses hybrid

# Force specific modes
./mux-geo analyze https://example.com --mode llm --provider claude
./mux-geo analyze https://example.com --mode hybrid --provider openai

# Interactive model selection
./mux-geo analyze https://example.com --interactive

# Advanced options
./mux-geo analyze https://example.com \
  --mode hybrid \
  --provider openai \
  --model gpt-4o \
  --output markdown
```

### Bulk URL Analysis

Create a file with URLs (one per line):

```bash
# urls.txt
https://example.com
https://docs.example.com
https://blog.example.com
```

Run bulk analysis:

```bash
./mux-geo bulk urls.txt --concurrent 3 --output json
```

### Directory Scanning

Scan a local directory for HTML files:

```bash
./mux-geo scan ./website --extensions .html,.htm --output markdown
```

## Command Options

### Global Options

- `--mode`: Analysis mode (`auto`, `local`, `llm`, `hybrid`) [default: auto]
- `--provider, -p`: LLM provider (`claude`, `openai`, `local`) [default: claude]
- `--model, -m`: Model to use (empty = recommended model)
- `--output, -o`: Output format (`text`, `json`, `markdown`) [default: text]
- `--interactive, -i`: Interactive model selection [default: false]

### New Commands

- `models [provider]`: List available models for providers
- `debug <url>`: Debug content extraction and analysis issues

### Bulk Command Options

- `--concurrent, -c`: Number of concurrent requests [default: 5]

### Scan Command Options

- `--extensions`: File extensions to scan [default: .html]

## Analysis Modes

### 🎯 **Auto Mode (Default & Recommended)**

- **🧠 Intelligent selection** - automatically chooses best method
- **📊 Hybrid when possible** - averages local + LLM scores for accuracy
- **🔄 Smart fallbacks** - uses local analysis when API unavailable
- **💡 Helpful guidance** - shows recommendations for enhanced analysis
- **⚡ Zero configuration** - works immediately without setup

### 🏠 **Local Mode**

- **🔧 No API keys required** - works completely offline
- **⚡ Instant results** - fast rule-based analysis
- **📊 Comprehensive scoring** across 5 key GEO factors:
  - Content Structure (25%) - heading hierarchy, organization
  - Semantic Clarity (25%) - readability, terminology
  - Context Richness (20%) - depth, examples, specifics
  - Authority Signals (15%) - citations, credibility
  - Accessibility (15%) - meta tags, structure
- **🎯 Detailed recommendations** for technical improvements
- **💨 Perfect for quick audits** and batch processing

### 🤖 **LLM Mode**

- **🔑 Requires API keys** for chosen provider
- **📊 Averages scores** - combines local rule-based + AI analysis automatically
- **🧠 AI-powered insights** using advanced language models
- **🎨 Beautiful formatting** with README-style output and table rendering
- **🔍 Contextual analysis** with nuanced understanding
- **📝 Natural language recommendations** with examples
- **🎯 Best for content strategy** and comprehensive optimization

### ⚖️ **Hybrid Mode (Best of Both)**

- **📊 Always averages scores** - mathematically combines local + LLM: `(Local + AI) / 2`
- **🔍 Dual perspective** - technical structure + content quality assessment
- **✅ Balanced accuracy** - LLM validates and enhances local analysis
- **🎯 Comprehensive insights** - both rule-based and AI-powered evaluation
- **🏆 Most accurate results** - recommended for professional use
- **📈 Score transparency** - shows breakdown: "Score: 65/100 (Local: 29 + AI: 78, averaged)"

## Supported LLM Providers

### 🧠 **Claude (Anthropic)**

- **API Key**: `CLAUDE_API_KEY` environment variable
- **Format**: Must start with `sk-ant-`
- **Recommended Models**:
  - ⭐ `claude-3-5-sonnet-20241022` - Latest, best for complex analysis
  - `claude-3-sonnet-20240229` - Balanced performance and cost
  - `claude-3-opus-20240229` - Most capable, higher cost
  - `claude-3-haiku-20240307` - Fastest, most economical

### 🤖 **OpenAI (GPT)**

- **API Key**: `OPENAI_API_KEY` environment variable
- **Format**: Must start with `sk-` or `sk-proj-`
- **Available Models**:
  - ⭐ `gpt-4o` - Latest multimodal model, best overall performance
  - `gpt-4o-mini` - Cost-effective, fast, excellent value for money
  - `gpt-4-turbo` - Large context window, strong reasoning capabilities
  - `gpt-4` - High quality reasoning, proven performance
  - `gpt-3.5-turbo` - Fast, economical, good for simple analysis tasks

### 🏠 **Local LLM (Ollama)**

- **Setup**: Runs locally, no API costs, privacy-focused
- **URL**: `OLLAMA_BASE_URL` or default `http://localhost:11434`
- **Available Models**:
  - ⭐ `llama2` - Open source, reliable, best for GEO analysis
  - `llama3` - Latest open source model, improved performance
  - `codellama` - Specialized for code analysis and technical content
  - `mistral` - Efficient alternative, good balance of speed/quality

#### **🛠️ Quick Ollama Setup**

**Step 1: Install Ollama**
```bash
# Linux/macOS/WSL
curl -fsSL https://ollama.com/install.sh | sh

# Windows: Download from https://ollama.com/download

# Docker alternative
docker run -d -v ollama:/root/.ollama -p 11434:11434 --name ollama ollama/ollama
```

**Step 2: Pull a Model**
```bash
# Recommended for GEO analysis
ollama pull llama2

# Or try the latest version
ollama pull llama3

# For technical content analysis
ollama pull codellama
```

**Step 3: Start Ollama**
```bash
# Start the Ollama server (if not auto-started)
ollama serve

# Verify it's running
curl http://localhost:11434/api/version
```

**Step 4: Use with GEO Checker**
```bash
# Use local provider with specific model
./mux-geo analyze https://example.com --provider local --model llama2

# Interactive selection
./mux-geo analyze https://example.com --interactive
# Choose: 3. local → 1. llama2

# See available local models
./mux-geo models local
```

### 📋 **Model Management**

```bash
# List all available models
./mux-geo models

# List models for specific provider
./mux-geo models openai
./mux-geo models claude
./mux-geo models local

# Interactive model selection
./mux-geo analyze <url> --interactive
```

## Output Formats

### 📄 **Text (Default)**

Beautiful terminal output with enhanced formatting:
- **🎨 README-style markdown** rendering in terminal
- **🌈 Colored headers** and visual hierarchy
- **📊 Score indicators** with visual progress
- **📋 Bullet points** with colored markers
- **💡 Code highlighting** for technical terms
- **📝 Quote blocks** for important insights

### JSON

Structured JSON output suitable for programmatic processing:

```json
{
  "url": "https://example.com",
  "title": "Example Page",
  "analysis": "...",
  "processed_at": "2024-01-01T12:00:00Z",
  "tokens_used": 1500,
  "metadata": {...}
}
```

### Markdown

Formatted markdown suitable for documentation and reports.

## Sample Output

```
                           MUX AI                            
             Generative Engine Optimization Tool             
                                                             
  🚀 Local Analysis  🤖 LLM Integration  📊 Smart Reporting  

▶ ANALYSIS DETAILS
──────────────────
  URL:         https://example.com
  Title:       Example Website
  Mode:        HYBRID
  Analyzed:    2024-01-15 10:30:45

▶ OVERALL SCORE
───────────────
  GEO Score:            78/100 (78.0%)
    📊 Averaged Score (Local: 72 + LLM: 84, averaged)

▶ DETAILED BREAKDOWN
────────────────────
  Content Structure:    85/100 (85.0%)
  Semantic Clarity:     72/100 (72.0%)
  Context Richness:     80/100 (80.0%)
  Authority Signals:    65/100 (65.0%)
  Accessibility:        88/100 (88.0%)

● Strengths
    ✓ Good heading hierarchy structure
    ✓ Well-organized content structure
    ✓ Content is clear and readable
    ✓ Good information density

● Recommendations
     1. Add more citations and credible references
     2. Include more expertise and credibility signals
     3. Define technical terms and concepts clearly
     4. Provide more context and background information

▶ AI INSIGHTS
─────────────

  Overall Score: 84/100

  ◆ 🎯 GEO Analysis Summary
  ──────────────────────────────────────────────────

  This webpage shows strong semantic foundations but needs 
  structural improvements for optimal GEO performance.

  ▶ ✅ Key Strengths

  • Excellent semantic clarity with consistent terminology
  • Clean content structure following logical flow
  • AI-friendly formatting with proper semantic elements

  ▶ 🚀 Priority Recommendations

  1. Implement proper heading hierarchy
  2. Add rich context and examples
  3. Enhance authority signals with citations

  │ Pro Tip: Focus on structural improvements first for
  │ maximum impact on AI understanding.
```

## Examples

### Example Files

See the `examples/` directory for:

- `urls.txt`: Sample URL list for bulk processing
- `sample.html`: Example HTML file for testing directory scanning

### Basic Usage Examples

```bash
# 🎯 Intelligent analysis (recommended)
./mux-geo analyze https://example.com  # Auto-detects best method
./mux-geo analyze https://example.com --interactive  # Guided setup
./mux-geo debug https://example.com  # Troubleshoot content issues

# 📊 Output formats
./mux-geo analyze https://example.com --output json > result.json
./mux-geo scan ./website --output markdown > report.md

# 🤖 Enhanced LLM analysis
export OPENAI_API_KEY="your-key"
./mux-geo analyze https://example.com  # Auto-uses hybrid mode
./mux-geo analyze https://example.com --mode llm --provider claude
./mux-geo analyze https://example.com --mode hybrid --provider openai --model gpt-4o

# 🏠 Local LLM analysis (privacy-focused, no API costs)
./mux-geo analyze https://example.com --provider local --model llama2
./mux-geo analyze https://example.com --mode hybrid --provider local --model llama3

# 📋 Model management
./mux-geo models  # List all available models
./mux-geo models openai  # List OpenAI models only

# ⚡ Bulk processing
./mux-geo bulk urls.txt --concurrent 10 --output json
./mux-geo bulk urls.txt --mode hybrid --provider claude --concurrent 3

# 📁 Directory scanning
./mux-geo scan ./src --extensions .html,.htm --mode auto
./mux-geo scan ./docs --mode hybrid --provider openai
```

## GEO Analysis

### Local Scoring Algorithm

The local analysis evaluates content across 5 key dimensions:

1. **Content Structure (25%)**
   - Heading hierarchy (H1 → H2 → H3)
   - Content organization and flow
   - Paragraph structure and length
   - Use of lists and bullet points

2. **Semantic Clarity (25%)**
   - Readability and sentence complexity
   - Terminology consistency
   - Definition clarity for technical terms
   - Unambiguous language usage

3. **Context Richness (20%)**
   - Content depth and detail level
   - Use of examples and specifics
   - Background information provision
   - Comprehensive coverage of topics

4. **Authority Signals (15%)**
   - Citations and references
   - Expertise indicators
   - Factual accuracy signals
   - Credible source integration

5. **Accessibility (15%)**
   - Meta information quality
   - Machine-readable structure
   - Information density balance
   - AI parsing friendliness

Each factor is scored 0-100, then weighted to produce an overall GEO score with specific, actionable recommendations.

### 🧠 **Intelligent Scoring System**

**Auto Mode** automatically provides the most accurate scoring:

| Scenario | Scoring Method | Formula | Accuracy |
|----------|----------------|---------|----------|
| 🔑 **API Key Available** | Always Averages | `(Local + LLM) / 2` | ⭐⭐⭐⭐⭐ Best |
| ❌ **No API Key** | Local + Recommendations | Rule-based only | ⭐⭐⭐ Good |
| 🔧 **Force Local** | Local Only | Rule-based only | ⭐⭐⭐ Technical |
| 🤖 **LLM/Hybrid Mode** | Always Averages | `(Local + LLM) / 2` | ⭐⭐⭐⭐⭐ Best |

**Why Averaging is Most Accurate:**
- **Local analysis** catches technical SEO structure issues
- **LLM analysis** evaluates actual content quality and meaning  
- **Mathematical averaging** provides balanced, unbiased assessment
- **Dual validation** ensures comprehensive coverage of all GEO factors
- **Consistent scoring** regardless of mode (LLM or Hybrid)

### 📊 **Score Interpretation**

- **90-100**: Excellent GEO optimization
- **80-89**: Good optimization with minor improvements needed
- **70-79**: Moderate optimization, several areas to improve
- **60-69**: Basic optimization, significant improvements needed
- **Below 60**: Poor optimization, major overhaul recommended

## Architecture

```
mux-geo-checker/
├── cmd/                  # CLI commands
│   ├── analyze.go        # Single URL analysis
│   ├── bulk.go           # Bulk URL processing
│   ├── root.go           # Root command & banner
│   └── scan.go           # Directory scanning
├── pkg/
│   ├── analyzer/         # Core analysis engine with intelligent mode selection
│   ├── config/           # Configuration management
│   ├── formatter/        # Enhanced output formatting with markdown support
│   ├── llm/              # LLM provider interfaces with error handling
│   │   ├── claude.go     # Anthropic Claude with validation
│   │   ├── openai.go     # OpenAI GPT with model management
│   │   ├── local.go      # Local LLM (Ollama) support
│   │   ├── provider.go   # Provider factory with auto-detection
│   │   ├── errors.go     # Structured error handling
│   │   └── interactive.go # Interactive model selection
│   ├── scorer/           # Local scoring algorithm
│   ├── scanner/          # Directory scanning
│   └── ui/               # Enhanced terminal UI with markdown rendering
├── internal/
│   ├── bulk/             # Bulk processing logic
│   └── webpage/          # Web scraping utilities
└── main.go               # Application entry point
```

## 🔧 **Troubleshooting**

### Common Issues

#### "Content cannot be empty" Error
```bash
# Debug content extraction
./mux-geo debug https://problematic-url.com

# Check if page requires JavaScript or has unusual structure
```

#### API Key Issues
```bash
# Check if key is set
echo $OPENAI_API_KEY

# Verify key format
# OpenAI: starts with "sk-" or "sk-proj-"
# Claude: starts with "sk-ant-"
```

#### Model Not Found
```bash
# List available models
./mux-geo models openai

# Use interactive selection
./mux-geo analyze <url> --interactive
```

#### Local LLM (Ollama) Issues
```bash
# Check if Ollama is running
curl http://localhost:11434/api/version

# Start Ollama if not running
ollama serve

# List installed models
ollama list

# Pull a model if none installed
ollama pull llama2

# Test Ollama directly
curl http://localhost:11434/api/generate -d '{
  "model": "llama2",
  "prompt": "Test prompt",
  "stream": false
}'
```

### Score Discrepancies

- **Low local, high LLM**: Good content, poor structure → Fix technical issues
- **High local, low LLM**: Good structure, poor content → Improve content quality  
- **Both low**: Needs comprehensive optimization

### 🆘 **Getting Help**

```bash
# Check command help
./mux-geo --help
./mux-geo analyze --help
./mux-geo models --help

# Debug content extraction
./mux-geo debug <url>

# Test with example content
./mux-geo analyze https://example.com
```

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests if applicable
5. Submit a pull request

## License

MIT License - see LICENSE file for details.

## 📚 **Additional Documentation**

- **[Error Handling](pkg/llm/errors.go)** - Comprehensive error types and handling

## 🎯 **Quick Reference**

| Task | Command | Notes |
|------|---------|-------|
| **Basic analysis** | `./mux-geo analyze <url>` | Auto-detects best method |
| **With API key** | `export OPENAI_API_KEY="key" && ./mux-geo analyze <url>` | Uses hybrid mode |
| **Local LLM** | `./mux-geo analyze <url> --provider local --model llama2` | Privacy-focused, no API costs |
| **Interactive setup** | `./mux-geo analyze <url> -i` | Guided model selection |
| **List models** | `./mux-geo models` | Shows all available models |
| **Local models only** | `./mux-geo models local` | Shows Ollama models |
| **Debug issues** | `./mux-geo debug <url>` | Troubleshoot content extraction |
| **Bulk analysis** | `./mux-geo bulk urls.txt` | Process multiple URLs |
| **Local files** | `./mux-geo scan ./website` | Analyze local HTML files |

## Support

For issues and feature requests, please use the GitHub issues tracker.