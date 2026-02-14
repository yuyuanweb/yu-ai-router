package yuairoutersdk

import "fmt"

type APIError struct {
	Code    int
	Message string
}

func (e *APIError) Error() string {
	return fmt.Sprintf("yu-ai-router error(code=%d): %s", e.Code, e.Message)
}

type AuthError struct {
	Message string
}

func (e *AuthError) Error() string {
	return "auth error: " + e.Message
}

type RateLimitError struct {
	Message string
}

func (e *RateLimitError) Error() string {
	return "rate limit error: " + e.Message
}
