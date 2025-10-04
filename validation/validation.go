package validation

import (
	"fmt"
	"html"
	"net/url"
	"regexp"
	"strconv"
	"strings"
)

// ValidationError represents a validation error
type ValidationError struct {
	Field   string
	Message string
}

func (e ValidationError) Error() string {
	return fmt.Sprintf("%s: %s", e.Field, e.Message)
}

// ValidationErrors represents multiple validation errors
type ValidationErrors []ValidationError

func (e ValidationErrors) Error() string {
	if len(e) == 0 {
		return ""
	}
	var msgs []string
	for _, err := range e {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// ValidateRequired checks if a string value is not empty
func ValidateRequired(value, fieldName string) error {
	if strings.TrimSpace(value) == "" {
		return ValidationError{Field: fieldName, Message: "is required"}
	}
	return nil
}

// ValidateMaxLength checks if a string value doesn't exceed maximum length
func ValidateMaxLength(value, fieldName string, maxLen int) error {
	if len(value) > maxLen {
		return ValidationError{Field: fieldName, Message: fmt.Sprintf("must be at most %d characters", maxLen)}
	}
	return nil
}

// ValidateMinLength checks if a string value meets minimum length
func ValidateMinLength(value, fieldName string, minLen int) error {
	if len(value) < minLen {
		return ValidationError{Field: fieldName, Message: fmt.Sprintf("must be at least %d characters", minLen)}
	}
	return nil
}

// ValidateRange checks if an integer value is within a range
func ValidateRange(value int, fieldName string, min, max int) error {
	if value < min || value > max {
		return ValidationError{Field: fieldName, Message: fmt.Sprintf("must be between %d and %d", min, max)}
	}
	return nil
}

// ValidateURL checks if a string is a valid URL
func ValidateURL(value, fieldName string) error {
	if value == "" {
		return nil // Allow empty URLs
	}
	u, err := url.Parse(value)
	if err != nil || u.Scheme == "" || u.Host == "" {
		return ValidationError{Field: fieldName, Message: "must be a valid URL"}
	}
	return nil
}

// ValidateEmail checks if a string is a valid email address
func ValidateEmail(value, fieldName string) error {
	if value == "" {
		return nil // Allow empty emails
	}
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	if !emailRegex.MatchString(value) {
		return ValidationError{Field: fieldName, Message: "must be a valid email address"}
	}
	return nil
}

// ValidateAlphanumeric checks if a string contains only alphanumeric characters
func ValidateAlphanumeric(value, fieldName string) error {
	if value == "" {
		return nil
	}
	alphanumericRegex := regexp.MustCompile(`^[a-zA-Z0-9]+$`)
	if !alphanumericRegex.MatchString(value) {
		return ValidationError{Field: fieldName, Message: "must contain only alphanumeric characters"}
	}
	return nil
}

// SanitizeString sanitizes a string by escaping HTML and trimming whitespace
func SanitizeString(value string) string {
	return strings.TrimSpace(html.EscapeString(value))
}

// SanitizeURL sanitizes a URL string
func SanitizeURL(value string) string {
	value = strings.TrimSpace(value)
	if value == "" {
		return value
	}

	// Parse and reconstruct URL to normalize it
	if u, err := url.Parse(value); err == nil {
		// Only allow http and https schemes
		if u.Scheme != "http" && u.Scheme != "https" {
			return ""
		}
		return u.String()
	}
	return ""
}

// ValidateAndSanitizeQueryParam validates and sanitizes a query parameter
func ValidateAndSanitizeQueryParam(query url.Values, paramName string, required bool, maxLen int) (string, error) {
	value := query.Get(paramName)

	if required && value == "" {
		return "", ValidationError{Field: paramName, Message: "is required"}
	}

	if value == "" {
		return "", nil
	}

	if maxLen > 0 {
		if err := ValidateMaxLength(value, paramName, maxLen); err != nil {
			return "", err
		}
	}

	return SanitizeString(value), nil
}

// ValidateAndSanitizeIntParam validates and sanitizes an integer query parameter
func ValidateAndSanitizeIntParam(query url.Values, paramName string, required bool, min, max int) (int, error) {
	value := query.Get(paramName)

	if required && value == "" {
		return 0, ValidationError{Field: paramName, Message: "is required"}
	}

	if value == "" {
		return 0, nil
	}

	intValue, err := strconv.Atoi(value)
	if err != nil {
		return 0, ValidationError{Field: paramName, Message: "must be a valid integer"}
	}

	if min != 0 || max != 0 { // Only validate range if bounds are specified
		if err := ValidateRange(intValue, paramName, min, max); err != nil {
			return 0, err
		}
	}

	return intValue, nil
}

// ValidateStruct validates a struct using reflection (placeholder for future enhancement)
func ValidateStruct(s interface{}) error {
	// This could be enhanced with struct tags and reflection
	// For now, return nil as we don't have complex structs to validate
	return nil
}
