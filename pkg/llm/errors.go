package llm

import (
	"errors"
	"fmt"
	"net/http"
	"strings"
)

// Error types for better error handling
var (
	ErrInvalidCredentials = errors.New("invalid API credentials")
	ErrRateLimited       = errors.New("rate limit exceeded")
	ErrQuotaExceeded     = errors.New("quota exceeded")
	ErrModelNotFound     = errors.New("model not found")
	ErrInvalidRequest    = errors.New("invalid request")
	ErrServiceUnavailable = errors.New("service temporarily unavailable")
	ErrTimeout           = errors.New("request timeout")
	ErrNetworkError      = errors.New("network error")
	ErrInvalidResponse   = errors.New("invalid response format")
	ErrContentFiltered   = errors.New("content filtered by provider")
)

// LLMError represents a structured error from LLM operations
type LLMError struct {
	Type       ErrorType              `json:"type"`
	Message    string                 `json:"message"`
	Provider   string                 `json:"provider"`
	Model      string                 `json:"model,omitempty"`
	StatusCode int                    `json:"status_code,omitempty"`
	Details    map[string]interface{} `json:"details,omitempty"`
	Retryable  bool                   `json:"retryable"`
}

type ErrorType string

const (
	ErrorTypeAuth         ErrorType = "authentication"
	ErrorTypeRateLimit    ErrorType = "rate_limit"
	ErrorTypeQuota        ErrorType = "quota"
	ErrorTypeModel        ErrorType = "model"
	ErrorTypeRequest      ErrorType = "request"
	ErrorTypeService      ErrorType = "service"
	ErrorTypeTimeout      ErrorType = "timeout"
	ErrorTypeNetwork      ErrorType = "network"
	ErrorTypeResponse     ErrorType = "response"
	ErrorTypeContent      ErrorType = "content"
	ErrorTypeUnknown      ErrorType = "unknown"
)

func (e *LLMError) Error() string {
	if e.Provider != "" {
		return fmt.Sprintf("[%s] %s: %s", e.Provider, e.Type, e.Message)
	}
	return fmt.Sprintf("%s: %s", e.Type, e.Message)
}

func (e *LLMError) Is(target error) bool {
	switch target {
	case ErrInvalidCredentials:
		return e.Type == ErrorTypeAuth
	case ErrRateLimited:
		return e.Type == ErrorTypeRateLimit
	case ErrQuotaExceeded:
		return e.Type == ErrorTypeQuota
	case ErrModelNotFound:
		return e.Type == ErrorTypeModel
	case ErrInvalidRequest:
		return e.Type == ErrorTypeRequest
	case ErrServiceUnavailable:
		return e.Type == ErrorTypeService
	case ErrTimeout:
		return e.Type == ErrorTypeTimeout
	case ErrNetworkError:
		return e.Type == ErrorTypeNetwork
	case ErrInvalidResponse:
		return e.Type == ErrorTypeResponse
	case ErrContentFiltered:
		return e.Type == ErrorTypeContent
	}
	return false
}

// NewLLMError creates a new LLMError
func NewLLMError(errorType ErrorType, message, provider string) *LLMError {
	return &LLMError{
		Type:      errorType,
		Message:   message,
		Provider:  provider,
		Retryable: isRetryable(errorType),
	}
}

// ParseHTTPError converts HTTP status codes to appropriate LLM errors
func ParseHTTPError(statusCode int, body []byte, provider string) *LLMError {
	bodyStr := string(body)
	
	switch statusCode {
	case http.StatusUnauthorized:
		return &LLMError{
			Type:       ErrorTypeAuth,
			Message:    "Invalid API key or authentication failed",
			Provider:   provider,
			StatusCode: statusCode,
			Retryable:  false,
		}
	case http.StatusForbidden:
		if strings.Contains(strings.ToLower(bodyStr), "quota") {
			return &LLMError{
				Type:       ErrorTypeQuota,
				Message:    "API quota exceeded",
				Provider:   provider,
				StatusCode: statusCode,
				Retryable:  false,
			}
		}
		return &LLMError{
			Type:       ErrorTypeAuth,
			Message:    "Access forbidden - check permissions",
			Provider:   provider,
			StatusCode: statusCode,
			Retryable:  false,
		}
	case http.StatusNotFound:
		return &LLMError{
			Type:       ErrorTypeModel,
			Message:    "Model not found or endpoint not available",
			Provider:   provider,
			StatusCode: statusCode,
			Retryable:  false,
		}
	case http.StatusTooManyRequests:
		return &LLMError{
			Type:       ErrorTypeRateLimit,
			Message:    "Rate limit exceeded - please wait before retrying",
			Provider:   provider,
			StatusCode: statusCode,
			Retryable:  true,
		}
	case http.StatusBadRequest:
		return &LLMError{
			Type:       ErrorTypeRequest,
			Message:    parseRequestError(bodyStr),
			Provider:   provider,
			StatusCode: statusCode,
			Retryable:  false,
		}
	case http.StatusInternalServerError, http.StatusBadGateway, http.StatusServiceUnavailable, http.StatusGatewayTimeout:
		return &LLMError{
			Type:       ErrorTypeService,
			Message:    "Service temporarily unavailable",
			Provider:   provider,
			StatusCode: statusCode,
			Retryable:  true,
		}
	default:
		return &LLMError{
			Type:       ErrorTypeUnknown,
			Message:    fmt.Sprintf("Unexpected error (status: %d): %s", statusCode, truncateBody(bodyStr, 200)),
			Provider:   provider,
			StatusCode: statusCode,
			Retryable:  statusCode >= 500,
		}
	}
}

// isRetryable determines if an error type is retryable
func isRetryable(errorType ErrorType) bool {
	switch errorType {
	case ErrorTypeRateLimit, ErrorTypeService, ErrorTypeTimeout, ErrorTypeNetwork:
		return true
	default:
		return false
	}
}

// parseRequestError extracts meaningful error messages from request errors
func parseRequestError(body string) string {
	body = strings.ToLower(body)
	
	if strings.Contains(body, "content policy") || strings.Contains(body, "safety") {
		return "Content was filtered due to safety policies"
	}
	if strings.Contains(body, "token") && strings.Contains(body, "limit") {
		return "Request exceeds token limits"
	}
	if strings.Contains(body, "model") {
		return "Invalid model specification"
	}
	if strings.Contains(body, "parameter") {
		return "Invalid request parameters"
	}
	
	return "Invalid request format"
}

// truncateBody truncates response body for error messages
func truncateBody(body string, maxLength int) string {
	if len(body) <= maxLength {
		return body
	}
	return body[:maxLength] + "..."
}

// WrapNetworkError wraps network-related errors
func WrapNetworkError(err error, provider string) *LLMError {
	return &LLMError{
		Type:      ErrorTypeNetwork,
		Message:   fmt.Sprintf("Network error: %v", err),
		Provider:  provider,
		Retryable: true,
	}
}

// WrapTimeoutError wraps timeout errors
func WrapTimeoutError(err error, provider string) *LLMError {
	return &LLMError{
		Type:      ErrorTypeTimeout,
		Message:   fmt.Sprintf("Request timeout: %v", err),
		Provider:  provider,
		Retryable: true,
	}
}

// WrapResponseError wraps response parsing errors
func WrapResponseError(err error, provider string) *LLMError {
	return &LLMError{
		Type:      ErrorTypeResponse,
		Message:   fmt.Sprintf("Invalid response format: %v", err),
		Provider:  provider,
		Retryable: false,
	}
}