// Package exitcode provides standardized exit codes for CLI commands.
package exitcode

// Standard exit codes for goossify commands.
const (
	Success        = 0 // Operation completed successfully
	Error          = 1 // Operation failed with error
	Warning        = 2 // Operation completed with warnings/partial success
	ValidationFail = 3 // Validation failed (e.g., score below threshold)
)

// Result represents a command execution result with exit code.
type Result struct {
	Code     int      `json:"exit_code"`
	Message  string   `json:"message,omitempty"`
	Warnings []string `json:"warnings,omitempty"`
	Errors   []string `json:"errors,omitempty"`
}

// NewSuccess creates a success result.
func NewSuccess(message string) *Result {
	return &Result{Code: Success, Message: message}
}

// NewError creates an error result.
func NewError(message string, errors ...string) *Result {
	return &Result{Code: Error, Message: message, Errors: errors}
}

// NewWarning creates a warning result.
func NewWarning(message string, warnings ...string) *Result {
	return &Result{Code: Warning, Message: message, Warnings: warnings}
}

// NewValidationFail creates a validation failure result.
func NewValidationFail(message string) *Result {
	return &Result{Code: ValidationFail, Message: message}
}

// HasWarnings returns true if there are warnings.
func (r *Result) HasWarnings() bool {
	return len(r.Warnings) > 0
}

// HasErrors returns true if there are errors.
func (r *Result) HasErrors() bool {
	return len(r.Errors) > 0
}

// AddWarning adds a warning to the result.
func (r *Result) AddWarning(warning string) {
	r.Warnings = append(r.Warnings, warning)
	if r.Code == Success {
		r.Code = Warning
	}
}

// AddError adds an error to the result.
func (r *Result) AddError(err string) {
	r.Errors = append(r.Errors, err)
	r.Code = Error
}
